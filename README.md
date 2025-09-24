# claim-extraction
# AI Ethics Claim Extraction Pipeline

This project is a **multi-stage NLP pipeline** designed to extract claims from documents in the domain of AI ethics and safety.  
It combines **Go**, **AWS Lambda**, **S3**, **Postgres**, and **spaCy** for a fully automated document ingestion and processing workflow.

---

## Table of Contents

- [Overview](#overview)  
- [Architecture](#architecture)  
- [Setup](#setup)  
- [Usage](#usage)  
- [Pipeline Flow](#pipeline-flow)  
- [Future Work](#future-work)  

---

## Overview

The goal of this project is to:

- Ingest documents (PDFs) into an AWS-backed pipeline  
- Extract textual content  
- Detect and label claims related to AI ethics and safety  
- Store and index metadata for search and further analysis  

This pipeline is designed to be modular so that it can be expanded for other domains in the future.

---

## Architecture

The pipeline consists of three main components:

1. **Go CLI Tool**  
   Uploads documents to S3 and logs metadata into a Postgres database.  
   Waits for processing completion notifications via AWS SNS/SQS.

2. **AWS Lambda (Python)**  
   Triggered on file upload.  
   Extracts text from documents and stores processed text back into S3.  
   Publishes completion notifications to SNS.

3. **Python NLP Processor**  
   Downloads processed text from S3.  
   Runs claim extraction using **spaCy** with custom rules and named entity recognition.  

---

## Setup

### Requirements

- Go ≥ 1.20  
- Python ≥ 3.10  
- Docker  
- AWS CLI configured  
- AWS account with IAM permissions for S3, Lambda, SNS, SQS  

### Steps

1. Clone this repository  
2. Set up Postgres locally (via Docker or otherwise)  
3. Create an S3 bucket and SNS topic in AWS  
4. Deploy Lambda function for text extraction  
5. Configure IAM roles and permissions  
6. Install Python dependencies:
   ```bash
   pip install -r requirements.txt
    ```
7. Install spaCy language models: 
```bash
    pip -m spacy download en_core_web_trf
```
---
# IN PROGRESS
