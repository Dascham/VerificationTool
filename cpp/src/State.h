#ifndef CPP_STATE_H
#define CPP_STATE_H


#include <vector>
#include <ostream>

#include "util.h"

struct State {
    std::vector<size_t> locations;
    std::vector<int> variables;
    // Clocks/Zones/DBM go here


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
