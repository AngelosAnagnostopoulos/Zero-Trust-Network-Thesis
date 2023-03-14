CREATE TABLE IF NOT EXISTS person (
    person_id   SERIAL PRIMARY KEY,
    username    varchar(255)
);

INSERT INTO person VALUES (1, 'mike');
INSERT INTO person VALUES (2, 'sike');
INSERT INTO person VALUES (3, 'like');

