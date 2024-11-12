DROP INDEX IF EXISTS idx_contents_search_vector;
ALTER TABLE contents DROP COLUMN IF EXISTS search_vector;
