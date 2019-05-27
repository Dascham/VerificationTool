#ifndef CPP_MODEL_H
#define CPP_MODEL_H

#include <utility>
#include <vector>
#include <string>

#include "State.h"
#include "Location.h"
#include "Edge.h"
#include "Block.h"

namespace model {

    struct Model;
    struct Automaton;

    struct Model {
        uint8_t numVariables;
        std::vector<Automaton> automata;

        Model(uint8_t numVariables, std::vector<Automaton> automata)
            : numVariables{numVariables}, automata{std::move(automata)} {}
    };

    struct Automaton {
        std::vector<Location> locations;
        size_t initialLocaton;

        Partitioning partitioning;

        explicit Automaton(std::vector<Location> locations, size_t initialLocation = 0, Partitioning partitioning = {})
            : locations{std::move(locations)}, initialLocaton{initialLocation}, partitioning{std::move(partitioning)} {}
    };

}


#endif //CPP_MODEL_H
