#include <iostream>
#include "util.h"

int main() {
    std::cout << "Hello, World!" << std::endl;

    size_t hash = 0;

    for (int i = 0; i < 100; ++i) {
        util::hash_combine(hash, 1);
        std::cout << hash % 8 << std::endl;
    }

    return 0;
}