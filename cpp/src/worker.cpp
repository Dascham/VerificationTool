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

    DistributedModelChecker modelChecker{workerID, Model{
        numVars,
        {
            Automaton{{
                Location{Invariant{}, {
                    Edge{1, Guard{{
                        Predicate{Term{Term::Type::Variable, 0}, Predicate::ComparisonOperator::LessThan, Term{Term::Type::Constant, 127}},
                        Predicate{Term{Term::Type::Variable, 0}, Predicate::ComparisonOperator::GreaterThan, Term{Term::Type::Constant, -128}}
                        }}, Update{{
                            Assignment{0, Assignment::AssignOperator::IncAssign, Term{Term::Type::Constant, -1}}
                        }}},
                }},
                Location{Invariant{}, {
                    Edge{0, Guard{{
                        Predicate{Term{Term::Type::Variable, 0}, Predicate::ComparisonOperator::LessThanEqual, Term{Term::Type::Constant, 127}}
                    }}}
                }}
            }}
        }
    }};
    modelChecker.checkModel();

    // TODO: ASSERT model loc,var sizes == constants, do first?


    return 0;
}