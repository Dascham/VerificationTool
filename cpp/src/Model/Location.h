#ifndef CPP_LOCATION_H
#define CPP_LOCATION_H


#include <vector>

#include "Edge.h"

namespace Model {

    struct Location{
        /*struct Invariant {
            //Variable
            //Comparison Operator
            //Constant
            //  What about comparing two variables

            virtual bool isSatisfied() const;
        };*/

        std::vector<Edge> edges;
    };
}

#endif //CPP_LOCATION_H
