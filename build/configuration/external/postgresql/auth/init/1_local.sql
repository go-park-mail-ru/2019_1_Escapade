-- This script is used for CI tests
-- Real one(1_secret.sql) is not in git index
CREATE USER authlocal WITH SUPERUSER PASSWORD 'authlocal';
CREATE DATABASE authbase OWNER authlocal;