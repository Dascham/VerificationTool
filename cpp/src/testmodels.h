#ifndef CPP_TESTMODELS_H
#define CPP_TESTMODELS_H

#include "model/Model.h"

namespace testmodels {

    model::Model testModel() {

        using namespace model;

        constexpr size_t numVars = 1;

        return Model{
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
        };

    }

}

#endif //CPP_TESTMODELS_H
