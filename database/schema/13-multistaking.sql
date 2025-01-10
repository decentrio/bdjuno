CREATE TYPE MS_COIN AS
(
    denom  TEXT,
    amount TEXT,
    bond_weight NUMERIC
);

CREATE TYPE MS_UNLOCK_ENTRY AS
(
    creation_height  BIGINT,
    denom  TEXT,
    amount TEXT,
    bond_weight NUMERIC
);

/* ---- PARAMS ---- */

CREATE TABLE ms_locks
(
    staker_addr TEXT    NOT NULL,
    val_addr    TEXT    NOT NULL,
    ms_lock MS_COIN[] NOT NULL,
    height     BIGINT  NOT NULL,
    PRIMARY KEY (staker_addr, val_addr)
);

CREATE TABLE ms_unlocks
(
    staker_addr TEXT    NOT NULL,
    val_addr    TEXT    NOT NULL,
    unlock_entry MS_UNLOCK_ENTRY[] NOT NULL,
    height     BIGINT  NOT NULL,
    PRIMARY KEY (staker_addr, val_addr)
);
CREATE INDEX ms_locks_height_index ON ms_locks (height);
CREATE INDEX ms_unlocks_height_index ON ms_unlocks (height);

CREATE TABLE validator_denom
(
    val_addr   TEXT NOT NULL REFERENCES validator (consensus_address) PRIMARY KEY,
    denom      TEXT NOT NULL,
    height     BIGINT     NOT NULL
);
CREATE INDEX validator_denom_height_index ON validator_denom (height);

CREATE TABLE ms_event
(
    height       BIGINT  NOT NULL REFERENCES block (height),
    name         TEXT    NOT NULL,
    val_addr         TEXT    NOT NULL,
    del_addr         TEXT    NOT NULL,
    amount         TEXT    NOT NULL
);
CREATE INDEX ms_event_height_index ON ms_event (height);

CREATE TABLE token_unbonding
(
    denom TEXT    NOT NULL PRIMARY KEY,
    amount    TEXT    NOT NULL,
    height     BIGINT  NOT NULL
);
CREATE INDEX token_unbonding_height_index ON token_unbonding (height);

CREATE TABLE token_bonded
(
    denom TEXT    NOT NULL PRIMARY KEY,
    amount    TEXT    NOT NULL,
    height     BIGINT  NOT NULL
);
CREATE INDEX token_bonded_height_index ON token_bonded (height);