#include "BaseModelChecker.h"

namespace modelcheckers {

    void BaseModelChecker::handleNewState(const State &newState, const State &oldState,
                                          std::vector<size_t> changedLocations) {
        addNewState(newState); // Default behavior is to simply add it to the local queue
    }

    bool BaseModelChecker::addNewState(const State &state) {
        if (state.locations.size() != model.automata.size() || state.variables.size() != model.numVariables) {
            fprintf(stderr, "locs:%zu == %zu    vars:%zu == %u\n",
                    state.locations.size(), model.automata.size(),
                    state.variables.size(), model.numVariables);

            fprintf(stderr, "DEBUG SKIPPING!\n"); // TODO: remove
            return false;
        }

        auto ret = encounteredStates.insert(state);
        bool isNew = ret.second;

        //assert(state.locations.size() ==
        //       model.automata.size()); // State should have one location index per automata

        if (isNew) {
            //std::cout << "New state encountered; adding to queue...\n"
            //          << state << std::endl;
            stateQueue.push(state);
        } else {
            //std::cout << "Already encountered state...\n"
            //          << state << std::endl;
        }

        return isNew;
    }

    bool BaseModelChecker::checkInvariants(const State &state) {
        for (int i = 0; i < model.automata.size(); ++i) {
            const model::Automaton &automaton = model.automata[i];

            const size_t loc_i = state.locations[i];
            const model::Location &loc = automaton.locations[loc_i];

            if (!loc.invariant(state)) return false;
        }

        return true;
    }

    void BaseModelChecker::generateSuccessors(const State &state) {
        using namespace model;

        statistics.exploredCounter++;

        for (size_t i = 0; i < model.automata.size(); ++i) { // For each automaton
            const Automaton &automaton = model.automata[i];

            const size_t loc_i = state.locations[i];
            const Location &loc = automaton.locations[loc_i];

            // TODO: use indexed for loop to catch edge number in case of exception(model error)
            for (const auto &edge : loc.edges) { // For each edge from the current location in the automaton

                // Guard
                if (edge.guard(state)) { // If guard is satisfied in the current state

                    if (edge.sync.type == Sync::Type::Recv) continue; // Nothing to do here

                    // Update
                    State newState = edge.update(state);
                    newState.locations[i] = edge.destination;

                    assert(edge.destination < automaton.locations.size());
                    assert(newState.locations.size() == state.locations.size());
                    assert(newState.variables.size() == state.variables.size());

                    // Handle synchronizing edges
                    if (edge.sync.type == Sync::Type::Send) {
                        for (size_t j = 0; j < model.automata.size(); ++j) { // For each automaton

                            if (i == j) continue; // Can't sync with itself(would result in double- or phantom-step?)

                            const Automaton &syncedAutomaton = model.automata[j];

                            const size_t syncedLoc_j = state.locations[j];
                            const Location &syncedLoc = syncedAutomaton.locations[syncedLoc_j];

                            // TODO: use indexed for loop to catch edge number in case of exception(model error)
                            for (auto &syncedEdge : syncedLoc.edges) {
                                if (syncedEdge.sync.type == Sync::Type::Recv &&
                                    syncedEdge.sync.channel == edge.sync.channel &&
                                    syncedEdge.guard(state)) // Check if guard was satisfied in the old state
                                {

                                    State syncedState = syncedEdge.update(newState);
                                    syncedState.locations[j] = syncedEdge.destination;

                                    assert(syncedEdge.destination < syncedAutomaton.locations.size());
                                    assert(syncedState.locations.size() == state.locations.size());
                                    assert(syncedState.variables.size() == state.variables.size());

                                    if (checkInvariants(syncedState)) {
                                        handleNewState(syncedState, state, {i,j});
                                        statistics.generatedCounter++;
                                    } // TODO: count both valid and non valid states
                                }
                            }
                        }
                    } else { // Handle non-synchronizing edges

                        // Invariants for new state TODO: DRY
                        if (checkInvariants(newState)) {
                            handleNewState(newState, state, {i});
                            statistics.generatedCounter++;
                        } // TODO: count both valid and non valid states
                    }
                }
            }
        }

    }

    void BaseModelChecker::addInitialState() {
        State initialState{};
        initialState.locations.resize(model.automata.size());
        initialState.variables.resize(model.numVariables);

        for (int i = 0; i < model.automata.size(); ++i) {
            initialState.locations[i] = model.automata[i].initialLocaton;
        }

        addNewState(initialState);
    }

    void BaseModelChecker::printStatistics() {
        printf("Statistics:\n"
               "\texplored: %zu\n"
               "\tgenerated: %zu\n"
               "\tduplicate: %zu\n",
               statistics.exploredCounter, statistics.generatedCounter, statistics.duplicateCounter);
    }

}