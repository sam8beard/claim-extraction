CREATE TABLE documents ( 
    id SERIAL PRIMARY KEY, 
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    file_name TEXT NOT NULL, 
    source TEXT NOT NULL, 
    text_extracted BOOLEAN NOT NULL DEFAULT FALSE, 
    content_hash TEXT NOT NULL
    s3_key TEXT NOT NULL
)