#ifndef CPP_TESTMODELS_H
#define CPP_TESTMODELS_H

#include "model/Model.h"
//#include "config.h"

namespace testmodels {

    model::Model testModel(uint8_t numAutomata = 1, uint8_t numVars = 1) {
        using namespace model;

        return Model{
            numVars,
            {
                Automaton{
                    {
                        Location{Invariant{}, {
                            Edge{0, Guard{{
                                Predicate{Term{Term::Type::Variable, 0}, Predicate::ComparisonOperator::LessThan, Term{Term::Type::Constant, 127}},
                                Predicate{Term{Term::Type::Variable, 0}, Predicate::ComparisonOperator::GreaterThan, Term{Term::Type::Constant, -128}}
                            }}, Update{{
                                Assignment{0, Assignment::AssignOperator::IncAssign, Term{Term::Type::Constant, 1}}
                            }}},
                        }},
                        /*Location{Invariant{}, {
                            Edge{0, Guard{{
                                Predicate{Term{Term::Type::Variable, 0}, Predicate::ComparisonOperator::LessThanEqual, Term{Term::Type::Constant, 127}}
                            }}}
                        }}*/
                    }, 0, Partitioning{{
                        Block{0, Block::Type::Head},
                        //Block{0}
                    }}
                },
                Automaton{
                    {
                        Location{Invariant{}, {
                            Edge{0, Guard{{
                                Predicate{Term{Term::Type::Variable, 1}, Predicate::ComparisonOperator::LessThan, Term{Term::Type::Constant, 127}},
                                Predicate{Term{Term::Type::Variable, 1}, Predicate::ComparisonOperator::GreaterThan, Term{Term::Type::Constant, -128}}
                            }}, Update{{
                                Assignment{1, Assignment::AssignOperator::IncAssign, Term{Term::Type::Constant, 1}}
                            }}},
                        }},
                            /*Location{Invariant{}, {
                                Edge{0, Guard{{
                                    Predicate{Term{Term::Type::Variable, 0}, Predicate::ComparisonOperator::LessThanEqual, Term{Term::Type::Constant, 127}}
                                }}}
                            }}*/
                    }, 0, Partitioning{{
                        Block{0, Block::Type::Head},
                        //Block{0}
                    }}
                },
                Automaton{
                    {
                        Location{Invariant{}, {
                            Edge{0, Guard{{
                                Predicate{Term{Term::Type::Variable, 2}, Predicate::ComparisonOperator::LessThan, Term{Term::Type::Constant, 127}},
                                Predicate{Term{Term::Type::Variable, 2}, Predicate::ComparisonOperator::GreaterThan, Term{Term::Type::Constant, -128}}
                            }}, Update{{
                                Assignment{2, Assignment::AssignOperator::IncAssign, Term{Term::Type::Constant, 1}}
                            }}},
                            }},
                            /*Location{Invariant{}, {
                                Edge{0, Guard{{
                                    Predicate{Term{Term::Type::Variable, 0}, Predicate::ComparisonOperator::LessThanEqual, Term{Term::Type::Constant, 127}}
                                }}}
                            }}*/
                        }, 0, Partitioning{{
                            Block{0, Block::Type::Head},
                            //Block{0}
                    }}
                }
            }
        };
    }

    model::Model badModel(uint8_t numAutomata) {
        using namespace model;

        const uint8_t numVars = numAutomata;

        std::vector<Automaton> automata{};//{3};

        for (uint8_t i = 0; i < numAutomata; ++i) {
            automata.emplace_back(Automaton{
                {
                    Location{Invariant{},{ //0
                        Edge{1, Guard{ //0-1
                            {
                                Predicate{Term{Term::Type::Variable, i}, Predicate::ComparisonOperator::GreaterThan, Term{Term::Type::Constant, 2}}
                            }
                        }},
                        Edge{3, Guard{ //0-3
                            {
                                Predicate{Term{Term::Type::Variable, i}, Predicate::ComparisonOperator::LessThan, Term{Term::Type::Constant, 16}}
                            }
                        }}
                    }},

                    Location{Invariant{},{ //1
                        Edge{2, Guard{{}}, Update{{
                            Assignment{i, Assignment::AssignOperator::IncAssign,Term{Term::Type::Constant, -1}}
                        }}},
                        Edge{3, Guard{{
                            Predicate{Term{Term::Type::Variable, i}, Predicate::ComparisonOperator::LessThan, Term{Term::Type::Constant, 16}}
                        }}, Update{{
                            Assignment{i, Assignment::AssignOperator::IncAssign,Term{Term::Type::Constant, 3}}
                        }}},
                    }},

                    Location{Invariant{},{ //2
                        Edge{4} //2-4
                    }},

                    Location{Invariant{}, { //3
                        Edge{4}, //3-4
                        Edge{6, Guard{}, Update{{ //3-6
                            Assignment{i, Assignment::AssignOperator::IncAssign,Term{Term::Type::Constant, 2}}
                        }}}
                    }},

                    Location{Invariant{},{ //4
                        Edge{5} //4-5
                    }},

                    Location{Invariant{},{ //5
                        Edge{6, Guard{{ //5-6
                            Predicate{Term{Term::Type::Variable, i}, Predicate::ComparisonOperator::LessThan, Term{Term::Type::Constant, 16}}
                        }}},
                        Edge{0, Guard{{ //5-0
                            Predicate{Term{Term::Type::Variable, i}, Predicate::ComparisonOperator::GreaterThan, Term{Term::Type::Constant, 5}}
                        }}, Update{{
                            Assignment{i, Assignment::AssignOperator::IncAssign,Term{Term::Type::Constant, -3}}
                        }}},
                    }},

                    Location{Invariant{}, { //6
                        Edge{0, Guard{}, Update{{ //6-0
                            Assignment{i, Assignment::AssignOperator::IncAssign,Term{Term::Type::Constant, 1}}
                        }}}
                    }},

                },
                0, {{
                    Block{0, Block::Type::Cover}, //0
                    Block{1, Block::Type::Cover}, //1
                    Block{2, Block::Type::Head}, //2
                    Block{3, Block::Type::Head}, //3
                    Block{4, Block::Type::Head}, //4
                    Block{4, Block::Type::Block}, //5
                    Block{5, Block::Type::Head}, //6
                }}
            });
        }

        return Model{numVars, automata};
    }


    model::Model goodModel(uint8_t numAutomata, bool hasSync = false) {
        using namespace model;

        const uint8_t numVars = numAutomata + (hasSync ? 1 : 0);
        const uint8_t syncVar = numAutomata; // which

        std::vector<Automaton> automata{};

        for (uint8_t i = 0; i < numAutomata; ++i) {
            automata.emplace_back(Automaton{{
                    Location{Invariant{},{ //0
                        Edge{1, Guard{ //0-1
                            {
                                Predicate{Term{Term::Type::Variable, i}, Predicate::ComparisonOperator::LessThan, Term{Term::Type::Constant, 16}}
                            }
                        }, Update{}, (!hasSync) ? Sync{} : Sync{Sync::Type::Recv, 0}},
                        Edge{3, Guard{ //0-3
                            {
                                Predicate{Term{Term::Type::Variable, i}, Predicate::ComparisonOperator::GreaterThan, Term{Term::Type::Constant, 2}}
                            }
                        }}
                    }},

                    Location{Invariant{},{ //1 (right)
                        Edge{2, Guard{{}}, Update{{ //1-2
                            Assignment{i, Assignment::AssignOperator::IncAssign,Term{Term::Type::Constant, 2}}
                        }}}
                    }},

                    Location{Invariant{},{ //2
                        Edge{5, Guard{{}}, (!hasSync) ? Update{} : Update{{ //2-5
                            Assignment{syncVar, Assignment::AssignOperator::IncAssign,Term{Term::Type::Constant, -1}}
                        }}}
                    }},

                    Location{Invariant{},{ //3 (left)
                        Edge{4, Guard{{}}, Update{{ //3-4
                            Assignment{i, Assignment::AssignOperator::IncAssign,Term{Term::Type::Constant, -1}}
                        }}}
                    }},

                    Location{Invariant{},{ //4
                        Edge{5, Guard{{}}, Update{{ //4-5
                            Assignment{i, Assignment::AssignOperator::IncAssign,Term{Term::Type::Constant, -1}}
                        }}}
                    }},

                    Location{Invariant{},{ //5 (join)
                            Edge{6, (!hasSync) ? Guard{} :Guard{{ //5-6
                                Predicate{Term{Term::Type::Variable, syncVar}, Predicate::ComparisonOperator::Equal, Term{Term::Type::Constant, 0}}
                            }}}
                    }},

                    Location{Invariant{},{ //6
                            Edge{7, Guard{}, Update{{ //6-7
                                Assignment{i, Assignment::AssignOperator::IncAssign,Term{Term::Type::Constant, -1}}
                            }}}
                    }},

                    Location{Invariant{},{ //7
                            Edge{0} // 7-0
                    }},

                },0, {{
                    Block{0, Block::Type::Cover}, //0
                    Block{1, Block::Type::Head}, //1
                    Block{1, Block::Type::Block}, //2
                    Block{2, Block::Type::Head}, //3
                    Block{2, Block::Type::Block}, //4
                    Block{3, Block::Type::Head}, //5
                    Block{3, Block::Type::Block}, //6
                    Block{3, Block::Type::Block}, //7
                }}
            });

        }

        if (hasSync) {
            automata.emplace_back(Automaton{{
                Location{Invariant{}, { //0
                        Edge{1, Guard{{ //0-1
                            Predicate{Term{Term::Type::Variable, syncVar},Predicate::ComparisonOperator::LessThanEqual,Term{Term::Type::Constant, 2}}
                        }}, Update{{
                            Assignment{syncVar, Assignment::AssignOperator::IncAssign,Term{Term::Type::Constant, -1}}
                        }}, Sync{Sync::Type::Send, 0}}},
                },
                 Location{Invariant{}, { //1
                         Edge{2}}, //1-2
                 },

                  Location{Invariant{}, { //2
                         Edge{0}}, //2-0
                 },
            },0, {{
                Block{0, Block::Type::Head}, //0
                Block{0, Block::Type::Block}, //1
                Block{0, Block::Type::Block}, //2
            }}});
        }

        return Model{numVars, automata};
    }

}

#endif //CPP_TESTMODELS_H
