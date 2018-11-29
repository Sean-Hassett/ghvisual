#!bin/bash

cd ghvisual && go build -o main . && cd ..
chmod +x ./ghvisual/main
./ghvisual/main