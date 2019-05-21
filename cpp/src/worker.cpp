#include <iostream>
#include <thread>
#include <chrono>

#include <kissnet.hpp>

#include "config.h"

using namespace std::chrono_literals;
namespace kn = kissnet;

int main() {
    const size_t MYID = 1;


    // parse params(worker id)


    // Connect to master
    kn::tcp_socket a_socket(kn::endpoint("0.0.0.0:" + std::to_string(MASTER_PORT)));
    a_socket.connect();




    // TODO: THESE
    // Start listening (MAYBE BEFORE RECEVING FROM MASTER)
    //   Wait for n connections (thread manager thread?)
    // Also open n connections

    // TODO: ASSERT model loc,var sizes == constants


    //Create a kissnet tco ipv4 socket


    //Create a "GET /" HTTP request, and send that packet into the socket
    auto get_index_request = std::string{"GET / HTTP/1.1\r\nHost: avalon.ybalird.info\r\n\r\n"};

    //Send request
    a_socket.send(reinterpret_cast<const std::byte *>(get_index_request.c_str()), get_index_request.size());

    //Receive data into a buffer
    kn::buffer<4096> static_buffer;

    //Useless wait, just to show how long the response was
    std::this_thread::sleep_for(1s);

    //Print how much data our OS has for us
    std::cout << "bytes available to read : " << a_socket.bytes_available() << '\n';

    //Get the data, and the lengh of data
    const auto[data_size, status_code] = a_socket.recv(static_buffer);

    //To print it as a good old C string, add a null terminator
    if (data_size < static_buffer.size())
        static_buffer[data_size] = std::byte{'\0'};

    //Print the raw data as text into the terminal (should display html/css code here)
    std::cout << reinterpret_cast<const char *>(static_buffer.data()) << '\n';


    //So long, and thanks for all the fish
    return 0;
}