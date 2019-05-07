#include <iostream>
#include <queue>
#include <unordered_set>

#include "util.h"
#include "State.h"
#include "model/Model.h"

class SimpleModelChecker {
    model::Model model;

    std::queue<State> stateQueue{};
    std::unordered_set<State> encounteredStates{};

    bool addNewState(const State &state) {
        auto ret = encounteredStates.insert(state);
        bool isNew = ret.second;

        if (isNew) {
            std::cout << "New state encountered; adding to queue...\n"
                      << state << std::endl;
            stateQueue.push(state);
        } else {
            std::cout << "Already encountered state...\n"
                      << state << std::endl;
        }

        return isNew;
    }

    size_t generateSuccessors(const State &state) {
        size_t count = 0;
        for (int i = 0; i < state.locations.size(); ++i) {
            State newState = state;
            newState.locations[i] += 1;
            newState.locations[i] %= 16;

            addNewState(newState);
            count++;
        }

        return count;
    }

public:
    void checkModel() {
        addNewState({{0,1,2},{0,0,0}});

        size_t exploredCounter = 0;
        size_t generatedCounter = 1; // Start with initial state "generated"

        while (!stateQueue.empty()) {
            State current = stateQueue.front();
            stateQueue.pop();

            exploredCounter++;
            generatedCounter += generateSuccessors(current);
        }

        std::cout << "Explored states: " << exploredCounter << std::endl;
        std::cout << "Generated states: " << generatedCounter << std::endl;
        std::cout << "Duplicate states: " << generatedCounter-exploredCounter << std::endl;
    }

    explicit SimpleModelChecker(model::Model model) : model{std::move(model)} {}
};

int main() {
    std::cout << "Welcome. Welcome to the simple_model_checker.\n"
                 "\n"
                 "You have chosen, or been chosen, to use one of our finest non-distributed model checkers.\n"
              << std::endl;

    using namespace model;

    constexpr size_t numVars = 1;

    SimpleModelChecker simpleModelChecker{Model{
        std::vector<int8_t>(numVars), // NOTE: Using parenthesis to create vector of given size
        {
            Automaton{{
                Location{Invariant{}, {
                    Edge{1, Guard{{}}, Update{{
                        std::vector<Assignment>{Assignment{1, Assignment::AssignOperator::Assign, Term{Term::Type::Constant, 42}}}
                    }}},

                }},
                Location{Invariant{}, {
                    Edge{0}
                }}
            }}/*,
            Automaton{{}, {}}*/
        } // Empty vector of automata
    }};
    simpleModelChecker.checkModel();

    return 0;
}