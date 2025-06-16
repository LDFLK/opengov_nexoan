#!/bin/bash

# Load gazettes data
echo "Loading gazettes data..."
./orgchart -data $(pwd)/data/sample_data/gazettes/ -init -type document

# Load Ranil Wickremesinghe's presidency data
echo "Loading Ranil Wickremesinghe's presidency data..."
./orgchart -data $(pwd)/data/sample_data/presidents/2025-01-01/ -type person
./orgchart -data $(pwd)/data/sample_data/presidents/2025-02-01/ -type person

./orgchart -data $(pwd)/data/sample_data/rw/orgchart/2025-01-03 
./orgchart -data $(pwd)/data/sample_data/rw/orgchart/2025-01-15

./orgchart -data $(pwd)/data/sample_data/rw/people/2025-01-06 -type person
./orgchart -data $(pwd)/data/sample_data/rw/people/2025-01-10 -type person
./orgchart -data $(pwd)/data/sample_data/rw/people/2025-01-25 -type person

# Transition to Anura Kumara's presidency
echo "Transitioning to Anura Kumara's presidency..."
./orgchart -data $(pwd)/data/sample_data/rw_to_ak/2025-02-01-org
./orgchart -data $(pwd)/data/sample_data/rw_to_ak/2025-02-01-person -type person

# Load Anura Kumara's presidency data
echo "Loading Anura Kumara's presidency data..."
./orgchart -data $(pwd)/data/sample_data/akd/orgchart/2025-02-05/1120_00
./orgchart -data $(pwd)/data/sample_data/akd/orgchart/2025-02-18/

./orgchart -data $(pwd)/data/sample_data/akd/people/2025-02-12/ -type person
./orgchart -data $(pwd)/data/sample_data/akd/people/2025-02-14/ -type person
./orgchart -data $(pwd)/data/sample_data/akd/people/2025-02-26/ -type person

echo "Sample data loading complete!" 