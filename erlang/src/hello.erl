-module(modelChecker).

-export([startMC/0, createProcess/0, explore/1, createNode/0, createEdge/0]).

%setup model, create processes and start exploring

startMC() -> 0.

%create process from nodes and edges

createProcess() -> 0.

%create node from name and invariants

createNode() -> 0.

%create edge from starting node, target node, guards and updates

createEdge() -> 0.

%returns all state vectors reachable in one step from the input vector

explore(stateVector) -> io:fwrite("hello, world\n").


