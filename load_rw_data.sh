#!/bin/bash

# Load Ranil's people and presidency gazettes
./orgchart -data "$(pwd)/data/documents/Ranil Wickremesinghe/person/" -type document

# Load Ranil's org gazettes
./orgchart -data "$(pwd)/data/documents/Ranil Wickremesinghe/organisation/" -type document

# Load Ranil's presidency data
./orgchart -data "$(pwd)/data/people/Ranil Wickremesinghe/2022-07-20/2289-34-1/" -type person

./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2022-07-20/"

./orgchart -data "$(pwd)/data/people/Ranil Wickremesinghe/2022-07-20/2289-34-2/" -type person

# # Load Ranil's org data - terminate old ministers and departments and add new ministers
./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2022-07-22/"

# ./orgchart -data "$(pwd)/data/people/Ranil Wickremesinghe/2022-07-26/" -type person
# ./orgchart -data "$(pwd)/data/people/Ranil Wickremesinghe/2022-08-04/" -type person

# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2022-09-16/"
# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2022-10-05/"
# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2022-10-26/"

# ./orgchart -data "$(pwd)/data/people/Ranil Wickremesinghe/2022-11-04/" -type person

# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2022-12-22/"


# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2023-01-19/"
# ./orgchart -data "$(pwd)/data/people/Ranil Wickremesinghe/2023-01-19/" -type person

# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2023-04-27/"
# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2023-05-30/"
# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2023-07-31/"

# ./orgchart -data "$(pwd)/data/people/Ranil Wickremesinghe/2023-10-12/" -type person

# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2023-10-23/2355-09/"
# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2023-10-23/2355-10/"

# ./orgchart -data "$(pwd)/data/people/Ranil Wickremesinghe/2023-10-23/" -type person

# ./orgchart -data "$(pwd)/data/people/Ranil Wickremesinghe/2023-12-01/" -type person

# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2023-12-22/"
# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2024-02-27/"
# ./orgchart -data "$(pwd)/data/orgchart/Ranil Wickremesinghe/2024-08-23/"





