#ifndef CPP_LOCATION_H
#define CPP_LOCATION_H

#include <utility>
#include <vector>

#include "Edge.h"

namespace model {

    struct Invariant : Guard {
        explicit Invariant(std::vector<Predicate> predicates = {}) : Guard{std::move(predicates)} {}
    };

    struct Location {
        Invariant invariant;
        std::vector<Edge> edges;

        explicit Location(Invariant invariant = Invariant{}, std::vector<Edge> edges = {})
            : invariant{std::move(invariant)}, edges{std::move(edges)} {}
    };
}

#endif //CPP_LOCATION_H
