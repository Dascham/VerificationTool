#include <modelcheckers/DistributedModelChecker.h>

int main() {
    // parse params(worker id)

    size_t workerID = 0;

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