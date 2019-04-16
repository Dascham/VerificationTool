#ifndef CPP_MODEL_H
#define CPP_MODEL_H

#include <vector>
#include <string>

#include "../State.h"
#include "Location.h"
#include "Edge.h"

namespace Model {

    struct Model;

    struct Variable;

    struct Template;
    struct Process;

    struct Model {
        std::vector<std::string> globalVariables;

        std::vector<Template> templates;
        std::vector<Process> processes;

    };

    struct Variable {
        std::string name;
        // TYPE
        // TODO: maybe just array/vector of numbered integer variables(and map for named lookup)
    };

    struct Template {

        std::vector<Location> locations;
        std::vector<Edge> edges;

        std::vector<Variable> localVariables;
    };

    struct Process {
        const Template &aTemplate;

        explicit Process(const Template &aTemplate) : aTemplate(aTemplate) {}
    };


}


#endif //CPP_MODEL_H
