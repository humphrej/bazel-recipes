#!/usr/bin/env bash

#MYTESTDB_PASSWORD=mytestuser

PGPASSWORD=$TEST_PASSWORD psql -d mytestdb -h localhost -U mytestuser <<_eof
create table my_table(a varchar not null);
_eof
