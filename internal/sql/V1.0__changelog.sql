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
CREATE TABLE "user"
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
    CONSTRAINT accesstoken_pkey PRIMARY KEY (id),
    CONSTRAINT accesstoken_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES "user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT accesstoken_refreshtoken_id_fkey FOREIGN KEY (refresh_token_id)
        REFERENCES "refresh_token" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);
--rollback DROP TABLE access_token;

--changeset daniarmas:7 labels:create-chat-table context:example-context
--comment: creating the chat table
CREATE TABLE chat
(
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" character varying(255),
    "create_time" timestamp without time zone NOT NULL,
    "update_time" timestamp without time zone NOT NULL,
    "delete_time" timestamp without time zone,
    CONSTRAINT chat_pkey PRIMARY KEY (id)
);
--rollback DROP TABLE chat;

--changeset daniarmas:8 labels:create-table-union-user-chat context:example-context
--comment: creating the union_user_and_chat table
CREATE TABLE union_user_and_chat
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    chat_id uuid NOT NULL,
    create_time timestamp with time zone NOT NULL,
    update_time timestamp with time zone NOT NULL,
    delete_time timestamp with time zone,
    CONSTRAINT union_user_and_chat_pkey PRIMARY KEY (id),
    CONSTRAINT union_user_and_chat_fkey_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES "user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT union_user_and_chat_fkey_chat_id_fkey FOREIGN KEY (chat_id)
        REFERENCES chat (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);
--rollback DROP TABLE "union_user_and_chat";

--changeset daniarmas:9 labels:create-message-table context:example-context
--comment: creating the message table
CREATE TABLE message
(
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "chat_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "content" character varying(255) NOT NULL,
    "create_time" timestamp without time zone NOT NULL,
    "update_time" timestamp without time zone NOT NULL,
    "delete_time" timestamp without time zone,
    CONSTRAINT message_pkey PRIMARY KEY (id),
    CONSTRAINT message_fkey FOREIGN KEY (user_id)
        REFERENCES "user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT message_chat_id_fkey FOREIGN KEY (chat_id)
        REFERENCES "chat" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);
--rollback DROP TABLE message;