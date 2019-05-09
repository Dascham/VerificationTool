#ifndef CPP_EDGE_H
#define CPP_EDGE_H

#include <vector>

#include "State.h"
#include "Abstract.h"

namespace model {

    struct Guard;
    struct Update;
    struct Sync;

    struct Edge;

    struct Guard{
        std::vector<Predicate> predicates;

        bool operator()(const State& state) const {
            for (auto pred : predicates) {
                if (!pred(state)) {
                    return false;
                }
            }
            return true;
        }

        explicit Guard(std::vector<Predicate> predicates = {}) : predicates{std::move(predicates)} {}
    };

    struct Update {
        std::vector<Assignment> assignments;

        State operator()(State state) const {
            for (auto ass : assignments) {
                state = ass(state);
            }
            return state;
        };

        explicit Update(std::vector<Assignment> assignments = {}) : assignments{std::move(assignments)} {}
    };

    struct Sync {
        enum class Type {
            None,
            Send,
            Recv
        };

        Type type;
        uint8_t channel;

        explicit Sync(Type type = Type::None, uint8_t channel = 0) : type{type}, channel{channel} {}
    };

    struct Edge{
        size_t destination;

        Guard guard;
        Update update;
        Sync sync;

        explicit Edge(size_t destination, Guard guard = Guard{}, Update update = Update{}, Sync sync = Sync{})
            : destination{destination}, guard{std::move(guard)}, update{std::move(update)}, sync{sync} {}
    };

}

#endif //CPP_EDGE_H
