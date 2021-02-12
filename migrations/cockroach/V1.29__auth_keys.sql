CREATE TABLE auth.authn_keys
(
    key_id          TEXT,
    object_id       TEXT,
    object_type     SMALLINT,
    auth_identifier TEXT,

    key_type        SMALLINT,
    sequence        BIGINT,
    expiration_date TIMESTAMPTZ,
    creation_date   TIMESTAMPTZ,
    public_key      BYTES,

    PRIMARY KEY (key_id, object_id, object_type, auth_identifier)
);

INSERT INTO auth.authn_keys (
    key_id,
    object_id,
    object_type,
    auth_identifier,
    key_type,
    sequence,
    expiration_date,
    creation_date,
    public_key
    )
    SELECT
        id,
        user_id,
        0,
        user_id,
        machine_type,
        sequence,
        expiration_date,
        creation_date,
        public_key
    FROM auth.machine_keys;

CREATE TABLE management.authn_keys
(
    key_id          TEXT,
    object_id       TEXT,
    object_type     SMALLINT,
    auth_identifier TEXT,

    key_type        SMALLINT,
    sequence        BIGINT,
    expiration_date TIMESTAMPTZ,
    creation_date   TIMESTAMPTZ,
    public_key      BYTES,

    PRIMARY KEY (key_id, object_id, object_type, auth_identifier)
);

INSERT INTO management.authn_keys (
    key_id,
    object_id,
    object_type,
    auth_identifier,
    key_type,
    sequence,
    expiration_date,
    creation_date,
    public_key
)
SELECT
    id,
    user_id,
    0,
    user_id,
    machine_type,
    sequence,
    expiration_date,
    creation_date,
    public_key
FROM management.machine_keys;

INSERT INTO auth.current_sequences (view_name, event_timestamp, current_sequence, last_successful_spooler_run)
    SELECT view_name, event_timestamp, current_sequence, last_successful_spooler_run FROM auth.current_sequences WHERE view_name = 'auth.machine_keys';

INSERT INTO management.current_sequences (view_name, event_timestamp, current_sequence, last_successful_spooler_run)
    SELECT view_name, event_timestamp, current_sequence, last_successful_spooler_run FROM management.current_sequences WHERE view_name = 'management.machine_keys';

ALTER TABLE auth.authn_keys OWNER TO admin;
ALTER TABLE management.authn_keys OWNER TO admin;
