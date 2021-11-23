CREATE TABLE zitadel.projections.label_policies (
    id STRING NOT NULL,
    creation_date TIMESTAMPTZ NOT NULL,
    change_date TIMESTAMPTZ NOT NULL,
    sequence INT8 NOT NULL,
    state INT2 NOT NULL,
    resource_owner TEXT NOT NULL,
    
    is_default BOOLEAN NOT NULL DEFAULT false,
    hide_login_name_suffix BOOLEAN NOT NULL DEFAULT false,
	font_url STRING,
	watermark_disabled BOOLEAN NOT NULL DEFAULT false,
	should_error_popup BOOLEAN NOT NULL DEFAULT false,
	light_primary_color STRING,
	light_warn_color STRING,
	light_background_color STRING,
	light_font_color STRING,
	light_logo_url STRING,
	light_icon_url STRING,
	dark_primary_color STRING,
	dark_warn_color STRING,
	dark_background_color STRING,
	dark_font_color STRING,
	dark_logo_url STRING,
	dark_icon_url STRING,

	PRIMARY KEY (id)
);