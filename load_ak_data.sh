#!/bin/bash

# Load Anura's presidency data (termination of Ranil's presidency and starting Anura's)
./orgchart -data $(pwd)/data/orgchart/akd/2024-09-23/
./orgchart -data $(pwd)/data/people/akd/2024-09-23/ -type person