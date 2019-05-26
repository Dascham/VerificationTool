#ifndef CPP_SOCKETTHREAD_H
#define CPP_SOCKETTHREAD_H

#include <thread>
#include <mutex>
#include <queue>
#include <atomic>

#include <kissnet.hpp>
#include <iomanip>
#include "config.h"
#include "State.h"

namespace kn = kissnet;

class SocketThread {
public:
    std::atomic<bool> flag;
private:
    std::atomic<bool> running;

    std::thread thread;
    kn::tcp_socket socket;

    std::mutex mutex{};
    std::queue<State> queue{};


    void run() {
        while (running) {
            flag = true;

            if (!socket.is_valid()) {
                fprintf(stderr, "Socket is invalid\n");
                exit(1);
            }

            if (socket.bytes_available() >= sizeof(uint8_t) * MODEL_AUTOMATA + sizeof(int8_t) * MODEL_VARIABLES) {

                kn::buffer<sizeof(uint8_t) * MODEL_AUTOMATA> locationBuffer;
                kn::buffer<sizeof(int8_t) * MODEL_AUTOMATA> variableBuffer;

                const auto [sizeLoc, statusLoc] = socket.recv(locationBuffer);
                if (statusLoc != kn::socket_status::valid) {
                    std::cerr << "Socket not valid when receiving location vector, error: " << std::hex << statusLoc << "bytes: " << sizeLoc << std::endl;
                    exit(1);
                }
                assert(sizeLoc == locationBuffer.size());
                assert(sizeLoc == THE_MODEL.automata.size());

                const auto [sizeVar, statusVar] = socket.recv(variableBuffer);
                if (statusLoc != kn::socket_status::valid) {
                    std::cerr << "Socket not valid when receiving variable vector, error: " << std::hex << statusVar << "bytes: " << sizeVar << std::endl;
                    exit(1);
                }
                assert(sizeVar == variableBuffer.size());
                assert(sizeLoc == THE_MODEL.numVariables);

                // Print received state
                if (true) {
                    std::cout << "Recv State:\n"
                                 "Loc: ";

                    for (const auto &loc : locationBuffer) {
                        //std::cout << static_cast<uint8_t>(loc) << " ";
                        printf("%u ", static_cast<uint8_t>(loc));
                    }

                    std::cout << "\n"
                                 "Var: ";

                    for (const auto &var : variableBuffer) {
                        //std::cout << static_cast<int8_t>(var) << " ";
                        printf("%d ", static_cast<int8_t>(var));

                    }

                    std::cout << "\n" << std::endl;
                }

                std::vector<uint8_t> locations(locationBuffer.size());
                std::vector<int8_t> variables(variableBuffer.size());



                for (int i = 0; i < locationBuffer.size(); ++i) {
                    locations[i] = static_cast<uint8_t>(locationBuffer[i]);
                }
                for (int i = 0; i < variableBuffer.size(); ++i) {
                    variables[i] = static_cast<int8_t>(variableBuffer[i]);
                }


                {
                    State newState{locations, variables};
                    std::unique_lock lock{mutex}; // Lock until leaving scope
                    queue.emplace(newState);
                }

                // send ACK
                kn::buffer<1> ackBuffer{static_cast<std::byte>(0x42)};
                socket.send(ackBuffer);

            } else if (socket.bytes_available() != 0) {
                // Just to debug if it gets stuck at some partial state
                //printf("SOCKBYTES: %zu/%zu\n", socket.bytes_available(), ((sizeof(uint8_t) * MODEL_AUTOMATA + sizeof(int8_t) * MODEL_VARIABLES)) );
            }
        }
    }

public:
    // Used to hand over the queue to the worker thread
    /*std::queue<State> &&*/ void stealQueue(std::queue<State> &stateQueue) {
        std::unique_lock lock{mutex};
        std::swap(queue, stateQueue);
        //auto stolenQueue = std::move(queue);
        //queue = {};
        //return std::move(stolenQueue);
        //return std::move(queue);
    }

    void join() {
        if (!running) throw std::logic_error("SocketThread not running");
        if (!thread.joinable()) throw std::logic_error("SocketThread not joinable");

        running = false;
        thread.join();

        socket.close();
    }

    void assignSocket(kn::tcp_socket newSocket) {
        if (running) throw std::logic_error("SocketThread already running while trying to assign new socket");
        if (thread.joinable()) throw std::logic_error("SocketThread was still joinable while trying to assign new socket");

        if (socket.is_valid()) std::clog << "Assigned new socket to SocketThread that already had a valid socket" << std::endl;

        assert( (std::unique_lock{mutex, std::defer_lock}.try_lock()) ); // Mutex should not be locked
        socket = std::move(newSocket);

        running = true;
        thread = std::thread(&SocketThread::run, this);
    }

    explicit SocketThread() : running{false} {}
};

#endif //CPP_SOCKETTHREAD_H
