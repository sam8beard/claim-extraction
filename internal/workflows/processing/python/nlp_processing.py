import sys
import base64
import json
import re
import unicodedata
import spacy
from collections import Counter
import warnings

warnings.filterwarnings("ignore")

BODY_DELIMETER = b"--END-BODY--\n"
BUF_SIZE = 4096


def chunk_text(text, chunk_size=3):
    raw_sents = re.split(r'[.!?]\s+(?=[A-Z0-9])', text)
    raw_sents = [s.strip() for s in raw_sents if len(s.strip()) > 20]

    for i in range(0, len(raw_sents), chunk_size):
        chunk_text = " ".join(raw_sents[i:i+chunk_size])
        yield chunk_text


def preprocess_text(text):
    text = unicodedata.normalize("NFC", text)
    text = "".join(c for c in text if c.isprintable())
    text = re.sub(r"\s+", " ", text)
    text = text.strip()
    return text


def compute_claim_score(spans):
    types = ["SOURCE", "CLAIM_VERB", "CLAIM_CONTENTS"]

    # Count spans per type
    counts = Counter(span["type"] for span in spans if span["type"] in types)
    total = sum(counts.values())

    if total == 0:
        return 0.0  # No spans, no score

    # Get proportions for each type
    proportions = [counts[t] / total for t in types]

    # The more even proportion, the lower variance, the higher score
    mean = sum(proportions) / len(types)
    variance = sum((p - mean) ** 2 for p in proportions) / len(types)

    score = 1 - variance
    return round(score, 3)


def run_spancat(object_key, raw_text, file_name, nlp_model):
    text = preprocess_text(raw_text)
    file_data = {
        "objectKey": object_key,
        "fileName": file_name,
        "error": "",
        "claimScore": 0.0,
        "claimSpans": []
    }

    for doc in nlp_model.pipe(chunk_text(text), batch_size=1):
        spans = doc.spans['sc']
        for span, confidence in zip(spans, spans.attrs["scores"]):
            file_data["claimSpans"].append({
                "text": span.text,
                "type": span.label_,
                "sent": span.sent.text,
                "confidence": float(confidence)
            })

    # Compute claim score
    file_data["claimScore"] = compute_claim_score(file_data["claimSpans"])

    return file_data


def main():
    # Load model once at startup
    nlp_model = spacy.load("spancat_v1.0")
    nlp_model.add_pipe("sentencizer", before="spancat")

    # Process files one at a time in streaming fashion
    for meta_line in sys.stdin:
        meta_line = meta_line.rstrip("\n")
        if not meta_line:
            continue

        # Parse metadata
        try:
            meta = json.loads(meta_line)
        except Exception as e:
            err_obj = {
                "error": f"invalid json metadata: {str(e)}",
                "objectKey": ""
            }
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue

        # Read body until sentinel
        try:
            leftover = b''
            chunks = []
            while True:
                chunk = sys.stdin.buffer.read(BUF_SIZE)
                if not chunk:
                    break
                combined = leftover + chunk

                idx = combined.find(BODY_DELIMETER)
                if idx >= 0:
                    chunks.append(combined[:idx])
                    leftover = combined[idx+len(BODY_DELIMETER):]
                    break
                else:
                    keep = combined[-(len(BODY_DELIMETER)-1):]
                    chunks.append(combined[:-len(keep)])
                    leftover = keep

            body_b64 = b''.join(chunks)
        except Exception as e:
            err_obj = {
                "error": f"error reading body: {str(e)}",
                "objectKey": meta.get("objectKey", "")
            }
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue

        # Process immediately
        try:
            # Decode base64
            missing_padding = len(body_b64) % 4
            if missing_padding:
                bytes_body = base64.b64decode(
                    body_b64 + b'=' * (4 - missing_padding))
            else:
                bytes_body = base64.b64decode(body_b64)

            # Decode text (try UTF-8, fallback to Latin-1)
            try:
                text_body = bytes_body.decode('utf-8')
            except UnicodeDecodeError:
                try:
                    text_body = bytes_body.decode('latin-1')
                except Exception as e:
                    err_obj = {
                        "error": f"text decode error: {str(e)}",
                        "objectKey": meta.get("objectKey", "")
                    }
                    print(json.dumps(err_obj), file=sys.stderr, flush=True)
                    continue

            # Run NLP processing
            object_key = meta.get("objectKey", "")
            file_name = meta.get("fileName", "")
            result = run_spancat(object_key, text_body, file_name, nlp_model)

            # Write result immediately to stdout
            print(json.dumps(result), file=sys.stdout, flush=True)

        except Exception as e:
            err_obj = {
                "error": f"processing error: {str(e)}",
                "objectKey": meta.get("objectKey", "")
            }
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue


if __name__ == "__main__":
    main()
