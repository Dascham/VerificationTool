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
    std::atomic<bool> running;

    std::thread thread;
    kn::tcp_socket socket;

    std::mutex mutex{};
    std::queue<State> queue{};


    void run() {
        while (running) {
            if (!socket.is_valid()) {
                fprintf(stderr, "Socket is invalid\n");
                exit(1);
            }

            if (true) {
                const size_t locSize = sizeof(uint8_t) * MODEL_AUTOMATA;
                const size_t varSize = sizeof(int8_t) * MODEL_VARIABLES;

                std::vector<uint8_t> locations(MODEL_AUTOMATA);
                std::vector<int8_t> variables(MODEL_VARIABLES);

                //printf("Receiving state...\n");
                for (size_t bytesRead = 0; bytesRead < MODEL_AUTOMATA; /**/) {
                    kn::buffer<1*sizeof(uint8_t)> buffer;
                    const auto [sizeLoc, statusLoc] = socket.recv(buffer);
                    if (sizeLoc > 0) {
                        assert(sizeLoc == sizeof(uint8_t));
                        locations[bytesRead] = static_cast<uint8_t>(buffer[0]);
                        ++bytesRead;
                    }
                }
                for (size_t bytesRead = 0; bytesRead < MODEL_VARIABLES; /**/) {
                    kn::buffer<1*sizeof(uint8_t)> buffer;
                    const auto [sizeVar, statusVar] = socket.recv(buffer);
                    if (sizeVar > 0) {
                        assert(sizeVar == sizeof(int8_t));
                        variables[bytesRead] = static_cast<int8_t>(buffer[0]);
                        ++bytesRead;
                    }
                }

                {
                    State newState{locations, variables};
                    std::unique_lock lock{mutex}; // Lock until leaving scope
                    queue.emplace(newState);
                }

                // send ACK
                kn::buffer<1> ackBuffer{static_cast<std::byte>(0x42)};
                socket.send(ackBuffer);

            }
        }
    }

public:
    // Used to hand over the queue to the worker thread
    void stealQueue(std::queue<State> &stateQueue) {
        std::unique_lock lock{mutex};
        std::swap(queue, stateQueue);
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
