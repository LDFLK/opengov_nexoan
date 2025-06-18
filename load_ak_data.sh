#!/bin/bash

# Load Anura's presidency data (termination of Ranil's presidency and starting Anura's)
./orgchart -data $(pwd)/data/people/akd/2024-09-23/ -type person
./orgchart -data $(pwd)/data/orgchart/akd/2024-09-23/

# Load Anura's first cabinet release
./orgchart -data $(pwd)/data/orgchart/akd/2024-09-27/

./orgchart -data $(pwd)/data/people/akd/2024-09-25/ -type person

# Data for 2024-11-25 is not ready yet