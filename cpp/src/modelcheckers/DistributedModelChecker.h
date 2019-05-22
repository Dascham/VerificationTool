#ifndef CPP_DISTRIBUTEDMODELCHECKER_H
#define CPP_DISTRIBUTEDMODELCHECKER_H

#include "BaseModelChecker.h"

#include <iostream>
#include <thread>
#include <chrono>

#include <kissnet.hpp>

#include "SocketThread.h"

#include "config.h"

using namespace std::chrono_literals;
namespace kn = kissnet;

namespace modelcheckers {
    
    // This class is intended to be the worker of the modelchecker. It depends on a master to coordinate.
    class DistributedModelChecker : public BaseModelChecker {

        std::vector<SocketThread> socketThreads{WORKER_COUNT - 1};
        std::vector<kn::tcp_socket> outgoingSockets{WORKER_COUNT};

        const size_t workerID;

    protected:
        bool addNewState(const State &state) override {

            //if shouldHashLocal(state); need old state? TODO: continue implementing

            return BaseModelChecker::addNewState(state);
        }

    public:
        void checkModel() override {

            // Connect to master
            printf("Connecting to master... \n");
            kn::tcp_socket a_socket(kn::endpoint("0.0.0.0", MASTER_PORT));
            while (!a_socket.connect()); // Keep trying until we succeed (or we run out of cake)

            printf("Connected to master.\n\n");

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
                    printf("%zu/%zu workers connected...\n", i + 1, WORKER_COUNT - 1);
                }

                printf("All workers connected!\n\n");

                sin1.close();
            }};

            kn::tcp_socket mout{kn::endpoint{localhost, MASTER_PORT}};
            std::cout << mout.connect() << std::endl;

            std::this_thread::sleep_for(1s);

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
            printf("All outgoing sockets connected.\n");

            /*
             * kn::tcp_socket out1{kn::endpoint{localhost, WORKER_PORT_FIRST + 0}};
            std::cout << "connecting out1" << std::endl;
            std::cout << out1.connect() << std::endl; // TODO: check return value
            std::cout << "connected out1" << std::endl;
            */

            printf("Waiting for all workers to connect...\n");
            incomingThread.join();

            printf("Joining all SocketThreads.\n");
            for (size_t i = 0; i < socketThreads.size(); ++i) {
                std::cout << "Trying to join SocketThread %zu..." << std::endl;
                socketThreads[i].join();
                printf("Joined SocketThread %zu", i);
            }




            addInitialState();

            while (!stateQueue.empty()) {
                State current = stateQueue.front();
                stateQueue.pop();

                int generatedCounter = generateSuccessors(current);
            }

        }

        explicit DistributedModelChecker(size_t workerID, model::Model model)
                : workerID{workerID}, BaseModelChecker{std::move(model)} {}
    };

}

#endif //CPP_DISTRIBUTEDMODELCHECKER_H
