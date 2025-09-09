
CREATE EXTENSION IF NOT EXISTS dblink;

DO
$do$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'pg-summary') THEN
      PERFORM dblink_exec('dbname=postgres user=user password=password', 'CREATE DATABASE "pg-summary"');
   END IF;
END
$do$;
