#include <ostream>
#include <iomanip>
#include <cassert>

#include "State.h"
#include "Block.h"

bool operator==(const Block &a, const Block &b) {
    assert((a.ID == b.ID) == (a.type == b.type));
    return a.ID == b.ID;

}

std::ostream &operator<<(std::ostream &os, Block const &m) {
    switch (m.type) {
        case Block::Type::Block:
            os << "Block:" << m.ID << "\n";
        break;
        case Block::Type::Cover:
            os << "Block(Cover):" << m.ID << "\n";
        break;
    }

    os << std::endl;

    return os;
}

bool operator==(const Herd &a, const Herd &b) {

    if (a.groups.size() != b.groups.size()) {
        return false;
    }

    for (int i = 0; i < a.groups.size(); ++i) {
        if (a.groups[i].ID != b.groups[i].ID) return false;
    }

    return false;
}

std::ostream &operator<<(std::ostream &os, Herd const &m) {
    os << "Herd{ ";

    for (const auto &item : m.groups) {
        os << item.ID << "("
        << ((item.type == Block::Type::Cover) ? "cover" : (item.type == Block::Type::Head) ? "head" : "block")
        << ")";
    }
    printf("}\n");

    os << "}" << std::endl;

    return os;
}