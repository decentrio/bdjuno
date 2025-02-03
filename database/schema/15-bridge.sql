CREATE TABLE bridge_in
(
    hash     TEXT   NOT NULL REFERENCES transaction (hash) PRIMARY KEY,
    amount   TEXT,
    denom    TEXT,
    receiver TEXT
);

CREATE TABLE bridge_out
(
    hash     TEXT   NOT NULL REFERENCES transaction (hash) PRIMARY KEY,
    amount   TEXT,
    denom    TEXT,
    sender   TEXT
);


CREATE TABLE rate_limit
(
    denom             TEXT   NOT NULL,
    rate_limit        TEXT,
    inflow            TEXT,
    height            BIGINT NOT NULL,
    PRIMARY KEY (denom, height)
);