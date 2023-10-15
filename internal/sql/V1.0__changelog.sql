--liquibase formatted sql

--changeset daniarmas:1 labels:create-uuid-extension context:example-context
--comment: creating the uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
--rollback DROP EXTENSION "uuid-ossp";


--changeset daniarmas:2 labels:create-postgis-extension context:example-context
--comment: creating the postgis extension
CREATE EXTENSION IF NOT EXISTS "postgis";
--rollback DROP EXTENSION "postgis";

--changeset daniarmas:3 labels:create-apikey-table    context:example-context
--comment: creating the province table
CREATE TABLE apikey
(
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "app_version" character varying(255) NOT NULL,
    "revoked" boolean NOT NULL,
    "expiration_time" timestamp without time zone NOT NULL,
    "create_time" timestamp without time zone NOT NULL,
    "update_time" timestamp without time zone NOT NULL,
    "delete_time" timestamp without time zone,
    CONSTRAINT apikey_pkey PRIMARY KEY (id)
);
--rollback DROP TABLE apikey;