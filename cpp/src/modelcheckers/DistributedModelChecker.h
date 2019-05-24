#ifndef CPP_DISTRIBUTEDMODELCHECKER_H
#define CPP_DISTRIBUTEDMODELCHECKER_H

#include "BaseModelChecker.h"

#include <iostream>
#include <thread>
#include <chrono>

#include <kissnet.hpp>

#include "Packet.h"
#include "SocketThread.h"

#include "config.h"

using namespace std::chrono_literals;
namespace kn = kissnet;

namespace modelcheckers {

    // This class is intended to be the worker of the modelchecker. It depends on a master to coordinate.
    class DistributedModelChecker : public BaseModelChecker {
    private:
        std::vector<SocketThread> socketThreads{WORKER_COUNT - 1};
        std::vector<kn::tcp_socket> outgoingSockets{WORKER_COUNT};

        kn::tcp_socket masterSocket{kn::endpoint{localhost, MASTER_PORT}};

        const size_t workerID;
        bool running = true;


        MasterPacket receiveMasterPacket() {
            kn::buffer<2> buffer;
            auto packet = masterSocket.recv(buffer);

            MasterPacket masterPacket{(MasterPacket::Type)buffer[0], (uint8_t)buffer[1]};
            return masterPacket;
        }

        void sendWorkerPacket(WorkerPacket workerPacket) {
            kn::buffer<2> buffer;
            buffer[0] = static_cast<std::byte>(workerPacket.type);
            buffer[1] = static_cast<std::byte>(workerPacket.data);

            masterSocket.send(buffer);
        }

        void sendState(State state, kn::tcp_socket &socket) {
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
        }

    protected:

        void addStateQueue(std::queue<State> newQueue) {
            while (!newQueue.empty()) {
                addNewState(newQueue.front());
                newQueue.pop();
            }
        }

        bool addNewState(const State &state) override {

            //if shouldHashLocal(state); need old state? TODO: continue implementing

            return BaseModelChecker::addNewState(state);
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

                for (size_t i = 0; i < WORKER_COUNT - 1; ++i) {
                    socketThreads[i].assignSocket(sin1.accept());
                    printf("--%zu/%zu workers connected...\n", i + 1, WORKER_COUNT - 1);
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
                    printf("Skipping own port, socket %zu/%zu.\n", i + 1,
                           WORKER_COUNT); // TODO: explicit check to avoid sending to self(and crashing?)
                    continue;
                }

                outgoingSockets[i] = kn::tcp_socket{kn::endpoint{localhost, kn::port_t(WORKER_PORT_FIRST + i)}};
                if (outgoingSockets[i].is_valid()) {
                    printf("Connected to workers %zu/%zu.\n", i + 1, WORKER_COUNT);
                } else {
                    fprintf(stderr, "Error connecting to worker id %zu!\n", i);
                    exit(1);
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

                if (!stateQueue.empty()) exploreStateQueue();

                // Loop through SocketThreads and take their state queue
                for (auto &socketThread : socketThreads) {
                    addStateQueue(socketThread.stealQueue());

                    if (!stateQueue.empty()) {
                        didWork = true;
                        exploreStateQueue();
                    }
                }

                if (didWork) {
                    doneState = DoneState::NotDone;
                }

                //bool justDone = false;

                /*for (int i = 0; i < 3; ++i)*/if (masterSocket.bytes_available() < 2) std::this_thread::sleep_for(.1s); else { // TODO: while master has bytes?
                    // TODO: actually read from master(respond to is(still)done)

                    //kn::buffer<2> buffer;
                    //auto packet = masterSocket.recv(buffer);

                    MasterPacket masterPacket = receiveMasterPacket();//{MasterPacket::Type::Terminate, 0};
                    switch (masterPacket.type) {
                        case MasterPacket::Type::InitialState:
                            printf(">Initial State:\n"
                                   "\tAdding initial state to queue!\n");
                            addInitialState();
                            goto continue_running;
                        case MasterPacket::Type::IsDone:
                            printf(">IsDone?:\n");
                            if (!didWork) {
                                printf("\tResponding FirstDone(%u) to master!\n", masterPacket.data);
                                doneState = DoneState::FirstDone;
                                doneID = masterPacket.data;

                                WorkerPacket response{WorkerPacket::Type::FirstDone, masterPacket.data};
                                sendWorkerPacket(response);
                            } else {
                                printf("\tResponding NotDone(%u) at IsDone to master!\n", masterPacket.data);
                                WorkerPacket response{WorkerPacket::Type::NotDone, masterPacket.data};
                                sendWorkerPacket(response);
                            }
                            goto continue_running;
                        case MasterPacket::Type::IsStillDone:
                            printf(">IsStillDone?:\n");
                            if (doneState != DoneState::FirstDone) {
                                printf("\tResponding NotDone(%u) to master!\n", masterPacket.data);
                                WorkerPacket response{WorkerPacket::Type::NotDone, masterPacket.data};
                                sendWorkerPacket(response);
                            } else {
                                printf("\tChecking SocketThread queues one last time...\n");
                                for (auto &socketThread : socketThreads) {
                                    auto queue = socketThread.stealQueue();

                                    if (!queue.empty()) {
                                        printf("\tResponding NotDone(%u) to master!\n", masterPacket.data);
                                        WorkerPacket response{WorkerPacket::Type::NotDone, masterPacket.data};
                                        sendWorkerPacket(response);

                                        addStateQueue(queue);

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
                            fprintf(stderr, ">Invalid packet from master (%u)", (uint8_t) masterPacket.type);
                            exit(1);
                    }
                }

                continue_running:;

            }
            finish_running:

            printf("Joining all SocketThreads...\n");
            for (size_t i = 0; i < socketThreads.size(); ++i) {
                printf("Trying to join SocketThread %zu...", i);
                socketThreads[i].join();
            }

            printf("Model Checker finished!\n");

            printStatistics();

        }

        explicit DistributedModelChecker(size_t workerID, model::Model model)
                : workerID{workerID}, BaseModelChecker{std::move(model)} {}
    };

}

#endif //CPP_DISTRIBUTEDMODELCHECKER_H
