CREATE DATABASE payments;

\connect payments

CREATE EXTENSION "uuid-ossp";

CREATE TABLE payments (
    id      uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    amount  numeric(12, 2) DEFAULT 0
);