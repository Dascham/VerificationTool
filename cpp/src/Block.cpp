#include <ostream>
#include <iomanip>
#include <cassert>

#include "State.h"
#include "Block.h"

bool operator==(const Block &a, const Block &b) {
    assert((a.blockID == b.blockID) == (a.type == b.type));
    return a.blockID == b.blockID;

}

std::ostream &operator<<(std::ostream &os, Block const &m) {
    switch (m.type) {
        case Block::Type::Block:
            os << "Block:" << m.blockID << "\n";
        break;
        case Block::Type::Cover:
            os << "Block(Cover):" << m.blockID << "\n";
        break;
    }

    os << std::endl;

    return os;
}