#include <iostream>
#include <queue>
#include <unordered_set>
#include <cassert>

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

        assert(state.locations.size() == model.automata.size()); // State should have one location index per automata

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

    bool checkInvariants(const State &state) {
        for (int i = 0; i < model.automata.size(); ++i) {
            const model::Automaton &automaton = model.automata[i];

            const size_t loc_i = state.locations[i];
            const model::Location &loc = automaton.locations[loc_i];

            if (!loc.invariant(state)) return false;
        }

        return true;
    }

    size_t generateSuccessors(const State &state) {
        size_t count = 0;
        for (int i = 0; i < model.automata.size(); ++i) { // For each automaton
            const model::Automaton &automaton = model.automata[i];

            const size_t loc_i = state.locations[i];
            const model::Location &loc = automaton.locations[loc_i];

            for (const auto& edge : loc.edges) { // For each edge from the current location in the automaton

                // Guard
                if (edge.guard(state)) { // If guard is satisfied in the current state

                    // TODO: Find Sync partner here if needed? Or just perform update, then use old state to find partner

                    // Update
                    State newState = edge.update(state);
                    newState.locations[i] = edge.destination;

                    assert(edge.destination < automaton.locations.size());
                    assert(newState.locations.size() == state.locations.size());
                    assert(newState.variables.size() == state.variables.size());

                    // Invariants for new state
                    if (checkInvariants(newState)) {
                        addNewState(newState);
                        count++;
                    } // TODO: count both valid and non valid states
                }
            }
        }

        return count;
    }

    void addInitialState() {
        State initialState{};
        initialState.locations.resize(model.automata.size());
        initialState.variables.resize(model.numVariables);

        for (int i = 0; i < model.automata.size(); ++i) {
            initialState.locations[i] = model.automata[i].initialLocaton;
        }

        addNewState(initialState);
    }

public:
    void checkModel() {
        addInitialState();

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
        numVars,
        {
            Automaton{{
                Location{Invariant{}, {
                    Edge{1, Guard{{
                        Predicate{Term{Term::Type::Variable, 0}, Predicate::ComparisonOperator::LessThan, Term{Term::Type::Constant, 127}}
                    }}, Update{{
                        {Assignment{0, Assignment::AssignOperator::IncAssign, Term{Term::Type::Constant, 1}}}
                    }}},

                }},
                Location{Invariant{}, {
                    Edge{0, Guard{{
                        Predicate{Term{Term::Type::Variable, 0}, Predicate::ComparisonOperator::LessThanEqual, Term{Term::Type::Constant, 127}}
                    }}}
                }}
            }}
        }
    }};
    simpleModelChecker.checkModel();

    return 0;
}