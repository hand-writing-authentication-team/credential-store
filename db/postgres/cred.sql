\c credstore;
CREATE EXTENSION pgcrypto;
SELECT gen_random_uuid();
CREATE TABLE user_cred(
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   username TEXT NOT NULL CHECK(username <> ''),
   hand_writing TEXT NOT NULL CHECK(hand_writing <> ''),
   pw_encoded TEXT NOT NULL CHECK(pw_encoded <> ''),
   race TEXT,
   created INTEGER,
   modified INTEGER CHECK(modified >= created),
   deleted BOOLEAN DEFAULT FALSE,
   UNIQUE(id)
);

CREATE UNIQUE INDEX user_cred_active_unique_index
ON user_cred(username, deleted)
WHERE deleted IS FALSE;

CREATE TABLE validate_handwriting (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   userId UUID,
   username TEXT NOT NULL CHECK(username <> ''),
   hand_writing TEXT NOT NULL CHECK(hand_writing <> ''),
   created INTEGER,
   modified INTEGER CHECK(modified >= created),
   deleted BOOLEAN DEFAULT FALSE,
   UNIQUE(id),
   FOREIGN KEY (userId) REFERENCES user_cred (id)
);

CREATE UNIQUE INDEX validate_handwriting_active_unique_index
ON validate_handwriting(username, deleted)
WHERE deleted IS FALSE;

GRANT SELECT, INSERT, UPDATE ON TABLE user_cred to credadmin;
GRANT SELECT, INSERT, UPDATE ON TABLE validate_handwriting to credadmin;