#include <iostream>

#include <kissnet.hpp>

#include "config.h"
#include "Packet.h"

namespace kn = kissnet;

class Master {

    std::array<kn::tcp_socket, WORKER_COUNT> connections{};

    void broadcastPacket(const MasterPacket &masterPacket) {
        kn::buffer<2> buffer;
        buffer[0] = static_cast<std::byte>(masterPacket.type);
        buffer[1] = static_cast<std::byte>(masterPacket.data);

        for (auto &connection : connections) {
            connection.send(buffer);
        }
    }

    WorkerPacket receivePacket(size_t socket) {
        kn::buffer<2> buffer;
        auto packet = connections[socket].recv(buffer);

        WorkerPacket workerPacket{(WorkerPacket::Type)buffer[0], (uint8_t)buffer[1]};
        return workerPacket;
    }

    template<MasterPacket::Type Question, WorkerPacket::Type ExpectedAnswer>
    bool checkDone(uint8_t i) {
        MasterPacket question{Question, i};
        broadcastPacket(question);

        // Wait for all workers to respond
        bool anyNotDone = false;

        bool respArray[WORKER_COUNT] = {};
        size_t respCount = 0;

        while (respCount < WORKER_COUNT) {
            for (size_t j = 0; j < WORKER_COUNT; ++j) {

                if (respArray[j]) continue; // If we already received a response

                WorkerPacket response = receivePacket(j); //{WorkerPacket::Type::FirstDone, i};

                if (response.data != i) continue; // We are only looking for responses to the current question

                j = true;
                ++respCount;

                if (response.type == WorkerPacket::Type::NotDone) {
                    anyNotDone = true;
                } else if (response.type != ExpectedAnswer) {
                    fprintf(stderr, "Unexpected positive answer to question, expected:%u got:%u\n",
                            (uint8_t) ExpectedAnswer, (uint8_t) response.type);
                    exit(1);
                }
            }
        }

        return !anyNotDone;
    }

public:
    void run() {

        printf("Opening socket to listen for connections.\n");
        //setup socket
        kn::tcp_socket server(kn::endpoint("0.0.0.0:" + std::to_string(MASTER_PORT)));
        server.bind();
        server.listen();

        //Wait for connections


        printf("Waiting for all workers to connect...\n");
        for (int i = 0; i < WORKER_COUNT; ++i) {
            connections[i] = server.accept();
            const auto &conn = connections[i];
            const auto &bindLoc = conn.get_bind_loc();
            std::cout << "\tAccept: " << i << " " << bindLoc.address << ":" << bindLoc.port << std::endl;
        }
        printf("All workers connected!\n\n");

        printf("Telling workers to start from initial state.\n\n");
        MasterPacket initialPacket{MasterPacket::Type::InitialState};
        broadcastPacket(initialPacket);

        printf("Waiting for workers to be done...\n");
        for (uint8_t i = 0;; ++i) {

            // Confirm workers are done in two stages (blocking calls, wait for all responses)
            printf("\tChecking if workers are done(%u)...\n", i);
            if (!checkDone<MasterPacket::Type::IsDone, WorkerPacket::Type::FirstDone>(i)) continue;
            printf("\tConfirming if workers are done(%u)...\n", i);
            if (!checkDone<MasterPacket::Type::IsStillDone, WorkerPacket::Type::SecondDone>(i)) continue;

            printf("All workers are done!\n\n");

            // Terminate all workers
            MasterPacket terminatePacket{MasterPacket::Type::Terminate, i};
            broadcastPacket(terminatePacket);

            printf("Exiting successfully\n");
            exit(0);
        }
    }

};

int main() {
    Master master;
    master.run();
}