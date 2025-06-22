CREATE TABLE entities ( 
    entity_id SERIAL PRIMARY KEY, 
    claim_id INTEGER REFERENCES claims(claim_id), 
    entity_label TEXT, 
    entity_text TEXT
)