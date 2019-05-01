cmake_minimum_required(VERSION 3.13)
project(cpp)

set(CMAKE_CXX_STANDARD 14)

include_directories(
        src/
)

add_executable(cpp
        src/main.cpp
        src/State.h
        src/Group.h
        src/Model/Model.h
        src/Model/Location.h
        src/Model/Edge.h
        src/util.h
        src/State.cpp)

add_executable(simple_model_checker
        src/simple_model_checker.cpp
        src/State.cpp)