#ifndef CPP_CONFIG_H
#define CPP_CONFIG_H

#include <kissnet.hpp>

namespace kn = kissnet;

constexpr kn::port_t MASTER_PORT = 20700;
constexpr kn::port_t WORKER_PORT_FIRST = 20701;

constexpr size_t WORKER_COUNT = 1;

constexpr size_t MODEL_AUTOMATA = 2;
constexpr size_t MODEL_VARIABLES = 2;

const auto localhost = "127.0.0.1";

#endif //CPP_CONFIG_H
