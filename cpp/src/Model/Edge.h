#ifndef CPP_EDGE_H
#define CPP_EDGE_H

#include <vector>

#include "State.h"

namespace Model {

    struct Edge{
        struct Guard{
            // Left hand
            // Comparison Operator
            // Right hand

            //std::vector<Predicate> predicates;

            bool isSatisfied() const {

                /*for (auto it : predicates) {
                    if (!it.isSatisfied()) {
                        return false;
                    }
                }*/

                return true;
            }
        };
        struct Update {
            // Variable
            // Assignment operator
            // Value

            State perform(State current) const {

            };
        };

        Guard guard;
        Update update;
        //Sync sync;

        size_t destination;
    };

}

#endif //CPP_EDGE_H
