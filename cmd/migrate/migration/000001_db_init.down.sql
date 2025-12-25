-- Drop pivot tables first
DROP TABLE IF EXISTS video_director;
DROP TABLE IF EXISTS video_studio;
DROP TABLE IF EXISTS video_tag;
DROP TABLE IF EXISTS video_actor;

-- Drop master tables
DROP TABLE IF EXISTS director;
DROP TABLE IF EXISTS studio;
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS actor;

-- Drop main table
DROP TABLE IF EXISTS video;

-- Drop extensions
DROP EXTENSION IF EXISTS citext;
DROP EXTENSION IF EXISTS pg_trgm;