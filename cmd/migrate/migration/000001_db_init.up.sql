CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE video (
    id BIGINT PRIMARY KEY,
    category CITEXT,
    name TEXT,
    slug TEXT UNIQUE,
    origin_name TEXT,
    poster_url TEXT,
    thumb_url TEXT,
    description TEXT,
    link_embed TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE actor (
    id SERIAL PRIMARY KEY,
    name CITEXT UNIQUE
);

CREATE TABLE tag (
    id SERIAL PRIMARY KEY,
    name CITEXT UNIQUE
);

CREATE TABLE studio (
    id SERIAL PRIMARY KEY,
    name CITEXT UNIQUE
);

CREATE TABLE director (
    id SERIAL PRIMARY KEY,
    name CITEXT UNIQUE
);

CREATE TABLE video_actor (
    video_id BIGINT REFERENCES video(id) ON DELETE CASCADE,
    actor_id INT REFERENCES actor(id) ON DELETE CASCADE,
    PRIMARY KEY (video_id, actor_id)
);

CREATE TABLE video_tag (
    video_id BIGINT REFERENCES video(id) ON DELETE CASCADE,
    tag_id INT REFERENCES tag(id) ON DELETE CASCADE,
    PRIMARY KEY (video_id, tag_id)
);

CREATE TABLE video_studio (
    video_id BIGINT REFERENCES video(id) ON DELETE CASCADE,
    studio_id INT REFERENCES studio(id) ON DELETE CASCADE,
    PRIMARY KEY (video_id, studio_id)
);

CREATE TABLE video_director (
    video_id BIGINT REFERENCES video(id) ON DELETE CASCADE,
    director_id INT REFERENCES director(id) ON DELETE CASCADE,
    PRIMARY KEY (video_id, director_id)
);

-- Indexes
CREATE INDEX idx_video_name_trgm ON video USING gin (name gin_trgm_ops);
CREATE INDEX idx_video_category ON video (category);
CREATE INDEX idx_video_created_at ON video (created_at DESC);
CREATE INDEX idx_video_updated_at ON video (updated_at DESC);

CREATE INDEX idx_video_actor ON video_actor (actor_id);
CREATE INDEX idx_video_tag ON video_tag (tag_id);
CREATE INDEX idx_video_studio ON video_studio (studio_id);
CREATE INDEX idx_video_director ON video_director (director_id);