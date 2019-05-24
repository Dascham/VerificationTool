#include <ostream>
#include <iomanip>

#include "State.h"

bool operator==(const State &a, const State &b) {
    if (a.locations.size() != b.locations.size()) return false;
    if (a.variables.size() != b.variables.size()) return false;
    for (int i = 0; i < a.locations.size(); ++i) {
        if(a.locations[i] != b.locations[i]) return false;
    }
    for (int i = 0; i < a.variables.size(); ++i) {
        if(a.variables[i] != b.variables[i]) return false;
    }

    return true;
}

std::ostream &operator<<(std::ostream &os, State const &m) {
    std::ios::fmtflags os_flags (os.flags());

    os << "State {\n"
          "  locations: < "
            ;

    for (auto loc : m.locations) {
        os << std::setw(2) << (int)loc << " ";
    }

    os << ">\n"
          "  variables: < "
            ;

    for (auto var : m.variables) {
        os << std::setw(4) << (int)var << " ";
    }
    os << ">\n"
          "}\n"
       << std::endl;

    os.flags(os_flags);
    return os;
}