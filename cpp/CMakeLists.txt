cmake_minimum_required(VERSION 3.13)
project(cpp)

set(CMAKE_CXX_STANDARD 17)

include_directories(
        src/
        lib/kissnet
)

#set(THREADS_PREFER_PTHREAD_FLAG ON)
find_package(Threads REQUIRED)

add_executable(cpp
        src/main.cpp
        src/State.h
        src/Block.h
        src/model/Model.h
        src/model/Location.h
        src/model/Edge.h
        src/util.h
        src/State.cpp)


add_executable(simple_model_checker
        src/simple_model_checker.cpp
        src/State.cpp
        src/modelcheckers/BaseModelChecker.cpp
        src/Block.cpp
        )

target_link_libraries(simple_model_checker Threads::Threads)

add_executable(master
        src/master.cpp
        )

add_executable(worker
        src/worker.cpp
        src/State.cpp
        src/modelcheckers/BaseModelChecker.cpp
        src/Block.cpp
        )
target_link_libraries(worker Threads::Threads ws2_32)

add_custom_target(distributed_both)
add_dependencies(distributed_both master worker)
