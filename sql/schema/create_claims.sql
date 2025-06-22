CREATE TABLE claims ( 
    claim_id SERIAL PRIMARY KEY, 
    document_id INTEGER REFERENCES documents(document_id), -- foreign key to documents 
    claim_text TEXT, 
    start_char INTEGER,
    end_char INTEGER, 
    confidence_score REAL,
    extracted_at TIMESTAMP,

    -- Add eventually: 
    -- 	•	NOT NULL on essentials like claim_text and document_id
	--  •	A CHECK on confidence_score for sanity
	--  •	ON DELETE CASCADE if your app will regularly delete documents
)