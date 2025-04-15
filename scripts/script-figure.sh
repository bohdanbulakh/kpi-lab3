#!/usr/bin/env bash

for i in {1..10}
do
  curl -d \
    "figure 0.0 0.0
    move 0.${i} 0.${i}
    update" \
  http://localhost:17000
  sleep 1
done