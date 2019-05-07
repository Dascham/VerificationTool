#ifndef CPP_ABSTRACT_H
#define CPP_ABSTRACT_H

namespace model {

    struct Term {
        enum class Type {
            Variable, Constant
        };

        Type type;
        int8_t valueOrVariable;

        int8_t operator()(State state) {
            switch (type) {
                case Type::Variable:
                    return state.variables[valueOrVariable];
                case Type::Constant:
                    return valueOrVariable;
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

        Term(Type type, int8_t valueOrVariable) : type{type}, valueOrVariable{valueOrVariable} {}
    };

    struct Assignment {
        enum class AssignOperator{
            Assign, IncAssign
        };
        AssignOperator assignOperator;

        int lhsVariable;
        Term rhsTerm;

        Assignment(int lhsVariable, AssignOperator assignOperator, Term rhsTerm)
            : lhsVariable{lhsVariable}, assignOperator{assignOperator}, rhsTerm{rhsTerm} {}

        State operator()(State state) {

            switch (assignOperator) {
                case AssignOperator::Assign:
                    state.variables[lhsVariable] = rhsTerm(state);
                    break;
                case AssignOperator::IncAssign:
                    state.variables[lhsVariable] += rhsTerm(state);
                    break;
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
            }
        }
    };

}

#endif //CPP_ABSTRACT_H
