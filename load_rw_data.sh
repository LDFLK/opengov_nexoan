#!/bin/bash

# Load Ranil's presidency
./orgchart -data $(pwd)/data/people/presidents/2022-07-20/ -init -type person

# Load Ranil's Org data
./orgchart -data $(pwd)/data/orgchart/rw/2022-07-22/
./orgchart -data $(pwd)/data/orgchart/rw/2022-09-16/
./orgchart -data $(pwd)/data/orgchart/rw/2022-10-05/

./orgchart -data $(pwd)/data/orgchart/rw/2022-10-26/
./orgchart -data $(pwd)/data/orgchart/rw/2022-12-22/
./orgchart -data $(pwd)/data/orgchart/rw/2023-01-19/

./orgchart -data $(pwd)/data/orgchart/rw/2023-04-27/
./orgchart -data $(pwd)/data/orgchart/rw/2023-05-30/
./orgchart -data $(pwd)/data/orgchart/rw/2023-07-31/

./orgchart -data $(pwd)/data/orgchart/rw/2023-10-23/2355-09/
./orgchart -data $(pwd)/data/orgchart/rw/2023-10-23/2355-10/

./orgchart -data $(pwd)/data/orgchart/rw/2023-12-22/
./orgchart -data $(pwd)/data/orgchart/rw/2024-02-27/

./orgchart -data $(pwd)/data/orgchart/rw/2024-08-23/