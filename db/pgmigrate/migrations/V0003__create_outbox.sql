CREATE TABLE
    IF NOT EXISTS outbox (
        id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        event_type text NOT NULL,
        event_version integer NOT NULL DEFAULT 1,
        aggregate_type text NOT NULL,
        aggregate_id bigint NOT NULL,
        payload jsonb NOT NULL,
        created_at timestamptz NOT NULL DEFAULT now (),
        published_at timestamptz,
        attempts integer NOT NULL DEFAULT 0,
        next_attempt_at timestamptz NOT NULL DEFAULT now (),
        last_error text,
        locked_until timestamptz,
        claim_token text,
        CHECK (event_version > 0),
        CHECK (attempts >= 0),
        CHECK ((locked_until IS NULL) = (claim_token IS NULL))
    );

CREATE INDEX outbox_unpublished_idx ON outbox (next_attempt_at, id)
WHERE
    published_at IS NULL;