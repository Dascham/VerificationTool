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

    SimpleModelChecker simpleModelChecker{Model{
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
    simpleModelChecker.checkModel();

    return 0;
}