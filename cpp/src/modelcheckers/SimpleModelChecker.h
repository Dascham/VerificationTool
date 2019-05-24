#ifndef CPP_SIMPLEMODELCHECKER_H
#define CPP_SIMPLEMODELCHECKER_H

#include "BaseModelChecker.h"

#include "util.h"
#include "State.h"
#include "model/Model.h"


namespace modelcheckers {

    class SimpleModelChecker : BaseModelChecker {

    public:
        void checkModel() override {
            addInitialState();

            while (!stateQueue.empty()) {
                State current = stateQueue.front();
                stateQueue.pop();

                generateSuccessors(current);
            }

            printStatistics();
            assert(statistics.duplicateCounter == statistics.generatedCounter-statistics.exploredCounter);
        }

        explicit SimpleModelChecker(model::Model model) : BaseModelChecker{std::move(model)} {}
    };

}

#endif //CPP_SIMPLEMODELCHECKER_H
