#!/usr/bin/env bash

go build -o lagrange_test.out ../.
#./lagrange_test.out -h
./lagrange_test.out -bit=896 -try=100
./lagrange_test.out -bit=896 -try=100
./lagrange_test.out -bit=896 -try=100

./lagrange_test.out -bit=1792 -try=100
./lagrange_test.out -bit=1792 -try=100
./lagrange_test.out -bit=1792 -try=100