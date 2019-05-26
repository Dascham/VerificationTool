#ifndef CPP_DISTRIBUTEDMODELCHECKER_H
#define CPP_DISTRIBUTEDMODELCHECKER_H

#include "BaseModelChecker.h"

#include <iostream>
#include <thread>
#include <chrono>

// TODO: does this work on linux too?
#define AI_ADDRCONFIG LUP_ADDRCONFIG
const char* inet_ntop(int af, const void* src, char* dst, int cnt);
#undef NOMINMAX
#include <kissnet.hpp>
const char* inet_ntop(int af, const void* src, char* dst, int cnt){

    struct sockaddr_in srcaddr;

    memset(&srcaddr, 0, sizeof(struct sockaddr_in));
    memcpy(&(srcaddr.sin_addr), src, sizeof(srcaddr.sin_addr));

    srcaddr.sin_family = af;
    if (WSAAddressToString((struct sockaddr*) &srcaddr, sizeof(struct sockaddr_in), 0, dst, (LPDWORD) &cnt) != 0) {
        DWORD rv = WSAGetLastError();
        printf("WSAAddressToString() : %d\n",rv);
        return NULL;
    }
    return dst;
}

#include "Packet.h"
#include "Block.h"
#include "SocketThread.h"

#include "config.h"

using namespace std::chrono_literals;
namespace kn = kissnet;

namespace modelcheckers {

    // This class is intended to be the worker of the modelchecker. It depends on a master to coordinate.
    class DistributedModelChecker : public BaseModelChecker {
    private:
        std::array<SocketThread, WORKER_COUNT-1> socketThreads;
        std::array<kn::tcp_socket, WORKER_COUNT> outgoingSockets;
        std::array<size_t, WORKER_COUNT> pendingAcks = {};

        kn::tcp_socket masterSocket{kn::endpoint{localhost, MASTER_PORT}};

        const size_t workerID;
        bool running = true;

        Herd herd; // TODO: this, pass to ctor?

        MasterPacket receiveMasterPacket() {
            kn::buffer<2> buffer;
            auto [size, status] = masterSocket.recv(buffer);

            assert(size == 2);

            MasterPacket masterPacket{(MasterPacket::Type)buffer[0], (uint8_t)buffer[1]};
            return masterPacket;
        }

        void sendWorkerPacket(WorkerPacket workerPacket) {
            kn::buffer<2> buffer;
            buffer[0] = static_cast<std::byte>(workerPacket.type);
            buffer[1] = static_cast<std::byte>(workerPacket.data);

            masterSocket.send(buffer);
        }

        size_t consumePendingAcks(size_t otherID) {
            if (otherID == workerID) {
                assert(pendingAcks[otherID] == 0);
                return 0;
            }

            kn::tcp_socket &socket{outgoingSockets[otherID]};

            assert(socket.is_valid());

            while (socket.bytes_available() > 0) {
                kn::buffer<1> buffer;
                const auto[data_size, status_code] = socket.recv(buffer);

                for (int i = 0; i < data_size; ++i) {
                    printf("!!!%zx %x\n", static_cast<size_t>(buffer[i]), 0x42);
                    assert(buffer[i] == static_cast<std::byte>(0x42));
                }

                assert(data_size <= pendingAcks[otherID]);
                pendingAcks[otherID] -= data_size;
            }

            return pendingAcks[otherID];
        }

        void sendState(State state, kn::port_t otherID) {
            // Don't send to self using socket
            if (otherID == workerID) {
                addNewState(state);
                std::cout << "Kept: " << state << std::endl;
                return;
            }

            kn::tcp_socket &socket{outgoingSockets[otherID]};
            pendingAcks[otherID]++;

            assert(state.locations.size() == MODEL_AUTOMATA);
            assert(state.variables.size() == MODEL_VARIABLES);

            kn::buffer<MODEL_AUTOMATA> locationBuffer;
            kn::buffer<MODEL_VARIABLES> variableBuffer;

            for (int i = 0; i < locationBuffer.size(); ++i) {
                locationBuffer[i] = static_cast<std::byte>(state.locations[i]);
            }

            for (int i = 0; i < variableBuffer.size(); ++i) {
                variableBuffer[i] = static_cast<std::byte>(state.variables[i]);
            }

            socket.send(locationBuffer);
            socket.send(variableBuffer);

            std::cout << "Sent: " << state << std::endl;

            consumePendingAcks(otherID);
        }

    protected:
        void
        handleNewState(const State &newState, const State &oldState, std::vector<size_t> changedLocations) override {
            if constexpr (!USE_HERD_HASHING) {
                // If we don't care about herds, then we simply hash and send
                sendState(newState, std::hash<State>{}(newState) % WORKER_COUNT);
                return;
            }

            bool shouldHash = false;
            for (const auto &loc : changedLocations) {
                // TODO: get old-/new-Group
                if (oldState.locations[loc]/*TODO: is cover*/ || oldState.locations[loc] != newState.locations[loc]/*TODO: the herd of old/new*/) {
                    shouldHash = true;
                    break;
                }
            }

            if (shouldHash) {
                sendState(newState, std::hash<State>{}(newState) % WORKER_COUNT);
            } else { // Keep the state locally
                addNewState(newState);
            }

        }

        void addStateQueue(std::queue<State> newQueue) {
            while (!newQueue.empty()) {
                addNewState(newQueue.front());
                newQueue.pop();
            }
        }

        void exploreStateQueue() {
            while (!stateQueue.empty()) {
                State current = stateQueue.front();
                stateQueue.pop();
                generateSuccessors(current);
            }
        }

    public:

        enum class DoneState {
            NotDone, FirstDone, SecondDone
        };

        void checkModel() override {

            printf("Starting to listen for incoming worker connections.\n\n");
            std::thread incomingThread{[this]() {
                printf("--incomingThread\n");
                auto sin1 = kn::tcp_socket{
                        kn::endpoint{"0.0.0.0", static_cast<kn::port_t>(WORKER_PORT_FIRST + workerID)}};

                try {
                    sin1.bind();
                    sin1.listen();
                } catch (std::runtime_error &ex) {
                    std::cerr << "Failed to start listener for worker connections(try waiting a while): \n\t"
                              << ex.what() << std::endl;
                    exit(1);
                }

                printf("--Ready to accept incoming sockets! (%s:%u)\n", sin1.get_bind_loc().address.c_str() ,sin1.get_bind_loc().port);
                for (size_t i = 0; i < WORKER_COUNT - 1; ++i) {
                    printf("--i: %zu\n", i);
                    auto sock = sin1.accept();
                    printf("--accepted: %zu\n", i);
                    printf("--%zu/%zu workers connected... (%s:%u)\n", i + 1, WORKER_COUNT - 1,
                           sock.get_bind_loc().address.c_str(), sock.get_bind_loc().port);

                    socketThreads[i].assignSocket(std::move(sock));
                }

                printf("--All workers connected!\n");

                sin1.close();
            }};

            // Connect to master
            printf("Connecting to master... \n");
            while (!masterSocket.connect()); // Keep trying until we succeed (or we run out of cake)
            printf("Connected to master!\n");

            printf("Connecting to other workers...\n");
            for (size_t i = 0; i < WORKER_COUNT; ++i) {
                if (i == workerID) {
                    printf("Skipping own port, socket %zu/%zu.\n", i + 1, WORKER_COUNT);
                    continue;
                }

                outgoingSockets[i] = kn::tcp_socket{kn::endpoint{localhost, kn::port_t(WORKER_PORT_FIRST + i)}};
                for (size_t j = 1; /*j <= 3*/; ++j) {
                    const auto &sock = outgoingSockets[i];
                    sock.set_non_blocking(true);

                    printf("Connecting to worker %zu, attempt %zu. (%s:%u)\n", i, j,
                           sock.get_bind_loc().address.c_str(), sock.get_bind_loc().port);

                    if (outgoingSockets[i].connect()) {
                        printf("Connected to workers %zu/%zu. (%s:%u)\n", i + 1, WORKER_COUNT,
                               sock.get_bind_loc().address.c_str(), sock.get_bind_loc().port);
                        break;
                    } else if (j == 3) {
                        fprintf(stderr, "Error connecting to worker id %zu!\n", i);
                        exit(1);
                    }
                }


            }
            printf("All outgoing sockets connected!\n\n");

            printf("Waiting for all workers to connect...\n");
            incomingThread.join();

            printf("Ready!\n\n");
            //addInitialState();

            DoneState doneState = DoneState::NotDone;
            uint8_t doneID = std::numeric_limits<uint8_t>::max();

            while (running) {
                bool didWork = false;

                if (!stateQueue.empty()) {
                    didWork = true;
                    exploreStateQueue();
                }

                // Loop through SocketThreads and take their state queue
                for (auto &socketThread : socketThreads) {
                    socketThread.stealQueue(stateQueue);
                    //addStateQueue(socketThread.stealQueue());


                    if (!stateQueue.empty()) {
                        didWork = true;
                        exploreStateQueue();
                    }
                }

                if (didWork) {
                    doneState = DoneState::NotDone;
                }

                /*for (int i = 0; i < 3; ++i)*/  /*if (masterSocket.bytes_available() < 2) std::this_thread::sleep_for(.1s); else*/
                if (masterSocket.bytes_available() >= 2) {

                    MasterPacket masterPacket = receiveMasterPacket();//{MasterPacket::Type::Terminate, 0};

                    switch (masterPacket.type) {
                        case MasterPacket::Type::InitialState:
                            printf(">Initial State:\n"
                                   "\tAdding initial state to queue!\n");
                            addInitialState();
                            goto continue_running;
                        case MasterPacket::Type::IsDone:
                            printf(">IsDone?:\n");
                            {
                                if (!didWork) {

                                    // Check if we still have any pending acks for states first
                                    bool stillPendingAcks = false;
                                    for (size_t i = 0; i < WORKER_COUNT; ++i) {
                                        if (consumePendingAcks(i) > 0) {
                                            stillPendingAcks = true;
                                            //printf("Pending %zu: %zu\n", i, pendingAcks[i]);
                                            //std::this_thread::sleep_for(100ms);
                                            break;
                                        }
                                    }

                                    if (!stillPendingAcks) {
                                        printf("\tResponding FirstDone(%u) to master!\n", masterPacket.data);
                                        doneState = DoneState::FirstDone;
                                        doneID = masterPacket.data;

                                        WorkerPacket response{WorkerPacket::Type::FirstDone, doneID};
                                        sendWorkerPacket(response);
                                        goto continue_running;
                                    }
                                }


                                printf("\tResponding NotDone(%u) at IsDone to master!\n", masterPacket.data);
                                WorkerPacket response{WorkerPacket::Type::NotDone, masterPacket.data};
                                sendWorkerPacket(response);

                            }
                            goto continue_running;
                        case MasterPacket::Type::IsStillDone:
                            printf(">IsStillDone?:\n");
                            if (doneState != DoneState::FirstDone || doneID != masterPacket.data) {
                                doneState = DoneState::NotDone;

                                printf("\tResponding NotDone(%u) to master!\n", masterPacket.data);
                                WorkerPacket response{WorkerPacket::Type::NotDone, masterPacket.data};
                                sendWorkerPacket(response);
                            } else {

                                //std::this_thread::sleep_for(5000ms);
                                printf("\tChecking SocketThread queues one last time...\n");
                                for (auto &socketThread : socketThreads) {

                                    //printf("\t\tWait for one SocketThread iteration...\n");
                                    socketThread.flag = false;
                                    while(!socketThread.flag) std::this_thread::yield(); // Wait for one iteration of the SocketThread loop
                                    //printf("\t\t\tSocketThread iteration confirmed!\n");

                                    //auto queue = socketThread.stealQueue(std::queue<State>());
                                    socketThread.stealQueue(stateQueue);

                                    if (!stateQueue.empty()) {
                                        printf("\tResponding NotDone(%u) to master!\n", masterPacket.data);
                                        WorkerPacket response{WorkerPacket::Type::NotDone, masterPacket.data};
                                        sendWorkerPacket(response);

                                        //addStateQueue(queue);

                                        goto continue_running;
                                    }
                                }

                                printf("\tResponding SecondDone(%u) to master!\n", masterPacket.data);
                                doneState = DoneState::SecondDone;
                                assert(doneID == masterPacket.data);

                                WorkerPacket response{WorkerPacket::Type::SecondDone, masterPacket.data};
                                sendWorkerPacket(response);
                            }

                            // Let it continue until it is told to terminate or receives a State.
                            goto continue_running;
                        case MasterPacket::Type::Terminate:
                            printf(">Terminate:\n");
                            if (didWork || doneState != DoneState::SecondDone) {
                                fprintf(stderr, "\tWorker terminated early!\n"
                                                "\t\tdidWork: %s\n"
                                                "\t\tdoneState: %s\n",
                                        didWork ? "true" : "false",
                                        (doneState == DoneState::NotDone) ? "NotDone" :
                                        (doneState == DoneState::FirstDone) ? "FirstDone" :
                                        (doneState == DoneState::SecondDone) ? "SecondDone" :
                                        "INVALID");
                                exit(1);
                            } else {
                                printf("\tTerminating worker successfully!\n");
                                goto finish_running;
                            }
                        default:
                            fprintf(stderr, ">Invalid packet from master (%u)\n", (uint8_t) masterPacket.type);
                            exit(1);
                    }
                }

                continue_running:;

            }
            finish_running:

            printf("Joining all SocketThreads...\n");
            for (size_t i = 0; i < socketThreads.size(); ++i) {
                printf("Trying to join SocketThread %zu...\n", i);
                socketThreads[i].join();
            }

            printf("Model Checker finished!\n");

            printStatistics();

        }


        explicit DistributedModelChecker(size_t workerID, model::Model model, Herd herd = {{}})
                : workerID{workerID}, BaseModelChecker{std::move(model)}, herd{std::move(herd)} {}
    };

}

#endif //CPP_DISTRIBUTEDMODELCHECKER_H
