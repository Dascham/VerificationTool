#ifndef CPP_ABSTRACT_H
#define CPP_ABSTRACT_H

#include <cassert>

namespace model {

    struct Term {
        enum class Type {
            Variable, Constant
        };

        Type type;
        int16_t valueOrVariable;

        size_t operator()(State state) {
            switch (type) {
                case Type::Variable:
                    assert((size_t)valueOrVariable < state.variables.size());
                    return state.variables[valueOrVariable];
                case Type::Constant:
                    return valueOrVariable;
                default:
                    throw std::domain_error("Invalid Type in Term");
            }
        }

        // Not using Term for variable assignments, so might not need below function
        // Maybe could be done better with a Variable type and return type overloading operator()
        // ... or overload on parameters(no state gives variable location instead)
        /*size_t variable() {
            switch (type) {
                case Type::Variable:
                    return valueOrVariable;
                case Type::Constant:
                    throw std::domain_error("Tried to use a constant as if it were a variable.");
            }
        }*/

        Term(Type type, int16_t valueOrVariable) : type{type}, valueOrVariable{valueOrVariable} {
            switch (type) {
                case Type::Variable:
                    assert(valueOrVariable >= std::numeric_limits<uint8_t>::min());
                    assert(valueOrVariable <= std::numeric_limits<uint8_t>::max());
                    break;
                case Type::Constant:
                    assert(valueOrVariable >= std::numeric_limits<int8_t>::min());
                    assert(valueOrVariable <= std::numeric_limits<int8_t>::max());
                    break;
                default:
                    throw std::domain_error("Invalid Type in Term");
            }
        }
    };

    struct Assignment {
        enum class AssignOperator{
            Assign, IncAssign
        };
        AssignOperator assignOperator;

        uint8_t lhsVariable;
        Term rhsTerm;

        Assignment(uint8_t lhsVariable, AssignOperator assignOperator, Term rhsTerm)
            : lhsVariable{lhsVariable}, assignOperator{assignOperator}, rhsTerm{rhsTerm} {}

        State operator()(State state) {
            const int8_t rhsVal = rhsTerm(state);

            assert(lhsVariable < state.variables.size());

            assert(rhsVal >= std::numeric_limits<int8_t>::min());
            assert(rhsVal <= std::numeric_limits<int8_t>::max());

            switch (assignOperator) {
                case AssignOperator::Assign:
                    state.variables[lhsVariable] = rhsVal;
                    break;
                case AssignOperator::IncAssign: {
                    const int8_t curVal = state.variables[lhsVariable];
                    const int8_t newVal = rhsVal + state.variables[lhsVariable];

                    util::enforceBoundsAddition(curVal, rhsVal,
                                                "Increment variable[" + std::to_string(lhsVariable) + "]");

                    state.variables[lhsVariable] = newVal;
                }
                    break;
                default:
                    throw std::domain_error("Invalid AssignmentOperator in Assignment");
            }

            return state;
        }
    };

    struct Predicate {

        enum class ComparisonOperator : uint8_t {
            Equal, NotEqual,
            LessThan, LessThanEqual,
            GreaterThan, GreaterThanEqual,
        };

        ComparisonOperator comparisonOperator;

        Term lhsTerm;
        Term rhsTerm;

        Predicate(Term lhsTerm, ComparisonOperator comparisonOperator, Term rhsTerm)
                : lhsTerm{lhsTerm}, comparisonOperator{comparisonOperator}, rhsTerm{rhsTerm} {}

        bool operator()(const State& state) {
            switch (comparisonOperator) {
                case ComparisonOperator::Equal:
                    return lhsTerm(state) == rhsTerm(state);
                case ComparisonOperator::NotEqual:
                    return lhsTerm(state) != rhsTerm(state);
                case ComparisonOperator::LessThan:
                    return lhsTerm(state) < rhsTerm(state);
                case ComparisonOperator::LessThanEqual:
                    return lhsTerm(state) <= rhsTerm(state);
                case ComparisonOperator::GreaterThan:
                    return lhsTerm(state) > rhsTerm(state);
                case ComparisonOperator::GreaterThanEqual:
                    return lhsTerm(state) >= rhsTerm(state);
                default:
                    throw std::domain_error("Invalid ComparisonOperator in Predicate");
            }
        }
    };

}

#endif //CPP_ABSTRACT_H
