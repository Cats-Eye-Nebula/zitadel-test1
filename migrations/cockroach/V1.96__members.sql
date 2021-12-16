CREATE TABLE zitadel.projections.org_members (
    org_id STRING NOT NULL
    , user_id STRING NOT NULL
    , roles STRING[]

    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , sequence INT8 NOT NULL
    , resource_owner STRING NOT NULL

    , PRIMARY KEY (org_id, user_id)
    , INDEX idx_user (user_id)
);

CREATE TABLE zitadel.projections.iam_members (
    iam_id STRING NOT NULL
    , user_id STRING NOT NULL
    , roles STRING[]

    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , sequence INT8 NOT NULL
    , resource_owner STRING NOT NULL

    , PRIMARY KEY (iam_id, user_id)
    , INDEX idx_user (user_id)
);

CREATE TABLE zitadel.projections.project_members (
    project_id STRING NOT NULL
    , user_id STRING NOT NULL
    , roles STRING[]

    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , sequence INT8 NOT NULL
    , resource_owner STRING NOT NULL

    , PRIMARY KEY (project_id, user_id)
    , INDEX idx_user (user_id)
);

CREATE TABLE zitadel.projections.project_grant_members (
    project_id STRING NOT NULL
    , user_id STRING NOT NULL
    , grant_id STRING
    , roles STRING[]

    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , sequence INT8 NOT NULL
    , resource_owner STRING NOT NULL

    , PRIMARY KEY (project_id, grant_id, user_id)
    , INDEX idx_user (user_id)
);
