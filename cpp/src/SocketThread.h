#ifndef CPP_SOCKETTHREAD_H
#define CPP_SOCKETTHREAD_H

#include <thread>
#include <mutex>
#include <queue>
#include <atomic>

#include <kissnet.hpp>
#include <iomanip>
#include "config.h"

namespace kn = kissnet;

class SocketThread {
    std::thread thread;
    kn::tcp_socket socket;

    std::mutex mutex{};
    std::queue<State> queue{};

    std::atomic<bool> running;

    void run() {

        while (running) {
            if (socket.bytes_available() > sizeof(uint8_t) * MODEL_AUTOMATA + sizeof(int8_t) * MODEL_VARIABLES) {

                kn::buffer<sizeof(uint8_t) * MODEL_AUTOMATA> locationBuffer;
                kn::buffer<sizeof(int8_t) * MODEL_AUTOMATA> variableBuffer;

                const auto [sizeLoc, statusLoc] = socket.recv(locationBuffer);
                if (statusLoc != kn::socket_status::valid) {
                    std::cout << "Socket not valid when receiving location vector, error: " << std::hex << statusLoc << "bytes: " << sizeLoc << std::endl;
                }

                const auto [sizeVar, statusVar] = socket.recv(locationBuffer);
                if (statusLoc != kn::socket_status::valid) {
                    std::cout << "Socket not valid when receiving variable vector, error: " << std::hex << statusVar << "bytes: " << sizeVar << std::endl;
                }

                std::cout << "Recv State:\n"
                             "Loc: ";

                for (const auto &loc : locationBuffer) {
                    std::cout << static_cast<unsigned int>(loc) << " ";
                }

                std::cout << "\n"
                             "Var: ";

                for (const auto &var : variableBuffer) {
                    std::cout << static_cast<int>(var) << " ";
                }

                std::cout << "\n" << std::endl;

                std::vector<uint8_t> locations(locationBuffer.size());
                std::vector<int8_t> variables(variableBuffer.size());

                std::copy(locationBuffer.begin(), locationBuffer.end(), std::back_inserter(locations));
                std::copy(variableBuffer.begin(), variableBuffer.end(), std::back_inserter(variables));


                {
                    std::unique_lock lock{mutex}; // Lock until leaving scope
                    queue.emplace(locations, variables);
                }
            }
        }
    }

public:
    // Used to hand over the queue to the worker thread
    std::queue<State> &&stealQueue() {
        std::unique_lock lock{mutex};
        return std::move(queue);
    }

    void join() {
        if (!thread.joinable()) {
            throw std::logic_error("Tried to terminate thread while it is not joinable");
        }

        // TODO: Close socket

        running = false;
        thread.join();

        // return success status? or assume at this point everything went well
    }

    explicit SocketThread(kn::tcp_socket socket) : running{true}, socket{std::move(socket)} {
        // TODO: open socket?

        thread = std::thread(&SocketThread::run, this);
    }
};

#endif //CPP_SOCKETTHREAD_H
