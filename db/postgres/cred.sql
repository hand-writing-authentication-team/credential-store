CREATE TABLE user_cred(
   id SERIAL UNIQUE PRIMARY KEY,
   username TEXT NOT NULL CHECK(username <> ''),
   hand_writing TEXT NOT NULL CHECK(hand_writing <> ''),
   pw_encoded TEXT NOT NULL CHECK(pw_encoded <> ''),
   created INTEGER,
   modified INTEGER CHECK(modified >= created),
   deleted BOOLEAN,
   UNIQUE(username)
);

GRANT SELECT, INSERT, UPDATE ON TABLE user_cred to credAdmin;