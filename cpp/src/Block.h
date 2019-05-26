#ifndef CPP_BLOCK_H
#define CPP_BLOCK_H


#include <vector>

#include "State.h"

struct Block {
    enum class Type {
        Block, Cover
    };

    Type type;
    size_t blockID;

    Block(size_t blockID, Type type = Type::Block) : blockID{blockID}, type{type} {} // NOLINT(google-explicit-constructor,hicpp-explicit-conversions)
};

bool operator==(const Block &a, const Block &b);
std::ostream &operator<<(std::ostream &os, Block const &m);


struct Partitioning {
    std::vector<Block> blocks;
};

struct Herd {
    std::vector<Block> groups;
};

#endif //CPP_BLOCK_H
