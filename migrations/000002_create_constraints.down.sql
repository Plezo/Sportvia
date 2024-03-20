DROP EXTENSION tsm_system_row;

ALTER TABLE game DROP CONSTRAINT check_attempt_maxAttempts;

DROP INDEX playerName_idx;