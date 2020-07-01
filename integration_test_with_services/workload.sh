#!/usr/bin/env bash

echo Running workload
PGPASSWORD=$TEST_PASSWORD psql -d mytestdb -h localhost -U mytestuser <<_eof
insert into my_table values('test value');
_eof

