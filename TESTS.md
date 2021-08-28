# Tests

## Definition

### Test chapters

Tests are grouped in test chapters.
All tests grouped in a chapter are run with the same postgres connection, but different connections can be defined for different test chapters. 
Each test chapter is defined as a yaml document. You can combine multiple chapters (separated by the '---' yaml doc separator) in one file or stream multiple test chapters into stdin.

### Details in a test chapter
Each test chapter can have the following information defined:
- a dsn, with all connection details to connect to postgres.
 - **Note** that instead of configuring in this chapter, the [environment variables](https://www.postgresql.org/docs/current/libpq-envars.html) can also be used.
 - Options configured in this chapter take precedence over environment variables
 - the dsn chapter allows to set different connection details for different test chapters.
- You can set the number of connection retries with the option `retries`.
 - Default is 0 retries (only one try)
 - The number of retries occur for every test in the chapter!!!
 - in docker-compose environments this (together with the `delay` option) can be useful to wait until postgres is available.
- you can set a delay for each retry with the `delay` option
- you can get more verbose output with the `debug` option
 - this option is in parity with the `-d` commandline option
 - the `debug` option can be set per test chapter and the `-d` commandline option is a global debug option (all chapters)
- Tests are defined as a list of maps

### Details per test
Each test can define
- a name, which defaults to the query when not set
- the query, which needs to be set
- the expected result which is a list of maps (associative arrays)
 - note that the check is very straight forward implemented and as such order matters.
  - always order the list of results in the expected order returned by postgres
  - always use the `ORDER BY` clause when more than one row is expected
 - every row is represented as a map (associative array) where the name of the column is the key and the value in the row is the value in the map (e.a. col: value)
 - **note** that pgtester converts the values as returned by postgres to strings and compares strings
- every test can have an option `reverse` to inverse the outcome
 - so that
  - queries that return the expected results are counted as errors
  - queries that do not return the expected results are counted as OK
 - this is mostly useful when checking for connection errors, but might also be used to define results that you don't want to see returned (when all else would be OK

## Example
Below you will see an example that could be used as input.
The yaml defines 2 test chapters with different connection options for a different database and user.
The second test in the second chapter:
- is missing a name on the second test (and this the query is used as name instead)
- has the revers option set (and since pg_databases usually does not exist, this would be counted as an OK).

```yaml
---
dsn:
  host: postgres
  port: 5432
  user: postgres
  password: pgtester
  dbname: postgres

retries: 5
delay: 1s
debug: false

tests:
- name: After initialization you normally have 3 databases
  query: "select count(*) total from pg_database"
  results:
  - total: 3
---
dsn:
  host: postgres
  port: 5432
  user: testuser
  password: testing123
  dbname: test

tests:
- name: After initialization you normally have the databases postgres, template0 and template1
  query: "select datname from pg_database order by 1"
  results:
  - datname: postgres
  - datname: template0
  - datname: template1
# This test would have name "select datname from pg_databases"
- query: "select datname from pg_databases"
  results: []
  reverse: true
```

## Example output:
If pgtester is downloaded to a folder in your path, and the above example is written to a file test.yaml, the below command would bring a result as listed:
```bash
# cat test.yaml | pgtester
2021-08-27T20:40:32+02:00	INFO	==============================
2021-08-27T20:40:32+02:00	INFO	Running tests from (stdin) (0)
2021-08-27T20:40:32+02:00	INFO	==============================
2021-08-27T20:40:32+02:00	DEBUG	connecting to host='postgres' port='5432' user='postgres' password='pgtester' dbname='postgres'
2021-08-27T20:40:32+02:00	DEBUG	succesfully connected
2021-08-27T20:40:32+02:00	DEBUG	running query select count(*) total from pg_database with arguments []
2021-08-27T20:40:32+02:00	INFO	success as expected on test 'After initialization you normally have 3 databases'
2021-08-27T20:40:32+02:00	INFO	==============================
2021-08-27T20:40:32+02:00	INFO	Running tests from (stdin) (1)
2021-08-27T20:40:32+02:00	INFO	==============================
2021-08-27T20:40:32+02:00	INFO	success as expected on test 'After initialization you normally have the databases postgres, template0 and template1'
2021-08-27T20:40:32+02:00	ERROR	expected error occurred on test 'select datname from pg_databases': ERROR: relation "pg_databases" does not exist (SQLSTATE 42P01)
2021-08-27T20:40:32+02:00	INFO	===============================================
2021-08-27T20:40:32+02:00	INFO	succesfully finished without unexpected results
2021-08-27T20:40:32+02:00	INFO	===============================================
```
Please **note** that:
- the debug output shows some more details
- although the reverse option is set for the last test, it is still recorded with an ERROR, but the run ends with `succesfully finished without unexpected results`.
