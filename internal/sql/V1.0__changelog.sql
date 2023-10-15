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

--changeset daniarmas:4 labels:create-user-table context:example-context
--comment: creating the user table
CREATE TABLE user
(
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "email" character varying(255) NOT NULL,
    "password" character varying(255) NOT NULL,
    "fullname" character varying(255) NOT NULL,
    "username" character varying(255) NOT NULL,
    "create_time" timestamp without time zone NOT NULL,
    "update_time" timestamp without time zone NOT NULL,
    "delete_time" timestamp without time zone,
    CONSTRAINT user_pkey PRIMARY KEY (id),
    CONSTRAINT user_email_key UNIQUE (email)
);
--rollback DROP TABLE user;

--changeset daniarmas:5 labels:create-refresh-token-table context:example-context
--comment: creating the refresh-token table
CREATE TABLE refresh_token
(
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL,
    "expiration_time" character varying(255) NOT NULL,
    "create_time" timestamp without time zone NOT NULL,
    "update_time" timestamp without time zone NOT NULL,
    "delete_time" timestamp without time zone,
    CONSTRAINT refreshtoken_pkey PRIMARY KEY (id),
    CONSTRAINT refreshtoken_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES "user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);
--rollback DROP TABLE refresh_token;

--changeset daniarmas:6 labels:create-access_token-table context:example-context
--comment: creating the access_token table
CREATE TABLE access_token
(
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL,
    "refresh_token_id" uuid NOT NULL,
    "expiration_time" character varying(255) NOT NULL,
    "create_time" timestamp without time zone NOT NULL,
    "update_time" timestamp without time zone NOT NULL,
    "delete_time" timestamp without time zone,
    CONSTRAINT refreshtoken_pkey PRIMARY KEY (id),
    CONSTRAINT refreshtoken_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES "user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT refreshtoken_refreshtoken_id_fkey FOREIGN KEY (refresh_token_id)
        REFERENCES "refresh_token" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);
--rollback DROP TABLE access_token;