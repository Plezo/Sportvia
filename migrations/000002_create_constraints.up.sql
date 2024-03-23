CREATE EXTENSION tsm_system_row;

ALTER TABLE game ADD CONSTRAINT check_attempt_maxAttempts CHECK (attempt <= maxAttempts);

CREATE INDEX playerName_idx ON player USING hash (playerName);