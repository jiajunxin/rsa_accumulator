#!/usr/bin/env bash

go build -o lagrange_test.out ../.
#./lagrange_test.out -h

# 1
./lagrange_test.out -bit=1792 -try=1200