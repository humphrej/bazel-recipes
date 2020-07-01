#!/usr/bin/env bash

echo "local mytestdb  all   md5" >> /etc/postgresql/12/main/pg_hba.conf

service postgresql start
