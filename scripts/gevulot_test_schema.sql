-- This SQL script (re-)creates test database.

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS companies;

CREATE TABLE companies (
  id   SERIAL       PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE users (
  id         SERIAL       PRIMARY KEY,
  company_id INTEGER      NOT NULL REFERENCES companies(id),
  name       VARCHAR(255) NOT NULL,
  email      VARCHAR(255) NOT NULL
);
