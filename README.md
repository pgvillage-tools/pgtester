# pgtester
A tool to run queries against postgres and check for an expected result

## The origin
While writing postgres software that manages objects in Postgres (like [pgfga](https://github.com/MannemSolutions/pgfga)), we needed a tool for easy integration testing.
As an integration test we just wanted to create an environment with Postgres, the tool, (and other components as required), run the tool and check the outcome in postgres.
We decided to build a tool which can run defined queries against Postgres, and check for expected results.
And thus [pgtester](https://github.com/MannemSolutions/pgtester) was born.

## Downloading pgtester
The most straight forward way is to download [pgtester](https://github.com/MannemSolutions/pgtester) directly from the [github release page](https://github.com/MannemSolutions/pgtester/releases).
But there are other options, like
- using the [container image from dockerhub](https://hub.docker.com/repository/docker/mannemsolutions/pgtester/general)
- direct build from source (if you feel you must)

Please refer to [our download instructions](DOWNLOAD_AND_RUN.md) for more details on all options.

## Usage
After downloading the binary to a folder in your path, you can run pgtester with a command like:
```bash
pgtester ./mytest*.yml ./andonemoretest.yml
```
Or using stdin:
```bash
cat ./mytests*.yml | pgtester
```

## Defining your tests
A more detailed description can be found in [our test definition guide](TESTS.md).

TLDR; you can define one or more test chapters as yaml documents (separated by the '---' yaml doc separator).
Each test chapter can have the following information defined:
- a dsn, whith all connection details to connect to postgres.
  - **Note** that instead of configuring in this chapter, the [libpq environment variables](https://www.postgresql.org/docs/current/libpq-envars.html) can also be used, but options configured in this chapter take precedence.
- You can set the number of retries, delay and debugging options
- Each test can define
  - a name (defaults to the query when not set),
  - the query
  - the expected result (a list of key/value pairs)
  - the option to reverse the outcome (Ok results are counted as errors and vice versa)

An example test definition could be:
```yaml
---
dsn:
  host: postgres
  port: 5432
  user: postgres
  password: pgtester

retries: 60
delay: 1s
debug: false

tests:
- name: After initialization you normally have 3 databases
  query: "select count(*) total from pg_database"
  results:
  - total: 3
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
