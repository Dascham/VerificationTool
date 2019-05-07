#ifndef CPP_GROUP_H
#define CPP_GROUP_H


#include <vector>

#include "State.h"

struct Grouping { // TODO: or Herd
    std::vector<size_t> groups;

    static Grouping fromState(const State &state);
};


#endif //CPP_GROUP_H
