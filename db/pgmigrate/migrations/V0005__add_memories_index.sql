CREATE INDEX memories_content_fts_idx
    ON memories USING gin (to_tsvector('russian', content));