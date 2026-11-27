CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    telegram_user_id BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE messages
    ADD COLUMN user_id BIGINT REFERENCES users (id),
    ADD CONSTRAINT messages_user_id_id_key UNIQUE (user_id, id);

CREATE TABLE profile_sections (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    section_key TEXT NOT NULL CHECK (section_key ~ '^[a-z][a-z0-9_]*$'),
    title TEXT NOT NULL CHECK (btrim(title) <> ''),
    description TEXT NOT NULL CHECK (btrim(description) <> ''),

    UNIQUE (user_id, section_key),
    UNIQUE (user_id, id)
);

CREATE TABLE profile_facts (
    user_id BIGINT NOT NULL,
    section_id BIGINT NOT NULL,
    fact_key TEXT NOT NULL CHECK (btrim(fact_key) <> ''),
    value JSONB NOT NULL,
    source_message_id INTEGER NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, section_id, fact_key),

    FOREIGN KEY (user_id, section_id)
        REFERENCES profile_sections (user_id, id) ON DELETE CASCADE,

    FOREIGN KEY (user_id, source_message_id)
        REFERENCES messages (user_id, id) ON DELETE CASCADE
);

CREATE INDEX profile_facts_source_idx
    ON profile_facts (user_id, source_message_id);