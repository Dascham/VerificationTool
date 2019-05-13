#include <iostream>
#include <queue>
#include <unordered_set>
#include <cassert>

#include "model/Model.h"
#include "modelcheckers/SimpleModelChecker.h"
#include "modelcheckers/DistributedModelChecker.h"

int main() {
    std::cout << "Welcome. Welcome to the simple_model_checker.\n"
                 "\n"
                 "You have chosen, or been chosen, to use one of our finest non-distributed model checkers.\n"
              << std::endl;

    using namespace model;
    using namespace modelcheckers;

    constexpr size_t numVars = 1;

    DistributedModelChecker simpleModelChecker{Model{
        numVars,
        {
            Automaton{{
                Location{Invariant{}, {
                    Edge{1, Guard{{
                        Predicate{Term{Term::Type::Variable, 0}, Predicate::ComparisonOperator::LessThan, Term{Term::Type::Constant, 127}}
                    }}, Update{{
                        Assignment{0, Assignment::AssignOperator::IncAssign, Term{Term::Type::Constant, 1}}
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
    simpleModelChecker.checkModel();

    return 0;
}