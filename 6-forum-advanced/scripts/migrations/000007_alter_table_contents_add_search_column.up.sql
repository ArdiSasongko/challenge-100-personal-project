ALTER TABLE contents ADD COLUMN search_vector tsvector;
UPDATE contents SET search_vector = to_tsvector('english', content_title || ' ' || content_body);
CREATE INDEX idx_contents_search_vector ON contents USING GIN (search_vector);