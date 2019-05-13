#ifndef CPP_SOCKETTHREAD_H
#define CPP_SOCKETTHREAD_H

#include <thread>
#include <mutex>
#include <queue>
#include <atomic>

class SocketThread {
    std::thread thread;

    std::mutex mutex{};
    std::queue<State> queue{};

    std::atomic<bool> running;

    void run() {

        while (running) {
            // do socket read stuff
        }
    }

public:
    // Used to hand over the queue to the worker thread
    std::queue<State> &&stealQueue() {
        std::unique_lock lock{mutex};
        return std::move(queue);
    }

    void join() {
        if (!thread.joinable()) {
            throw std::logic_error("Tried to terminate thread while it is not joinable");
        }

        // TODO: Close socket

        running = false;
        thread.join();

        // return success status? or assume at this point everything went well
    }

    SocketThread(/*socket? here*/) : running{true} {
        running = true;

        // TODO: open socket?

        thread = std::thread(&SocketThread::run, this);
    }
};

#endif //CPP_SOCKETTHREAD_H
