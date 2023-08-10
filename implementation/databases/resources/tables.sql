CREATE TABLE IF NOT EXISTS person (
    person_id   SERIAL PRIMARY KEY,
    username    varchar(255)
);

INSERT INTO person VALUES (1, 'resource1');
INSERT INTO person VALUES (2, 'resource2');
INSERT INTO person VALUES (3, 'resource3');

