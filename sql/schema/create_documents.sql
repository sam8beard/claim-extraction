CREATE TABLE documents ( 
    document_id SERIAL PRIMARY KEY, 
    uploaded_at TIMESTAMP, 
    file_name TEXT, 
    source TEXT, 
    text_extracted TIMESTAMP, 
    content_hash TEXT
    s3_key TEXT 
)