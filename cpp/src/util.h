#ifndef CPP_UTIL_H
#define CPP_UTIL_H

#include <functional>

namespace util {

    template<class T>
    inline void hash_combine(size_t &seed, const T &v) {
        std::hash<T> hasher{};
        seed ^= hasher(v) + 0x9e3779b9 + (seed << 6) + (seed >> 2);
    }

}

#endif //CPP_UTIL_H
