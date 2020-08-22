#!/usr/bin/env bash

yell() { echo "$0: $*" >&2; }
die() { yell "$*"; exit 111; }
try() { "$@" || die "cannot $*"; }

export DATABASE_URL

echo Running postgres data access layer tests
try ./go_postgres_storage_test

echo Running turps acceptance tests
try ./go_acceptance_test
