CREATE TABLE IF NOT EXISTS Users (
    user_id   SERIAL PRIMARY KEY,
    username    varchar(255),
    passwd_hash    varchar(255),
    lemao   varchar(255)
);

INSERT INTO Users VALUES (1, 'tiko', 'testpass1','123');
INSERT INTO Users VALUES (2, 'tico', 'testpass2','123');
INSERT INTO Users VALUES (3, 'rico', 'testpass3','123');
