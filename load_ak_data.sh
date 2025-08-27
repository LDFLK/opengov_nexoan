#!/bin/bash

# Load Anura's people and presidency gazettes
./orgchart -data "$(pwd)/data/documents/Anura Kumara Dissanayake/person/" -type document

# Load Anura's org gazettes
./orgchart -data "$(pwd)/data/documents/Anura Kumara Dissanayake/organisation/" -type document

# Load Anura's presidency data (termination of Ranil's presidency and starting Anura's)
./orgchart -data "$(pwd)/data/people/Anura Kumara Dissanayake/2024-09-23/2403-03_1" -type person
./orgchart -data "$(pwd)/data/orgchart/Anura Kumara Dissanayake/2024-09-23/"
./orgchart -data "$(pwd)/data/people/Anura Kumara Dissanayake/2024-09-23/2403-03_2" -type person

# Load Anura's org data
# ./orgchart -data "$(pwd)/data/orgchart/Anura Kumara Dissanayake/2024-09-25/2403-38"
# ./orgchart -data "$(pwd)/data/orgchart/Anura Kumara Dissanayake/2024-09-25/2403-39"


# terminate all the old Ranil's portfolios which were transferred to Anura and all the people assigned (all Anura)
./orgchart -data "$(pwd)/data/orgchart/Anura Kumara Dissanayake/2024-09-25/2403-38-1"
./orgchart -data "$(pwd)/data/people/Anura Kumara Dissanayake/2024-09-25/2403-38-1" -type person

# Add Anura's new cabinet
./orgchart -data "$(pwd)/data/orgchart/Anura Kumara Dissanayake/2024-09-25/2403-38"
./orgchart -data "$(pwd)/data/orgchart/Anura Kumara Dissanayake/2024-09-25/2403-39"

# Add Anura's new cabinet people
./orgchart -data "$(pwd)/data/people/Anura Kumara Dissanayake/2024-09-25/2403-38" -type person
./orgchart -data "$(pwd)/data/people/Anura Kumara Dissanayake/2024-09-25/2403-39" -type person


./orgchart -data "$(pwd)/data/orgchart/Anura Kumara Dissanayake/2024-09-27/"

# load the rest of Anura's org data
./orgchart -data "$(pwd)/data/orgchart/Anura Kumara Dissanayake/2024-11-18/2411-09"
./orgchart -data "$(pwd)/data/orgchart/Anura Kumara Dissanayake/2024-11-18/2411-10"


# Load Anura's people data
./orgchart -data "$(pwd)/data/people/Anura Kumara Dissanayake/2024-11-18/2411-09/" -type person
./orgchart -data "$(pwd)/data/people/Anura Kumara Dissanayake/2024-11-18/2411-10/" -type person

./orgchart -data "$(pwd)/data/orgchart/Anura Kumara Dissanayake/2024-11-25"

