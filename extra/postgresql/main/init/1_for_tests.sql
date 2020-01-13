-- This script is used for CI tests
-- Real one(1_secret.sql) is not in git index
CREATE USER test WITH SUPERUSER PASSWORD 'testtest';
CREATE DATABASE testbase OWNER test;