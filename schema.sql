
CREATE TABLE apikey
(
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "app_version" character varying(255) NOT NULL,
    "revoked" boolean NOT NULL,
    "expiration_time" timestamp without time zone NOT NULL,
    "create_time" timestamp without time zone NOT NULL,
    CONSTRAINT apikey_pkey PRIMARY KEY (id)
);

CREATE TABLE "user"
(
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "email" character varying(255) NOT NULL,
    "password" character varying(255) NOT NULL,
    "fullname" character varying(255) NOT NULL,
    "username" character varying(255) NOT NULL,
    "create_time" timestamp without time zone NOT NULL,
    CONSTRAINT user_pkey PRIMARY KEY (id),
    CONSTRAINT user_email_key UNIQUE (email)
);

CREATE TABLE refresh_token
(
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL,
    "expiration_time" timestamp without time zone NOT NULL,
    "create_time" timestamp without time zone NOT NULL,
    CONSTRAINT refreshtoken_pkey PRIMARY KEY (id),
    CONSTRAINT refreshtoken_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES "user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE access_token
(
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL,
    "refresh_token_id" uuid NOT NULL,
    "expiration_time" character varying(255) NOT NULL,
    "create_time" timestamp without time zone NOT NULL,
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

CREATE TABLE chat
(
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" character varying(255),
    "create_time" timestamp without time zone NOT NULL,
    CONSTRAINT chat_pkey PRIMARY KEY (id)
);

CREATE TABLE union_user_and_chat
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    chat_id uuid NOT NULL,
    create_time timestamp with time zone NOT NULL,
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

CREATE TABLE message
(
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "chat_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "content" character varying(255) NOT NULL,
    "create_time" timestamp without time zone NOT NULL,
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