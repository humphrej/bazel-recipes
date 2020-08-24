#!/usr/bin/env bash

#MYTESTDB_PASSWORD=mytestuser

if [ "x$TEST_PASSWORD" == "x" ]; then
  echo Environment variable TEST_PASSWORD must be set
  exit 1
fi

su - postgres -c "psql <<_eof
create database mytestdb;
create user mytestuser with password '$TEST_PASSWORD';
_eof
"
