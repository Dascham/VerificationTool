#include <ostream>
#include <iomanip>
#include <cassert>

#include "State.h"
#include "Group.h"

bool operator==(const Block &a, const Block &b) {
    assert((a.groupID == b.groupID) == (a.type==b.type));
    return a.groupID == b.groupID;

}

std::ostream &operator<<(std::ostream &os, Block const &m) {
    switch (m.type) {
        case Block::Type::Block:
            os << "Block:" << m.groupID << "\n";
        break;
        case Block::Type::Cover:
            os << "Block(Cover):" << m.groupID << "\n";
        break;
    }

    os << std::endl;

    return os;
}