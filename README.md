# ClaimEX
A CLI data aggregation tool for extracting structured information from argument-rich/claim-rich PDF documents

Features
---
*	Search for a topic via [SearXNG][searxng] 
*	Specify a number of files to aggregate
*	Process files using a pre-trained Span Categorizer [(SpanCat)][spancat] model to detect and extract claim spans with the click of a button
*	View analysis for each document, including:
      * Sources
      * Extracted spans
      * Claim density score (indicating how claim-heavy the document is)

## Powered By
 #### [**spaCy**][spacy] for NLP 
 #### [**SearXNG**][searxng] for search
 #### [**MinIO**][minio] for object storage

----
# IN PROGRESS...


[spacy]:https://github.com/explosion/spaCy
[searxng]: https://github.com/searxng/searxng
[minio]: https://github.com/minio/minio
[spancat]:https://github.com/explosion/spaCy/blob/master/spacy/pipeline/spancat.py
