#!/bin/bash
cd cmake-build-debug/

./worker 0 & ./worker 1 & ./worker 2 & ./worker 3 & ./master