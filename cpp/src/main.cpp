#include <iostream>

#include "util.h"
#include "State.h"
#include "SocketThread.h"

int main() {
    std::cout << "Hello, World!" << std::endl;

    //size_t hash = 0;
    /*for (int i = 0; i < 100; ++i) {
        util::hash_combine(hash, 1);
        std::cout << hash % 8 << std::endl;
    }*/

    //ThreadSafeQueue<State> temp;
    SocketThread socketThread{};

    std::queue<int> test;

    for (int i = 0; i < 10; ++i) {
        test.emplace(i);
    }

    printf("%zu\n", test.size());

    auto tmp = std::move(test);
    test = decltype(test){};

    printf("%zu\n", test.size());

    return 0;
}