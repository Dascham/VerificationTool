#include <kissnet.hpp>
#include <iostream>

#include "config.h"

namespace kn = kissnet;


int main() {
    std::cout << "Starting" << std::endl;
    //setup socket
    kn::tcp_socket server(kn::endpoint("0.0.0.0:" + std::to_string(MASTER_PORT)));
    std::cout << "Bind" << std::endl;
    server.bind();
    std::cout << "Listen" << std::endl;
    server.listen();
    std::cout << "Listening" << std::endl;


    //Wait for connections

    std::array<kn::tcp_socket, WORKER_COUNT> connections;

    std::cout << "Waiting for all workers to connect..." << std::endl;
    for (int i = 0; i < WORKER_COUNT; ++i) {
        connections[i] = server.accept();
        const auto &conn = connections[i];
        const auto &bindLoc = conn.get_bind_loc();
        std::cout << "Accept: " << i << " " << bindLoc.address << ":" << bindLoc.port << std::endl;
    }


    // TODO: send start msg
    std::cout << "Not sending any start signal right now..." << std::endl;


    // TODO: wait for finish OR JUST WAIT FOR PORT TO CLOSE?

    std::cout << "Done" << std::endl;

    //auto client = server.accept();

    //Read once in a 1k buffer
    /*kn::buffer<1024> buff;
    const auto [size, status] = client.recv(buff);

    //Add null terminator, and print as string
    if(size < buff.size()) buff[size] = std::byte{ 0 };
    std::cout << reinterpret_cast<const char*>(buff.data()) << '\n';*/

    return 0;
}