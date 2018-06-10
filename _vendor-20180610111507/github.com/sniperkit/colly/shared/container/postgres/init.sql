DROP TABLE IF EXISTS colly_cache_response;

CREATE TABLE colly_cache_response (
  id VARCHAR(32) PRIMARY KEY,
  body TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL
);