CREATE USER credadmin WITH PASSWORD 'Test1234@' SUPERUSER CREATEDB LOGIN;
\c credstore;
CREATE EXTENSION pgcrypto;
SELECT gen_random_uuid();
CREATE TABLE user_cred(
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   username TEXT NOT NULL CHECK(username <> ''),
   hand_writing TEXT NOT NULL CHECK(hand_writing <> ''),
   pw_encoded TEXT NOT NULL CHECK(pw_encoded <> ''),
   created INTEGER,
   modified INTEGER CHECK(modified >= created),
   deleted BOOLEAN DEFAULT FALSE,
   UNIQUE(id)
);

CREATE UNIQUE INDEX user_cred_active_unique_index
ON user_cred(username, deleted)
WHERE deleted IS FALSE;

GRANT SELECT, INSERT, UPDATE ON TABLE user_cred to credadmin;