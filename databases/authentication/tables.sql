\c auth_db;

CREATE TABLE IF NOT EXISTS Users (
    user_id   SERIAL PRIMARY KEY,
    username    varchar(255),
    passwd_hash    varchar(255),
    groups  varchar(255)[],
    trust   int
);

INSERT INTO 
    Users (user_id, username, passwd_hash, groups, trust)
VALUES
    (1, 'tiko', 'pass1','{sudoers}', '3'),
    (2, 'tico', 'pass2','{web}', '2'),
    (3, 'rico', 'pass3','{testers}', '1');
