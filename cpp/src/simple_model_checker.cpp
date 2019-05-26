#include <iostream>
#include <queue>
#include <unordered_set>
#include <cassert>

#include "model/Model.h"
#include "modelcheckers/SimpleModelChecker.h"
#include "modelcheckers/DistributedModelChecker.h"

int main() {
    printf("Welcome to the simple_model_checker.\n");

    using namespace model;
    using namespace modelcheckers;

    constexpr size_t numVars = 1;

    SimpleModelChecker simpleModelChecker{THE_MODEL};
    simpleModelChecker.checkModel();

    return 0;
}