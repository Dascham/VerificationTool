#ifndef CPP_UTIL_H
#define CPP_UTIL_H

#include <functional>

namespace util {

    template<class T>
    inline void hash_combine(size_t &hashVal, const T &v) {
        std::hash<T> hasher{};
        hashVal ^= hasher(v) + 0x9e3779b9'9e3779b9
                + (hashVal << (size_t) 6) + (hashVal >> (size_t) 2)
                + (hashVal << (sizeof(size_t) - 8)) + (hashVal >> (sizeof(size_t) - 8))
                ;
    }


    template<typename T = std::runtime_error>
    inline void _assert(bool success, std::string msg = "") {
        if (!success) throw T(msg);
    }

    template<typename T>
    inline void enforceBoundsAddition(T left, T right, const std::string& prefix = "") {
        // left + right >= min .. left >= min - right
        _assert(left >= std::numeric_limits<T>::min() - right,
                prefix
                + " underflow: " + std::to_string(left + right)
                + " < min:" + std::to_string(std::numeric_limits<T>::min())
                + " (" + typeid(T).name() + ")"
        );
        // left + right <= max .. rhs <= max - right
        _assert(left <= std::numeric_limits<T>::max() - right,
                prefix
                + " overflow: " + std::to_string(left + right)
                + " > max:" + std::to_string(std::numeric_limits<T>::max())
                + " (" + typeid(T).name() + ")"
        );
    }

}

#endif //CPP_UTIL_H
