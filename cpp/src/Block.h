#ifndef CPP_BLOCK_H
#define CPP_BLOCK_H


#include <vector>

#include "State.h"

struct Block {
    enum class Type {
        Block, Cover, Head
    };

    Type type;
    size_t ID;

    bool isCover() const {
        return type == Type::Cover;
    }

    bool isHead() const {
        return type == Type::Head;
    }

    Block(size_t ID, Type type = Type::Block) : ID{ID}, type{type} {} // NOLINT(google-explicit-constructor,hicpp-explicit-conversions)
};

bool operator==(const Block &a, const Block &b);
std::ostream &operator<<(std::ostream &os, const Block &m);




struct Partitioning {
    std::vector<Block> blocks;
};

struct Herd {
    std::vector<Block> groups;
};
bool operator==(const Herd &a, const Herd &b);
std::ostream &operator<<(std::ostream &os, const Herd &m);

namespace std {
    template <>
    struct hash<Herd> {
        size_t operator()(const Herd & x) const {
            size_t hash{};
            for(auto i : x.groups) {
                util::hash_combine(hash, i.ID);
            }

            return hash;
        }
    };
}

#endif //CPP_BLOCK_H
