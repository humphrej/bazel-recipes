CREATE TABLE IF NOT EXISTS change_list (
  id VARCHAR PRIMARY KEY,
  timestamp TIMESTAMPTZ NOT NULL,
  test_ids varchar[] NULL
);


CREATE TABLE IF NOT EXISTS test_run (
  id VARCHAR PRIMARY KEY,
  change_list_id VARCHAR NOT NULL REFERENCES change_list(id),
  output_url VARCHAR NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL,
  UNIQUE (id)
);

CREATE TABLE IF NOT EXISTS test_run_result (
  run_id VARCHAR NOT NULL REFERENCES test_run(id),
  test_id VARCHAR NOT NULL,
  count_runs INT NOT NULL,
  count_fails INT NOT NULL,
  UNIQUE (run_id, test_id)
);


