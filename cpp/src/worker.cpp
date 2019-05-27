#include <string>

#include <modelcheckers/DistributedModelChecker.h>

int main(int argc, char *argv[]) {
    size_t workerID = 0;

    if (argc > 1) {
        int wID = std::stoi(argv[1]);
        if (wID < 0) {
            fprintf(stderr, "workerID(%d) cannot be negative", wID);
            exit(1);
        }
        workerID = wID;
    }
    printf("workerID: %zu\n", workerID);

    using namespace model;
    using namespace modelcheckers;

    constexpr size_t numVars = 1;

    DistributedModelChecker modelChecker{workerID, THE_MODEL};
    modelChecker.checkModel();

    return 0;
}