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

            size_t exploredCounter = 0;
            size_t generatedCounter = 1; // Start with initial state "generated"

            while (!stateQueue.empty()) {
                State current = stateQueue.front();
                stateQueue.pop();

                exploredCounter++;
                generatedCounter += generateSuccessors(current);
            }

            std::cout << "Explored states: " << exploredCounter << std::endl;
            std::cout << "Generated states: " << generatedCounter << std::endl;
            std::cout << "Duplicate states: " << generatedCounter - exploredCounter << std::endl;
        }

        explicit SimpleModelChecker(model::Model model) : BaseModelChecker{std::move(model)} {}
    };

}

#endif //CPP_SIMPLEMODELCHECKER_H
