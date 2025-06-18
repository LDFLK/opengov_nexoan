#!/bin/bash

# Load Anura's people and presidency gazettes
./orgchart -data $(pwd)/data/documents/akd/person/ -type document

# Load Anura's org gazettes
./orgchart -data $(pwd)/data/documents/akd/organisation/ -type document

# Load Anura's presidency data (termination of Ranil's presidency and starting Anura's)
./orgchart -data $(pwd)/data/people/akd/2024-09-23/ -type person
./orgchart -data $(pwd)/data/orgchart/akd/2024-09-23/

# Load Anura's org data
./orgchart -data $(pwd)/data/orgchart/akd/2024-09-25/2403-38
./orgchart -data $(pwd)/data/orgchart/akd/2024-09-25/2403-39

./orgchart -data $(pwd)/data/orgchart/akd/2024-09-27/

./orgchart -data $(pwd)/data/orgchart/akd/2024-11-18/2411-09
./orgchart -data $(pwd)/data/orgchart/akd/2024-11-18/2411-10

./orgchart -data $(pwd)/data/orgchart/akd/2024-11-25

# Load Anura's people data
./orgchart -data $(pwd)/data/people/akd/2024-09-25/ -type person

./orgchart -data $(pwd)/data/people/akd/2024-11-18/2411-09/ -type person
./orgchart -data $(pwd)/data/people/akd/2024-11-18/2411-10/ -type person