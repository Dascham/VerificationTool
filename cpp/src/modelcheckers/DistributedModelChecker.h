#ifndef CPP_DISTRIBUTEDMODELCHECKER_H
#define CPP_DISTRIBUTEDMODELCHECKER_H

#include "BaseModelChecker.h"

#include "SocketThread.h"


namespace modelcheckers {

    class DistributedModelChecker : public BaseModelChecker {
        std::vector<SocketThread> socketThreads;

    protected:
        bool addNewState(const State &state) override {

            //if shouldHashLocal(state); need old state? TODO: continue implementing

            return BaseModelChecker::addNewState(state);
        }

    public:
        void checkModel() override {
            addInitialState();

            while (!stateQueue.empty()) {
                State current = stateQueue.front();
                stateQueue.pop();

                int generatedCounter = generateSuccessors(current);
            }

        }

        explicit DistributedModelChecker(model::Model model) : BaseModelChecker{std::move(model)} {}
    };

}

#endif //CPP_DISTRIBUTEDMODELCHECKER_H
