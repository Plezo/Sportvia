CREATE TABLE IF NOT EXISTS game (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    userID          UUID        NOT NULL,
    playerName      TEXT        NOT NULL,
    age             SMALLINT    NOT NULL,
    height          SMALLINT    NOT NULL,
    team            TEXT        NOT NULL,
    conference      TEXT        NOT NULL,
    division        TEXT        NOT NULL,
    position        TEXT        NOT NULL,
    playerNumber    SMALLINT    NOT NULL,
    playerImage     TEXT        NOT NULL,
    attempt         SMALLINT    NOT NULL,
    maxAttempts     SMALLINT    NOT NULL,
    win             BOOLEAN     NOT NULL DEFAULT FALSE,
    createdAt       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS player (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    playerName      TEXT        NOT NULL UNIQUE,
    age             SMALLINT    NOT NULL,
    height          SMALLINT    NOT NULL,
    team            TEXT        NOT NULL,
    conference      TEXT        NOT NULL,
    division        TEXT        NOT NULL,
    position        TEXT        NOT NULL,
    playerNumber    SMALLINT    NOT NULL,
    playerImage     TEXT        NOT NULL,
    updatedAt       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);