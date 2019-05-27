#ifndef CPP_PACKET_H
#define CPP_PACKET_H

struct MasterPacket {
    enum class Type : uint8_t {
        Invalid, InitialState, IsDone, IsStillDone, Terminate
    };

    Type type;
    uint8_t data;

    explicit MasterPacket(Type type, uint8_t data = 0) : type{type}, data{data} {}
};

struct WorkerPacket {
    enum class Type : uint8_t {
        Invalid, NotDone, FirstDone, SecondDone
    };

    Type type;
    uint8_t data;

    explicit WorkerPacket(Type type, uint8_t data = 0) : type{type}, data{data} {}
};

#endif //CPP_PACKET_H
