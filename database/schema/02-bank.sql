/* ---- SUPPLY ---- */

CREATE TABLE supply
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    coins      COIN[]  NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX supply_height_index ON supply (height);

CREATE TABLE token_holder
(
    denom      TEXT    NOT NULL PRIMARY KEY,
    num_holder BIGINT  NOT NULL,
    height     BIGINT  NOT NULL
);

CREATE INDEX token_holder_height_index ON token_holder (height);