CREATE TABLE IF NOT EXISTS auth.users3 (
    instance_id STRING NOT NULL,
    id STRING NOT NULL,
    password_set BOOL NULL,
    password_change TIMESTAMPTZ NULL,
    last_login TIMESTAMPTZ NULL,
    init_required BOOL NULL,
    mfa_init_skipped TIMESTAMPTZ NULL,
    username_change_required BOOL NULL,
    passwordless_init_required BOOL NULL,
    password_init_required BOOL NULL,

    PRIMARY KEY (instance_id ASC, id ASC)
)
