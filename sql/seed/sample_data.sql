-- sample documents
INSERT INTO documents (uploaded_at, file_name, source, text_extracted, content_hash, s3_key) VALUES
('2025-06-20 10:00:00', 'climate_report.pdf', 'UN Climate Agency', '2025-06-20 10:05:00', 'abc123def', 'documents/climate_report.pdf'),
('2025-06-21 14:30:00', 'tech_news.txt', 'TechCrunch', '2025-06-21 14:32:00', 'xyz789ghi', 'documents/tech_news.txt');

-- sample claims
INSERT INTO claims (document_id, claim_text, start_char, end_char, confidence_score, extracted_at) VALUES
(1, 'Global temperatures have risen by 1.2°C since 1900.', 0, 49, 0.95, '2025-06-20 10:06:00'),
(1, 'Sea levels are expected to rise by 30cm by 2050.', 51, 98, 0.88, '2025-06-20 10:06:05'),
(2, 'AI startups raised $5B in Q2 2025.', 0, 33, 0.92, '2025-06-21 14:33:00');

-- sample entities
INSERT INTO entities (claim_id, entity_label, entity_text) VALUES
(1, 'DATE', 'since 1900'),
(1, 'MEASURE', '1.2°C'),
(2, 'DATE', 'by 2050'),
(2, 'MEASURE', '30cm'),
(3, 'ORG', 'AI startups'),
(3, 'MONEY', '$5B'),
(3, 'DATE','Q2 2025');