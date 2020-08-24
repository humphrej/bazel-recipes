#!/usr/bin/env bash

if [ "x$TEST_PASSWORD" == "x" ]; then
  echo Environment variable TEST_PASSWORD must be set
  exit 1
fi

su - postgres -c "psql <<_eof
create database turps;
create user turps with password '$TEST_PASSWORD';
_eof
"
