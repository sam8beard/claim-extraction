# claimEx
A command line tool for extracting structured information from argument-rich/claim-rich PDF documents

## How It Works
*	Search for any topic via [SearXNG][searxng] 
*	Specify a number of files to aggregate
*	Files are processed using a [spaCy][spacy] Span Categorizer [(SpanCat)][spancat] model trained on ~1500 silver labels to detect and extract claim spans
*	View analysis for each document returned in JSON format 
--- 
## Pipeline Results
> ### **Spans**
* **Sources** 
    * *Who* made the claim
* **Claim Verbs**
    * The *verb* used to make the claim
* **Claim Modifiers** 
    * *Modifier(s)* that indicate the strength/degree with which the claim was made
* **Claim Contents** 
    * The *claim* being made

> ### **Other Info**
* **Origin Document** 
    * The **document** spans were extracted from
* **Origin Sentence**
    * The *sentence* that contains a given span
* **Claim Density Score**
    * A value representing how *claim-heavy* a given document is
* **Confidence Score**
    * How *confident* the model was at predicting a given span

> #### NOTE: File sizes of up to ~200MB are recommended for optimal performance
---
## Powered By
 #### [**spaCy**][spacy] for NLP 
 #### [**SearXNG**][searxng] for search
 #### [**MinIO**][minio] for object storage
----
## Getting Started
> Project in progress 

[spacy]:https://github.com/explosion/spaCy
[searxng]: https://github.com/searxng/searxng
[minio]: https://github.com/minio/minio
[spancat]:https://github.com/explosion/spaCy/blob/master/spacy/pipeline/spancat.py
