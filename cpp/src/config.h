#ifndef CPP_CONFIG_H
#define CPP_CONFIG_H

#include <kissnet.hpp>
#include "testmodels.h"


namespace kn = kissnet;

constexpr kn::port_t MASTER_PORT = 20700;
constexpr kn::port_t WORKER_PORT_FIRST = 20701;

constexpr size_t WORKER_COUNT = 4;

const auto localhost = "127.0.0.1";

constexpr bool USE_HERD_HASHING = true;

/*
constexpr size_t MODEL_AUTOMATA = 3;
constexpr size_t MODEL_VARIABLES = MODEL_AUTOMATA;
const model::Model THE_MODEL = testmodels::testModel(MODEL_AUTOMATA, MODEL_VARIABLES);
*/


// "Bad" model
/*
constexpr size_t MODEL_AUTOMATA = 3;
constexpr size_t MODEL_VARIABLES = MODEL_AUTOMATA;
const model::Model THE_MODEL = testmodels::badModel(MODEL_AUTOMATA);
*/

// "Good" model

constexpr size_t MODEL_AUTOMATA = 3;
constexpr size_t MODEL_VARIABLES = MODEL_AUTOMATA;
const model::Model THE_MODEL = testmodels::goodModel(MODEL_AUTOMATA,false);

// "Good Sync" model
/*
constexpr size_t MODEL_AUTOMATA = 8;
constexpr size_t MODEL_VARIABLES = MODEL_AUTOMATA;
const model::Model THE_MODEL = testmodels::goodModel(MODEL_AUTOMATA-1, true);
*/
#endif //CPP_CONFIG_H
