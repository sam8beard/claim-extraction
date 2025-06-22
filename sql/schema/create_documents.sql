CREATE TABLE documents ( 
    document_id serial PRIMARY KEY, 
    uploaded_at timestamp, 
    file_name text, 
    source text, 
    text_extracted timestamp, 
    content_hash text
)