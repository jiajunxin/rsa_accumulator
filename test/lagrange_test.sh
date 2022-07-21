#!/usr/bin/env bash

go build -o lagrange_test.out ../.
#./lagrange_test.out -h

# 1
./lagrange_test.out -bit=896 -try=50
./lagrange_test.out -bit=1792 -try=50
# 2
./lagrange_test.out -bit=896 -try=50
./lagrange_test.out -bit=1792 -try=50
# 3
./lagrange_test.out -bit=896 -try=50
./lagrange_test.out -bit=1792 -try=50
# 4
./lagrange_test.out -bit=896 -try=50
./lagrange_test.out -bit=1792 -try=50
# 5
./lagrange_test.out -bit=896 -try=50
./lagrange_test.out -bit=1792 -try=50
# 6
./lagrange_test.out -bit=896 -try=50
./lagrange_test.out -bit=1792 -try=50