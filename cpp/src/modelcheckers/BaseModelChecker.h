#ifndef CPP_BASEMODELCHECKER_H
#define CPP_BASEMODELCHECKER_H

#include <queue>
#include <unordered_set>
#include <iostream>

#include "util.h"
#include "State.h"
#include "Group.h"
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

            bool addNewState(const State &state);
            virtual void handleNewState(const State &newState, const State &oldState, std::vector<size_t> changedLocations);
            bool checkInvariants(const State &state);
            void generateSuccessors(const State &state);

            void addInitialState();

            void printStatistics();
    public:
            virtual void checkModel() = 0;

            explicit BaseModelChecker(model::Model model, const Herd& herd = {}) : model{std::move(model)} {}
    };

}

#endif //CPP_BASEMODELCHECKER_H
