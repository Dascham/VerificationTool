#ifndef CPP_BASEMODELCHECKER_H
#define CPP_BASEMODELCHECKER_H

#include <queue>
#include <unordered_set>
#include <iostream>

#include "util.h"
#include "State.h"
#include "model/Model.h"

namespace modelcheckers {

    struct Statistics {
        size_t exploredCounter = 0;
        size_t generatedCounter = 1; // Start with initial state "generated"
        size_t duplicateCounter = 0;
    };

    class BaseModelChecker {
    protected:
            Statistics statistics;

            model::Model model;

            std::queue<State> stateQueue{};
            std::unordered_set<State> encounteredStates{};

            virtual bool addNewState(const State &state);
            bool checkInvariants(const State &state);
            size_t generateSuccessors(const State &state);

            void addInitialState();

    public:
            virtual void checkModel() = 0;

            explicit BaseModelChecker(model::Model model) : model{std::move(model)} {}
    };

}

#endif //CPP_BASEMODELCHECKER_H
