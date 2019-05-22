#ifndef CPP_STATE_H
#define CPP_STATE_H


#include <array>
#include <ostream>

#include "util.h"

struct State {
    std::vector<uint8_t> locations;
    std::vector<int8_t> variables;
    // Clocks/Zones/DBM go here

    State() = default;
    State(std::vector<uint8_t> locations, std::vector<int8_t> variables) : locations{std::move(locations)}, variables{std::move(variables)} {}
};

bool operator==(const State &a, const State &b);
std::ostream &operator<<(std::ostream &os, State const &m);

namespace std {
    template <>
    struct hash<State> {
        size_t operator()(const State & x) const {
            size_t hash{};
            for(auto i : x.locations) {
                util::hash_combine(hash, i);
            }

            for(auto i : x.variables) {
                util::hash_combine(hash, i);
            }

            return hash;
        }
    };
}

#endif //CPP_STATE_H
