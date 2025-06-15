# Description of sample data

_Note: To add all data, run the shell script ./load_sample_data.sh in the root directory_

Presidencies:

2025/01/01- 2025/01/31 President Ranil Wickremesinghe
2025/02/01- 2025/02/28: President Anura Kumara

Ranil

    Organisation:

        2025/01/3:
        - Add Minister of Defence: Sri Lanka Army, Sri Lanka Navy
        - Add Minister of Health: Sri Lanka Medical Council, Department of Ayurveda
        - Add Minister of Finance: General Treasury, Department of National Budget

        2025/01/15
        - Move General Treasury from Minister of Finance to Minister of Defence
        - Terminate Department of Ayurveda in Minister of Health

    People:

        2025/01/6:
        - Minister of Health: George Washington

        2025/01/10:
        - Minister of Finance: Hamilton

        2025/01/25:
        - Minister of Health: Kanye West

Anura

    Organisation:

        2025/02/05:
        - Add Minister of Defence: Sri Lanka Army, Sri Lanka Navy, Sri Lanka Air Force
        - Add Minister of Health: Sri Lanka Medical Council, Department of Ayurveda
        - Add Minister of Finance and Economy: General Treasury, Department of National Budget, Department of Public Finance

        2025/02/18
        - Move Sri Lanka Medical Council to Minister of Defence
        - Terminate Department of National Budget

    People:

        2025/02/12:
        - Minister of Finance and Economy: Vibhatha Abeykoon

        2025/02/14:
        - Minister of Health: Sanjiva Weerawarana
        - Minister of Finance and Economy: Ranil Wickremesinghe

        2025/02/26:
        - Minister of Health: Kanye West

To add all sample data:

For Ranil Wickremesinghe:

```bash
./orgchart -data $(pwd)/data/sample_data/presidents/2025-01-01/ -init -type person
./orgchart -data $(pwd)/data/sample_data/presidents/2025-01-31/ -type person
./orgchart -data $(pwd)/data/sample_data/presidents/2025-02-01/ -type person
./orgchart -data $(pwd)/data/sample_data/presidents/2025-02-28/ -type person

./orgchart -data $(pwd)/data/sample_data/rw/orgchart/2025-01-03 
./orgchart -data $(pwd)/data/sample_data/rw/orgchart/2025-01-15

./orgchart -data $(pwd)/data/sample_data/rw/people/2025-01-06 -type person
./orgchart -data $(pwd)/data/sample_data/rw/people/2025-01-10 -type person
./orgchart -data $(pwd)/data/sample_data/rw/people/2025-01-25 -type person
```

For Anura Kumara:

First to transition from Ranil Wickremesinghe's tenure to Anura Kumara's run the following:

```bash
./orgchart -data $(pwd)/data/sample_data/rw_to_ak/2025-02-01-org
./orgchart -data $(pwd)/data/sample_data/rw_to_ak/2025-02-01-person -type person
```

Now run the following for Anura Kumara's tenure:

```bash
./orgchart -data $(pwd)/data/sample_data/akd/orgchart/2025-02-05/1120_00
./orgchart -data $(pwd)/data/sample_data/akd/orgchart/2025-02-18/

./orgchart -data $(pwd)/data/sample_data/akd/people/2025-02-12/ -type person
./orgchart -data $(pwd)/data/sample_data/akd/people/2025-02-14/ -type person
./orgchart -data $(pwd)/data/sample_data/akd/people/2025-02-26/ -type person
```

This should be the state of the data for each date:

RANIL:

2025/01/01
- President: Ranil Wickremesinghe

2025/01/3:
- Minister of Defence: Sri Lanka Army, Sri Lanka Navy
- Minister of Health: Sri Lanka Medical Council, Department of Ayurveda
- Minister of Finance: General Treasury, Department of National Budget

2025/01/06:
- Minister of Defence: Sri Lanka Army, Sri Lanka Navy
- Minister of Health: Sri Lanka Medical Council, Department of Ayurveda
- Minister of Finance: General Treasury, Department of National Budget
- Minister of Health PERSON: George Washington

2025/01/10:
- Minister of Defence: Sri Lanka Army, Sri Lanka Navy
- Minister of Health: Sri Lanka Medical Council, Department of Ayurveda
- Minister of Finance: General Treasury, Department of National Budget
- Minister of Health PERSON: George Washington
- Minister of Finance PERSON: Hamilton

2025/01/15
- Minister of Defence: Sri Lanka Army, Sri Lanka Navy, General Treasury
- Minister of Health: Sri Lanka Medical Council
- Minister of Finance: Department of National Budget
- Minister of Health PERSON: George Washington
- Minister of Finance PERSON: Hamilton

2025/01/25
- Minister of Defence: Sri Lanka Army, Sri Lanka Navy, General Treasury
- Minister of Health: Sri Lanka Medical Council
- Minister of Finance: Department of National Budget
- Minister of Health PERSON: Kanye West
- Minister of Finance PERSON: Hamilton

ANURA:

2025/02/01
- President: Anura Kumara

2025/02/05
- Minister of Defence: Sri Lanka Army , Sri Lanka Navy, Sri Lanka Air Force
- Minister of Health: Sri Lanka Medical Council, Department of Ayurveda
- Minister of Finance and Economy: Department of National Budget, General Treasury

2025/02/12
- Minister of Defence: Sri Lanka Army , Sri Lanka Navy, Sri Lanka Air Force
- Minister of Health: Sri Lanka Medical Council, Department of Ayurveda
- Minister of Finance and Economy: Department of National Budget, General Treasury
- Minister of Finance and Economy PERSON: Vibhatha Abeykoon

2025/02/14
- Minister of Defence: Sri Lanka Army , Sri Lanka Navy, Sri Lanka Air Force
- Minister of Health: Sri Lanka Medical Council, Department of Ayurveda
- Minister of Finance and Economy: Department of National Budget, General Treasury
- Minister of Finance and Economy PERSON: Ranil Wickremesinghe
- Minister of Health PERSON: Sanjiva Weerawarana

2025/02/18
- Minister of Defence: Sri Lanka Army , Sri Lanka Navy, Sri Lanka Air Force, Sri Lanka Medical Council
- Minister of Health: Department of Ayurveda
- Minister of Finance and Economy: General Treasury
- Minister of Finance and Economy PERSON: Ranil Wickremesinghe
- Minister of Health PERSON: Sanjiva Weerawarana

2025/02/26
- Minister of Defence: Sri Lanka Army , Sri Lanka Navy, Sri Lanka Air Force, Sri Lanka Medical Council
- Minister of Health: Department of Ayurveda
- Minister of Finance and Economy: General Treasury
- Minister of Finance and Economy PERSON: Ranil Wickremesinghe
- Minister of Health PERSON: Kanye West

Cypher Query for checking each date:

```bash
WITH date("2025-01-25") AS targetDate

// Government -> Minister
OPTIONAL MATCH (gov:Organisation {MinorKind: "government"})-[r1:AS_MINISTER]->(min:Organisation {MinorKind: "minister"})
WHERE date(r1.Created) <= targetDate AND (date(r1.Terminated) IS NULL OR date(r1.Terminated) > targetDate)

// Minister -> Department
OPTIONAL MATCH (min)-[r2:AS_DEPARTMENT]->(dep:Organisation {MinorKind: "department"})
WHERE date(r2.Created) <= targetDate AND (date(r2.Terminated) IS NULL OR date(r2.Terminated) > targetDate)

// Government -> Citizen
MATCH (gov)-[r3:AS_APPOINTED]->(cit1:Person {MinorKind: "citizen"})
WHERE date(r3.Created) <= targetDate AND (date(r3.Terminated) IS NULL OR date(r3.Terminated) > targetDate)

// Minister -> Citizen
OPTIONAL MATCH (min)-[r4:AS_APPOINTED]->(cit2:Person {MinorKind: "citizen"})
WHERE date(r4.Created) <= targetDate AND (date(r4.Terminated) IS NULL OR date(r4.Terminated) > targetDate)

RETURN DISTINCT
  gov, r1, min,
  r2, dep,
  r3, cit1,
  r4, cit2

```
 


