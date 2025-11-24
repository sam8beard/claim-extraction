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
BUF_SIZE = 50000


# def test():
#    test_files = [
#        {
#            "objectKey": "file1.txt",
#            "fileName": "file1",
#            "content": (
#                "Dr. Smith claims that the vaccine is highly effective. "
#                "However, recent studies suggest otherwise. "
#                "The WHO confirms that further testing is required. "
#                "Experts debate the methodology used in these studies. "
#                "Ultimately, the evidence remains inconclusive."
#            )
#        },
#        {
#            "objectKey": "file2.txt",
#            "content": (
#                "Alice asserts that the new policy will reduce emissions. "
#                "Bob counters that the economic impact will be severe. "
#                "The government releases official figures supporting Alice's claim. "
#                "Environmental groups respond positively, emphasizing long-term benefits. "
#                "Analysts remain skeptical about the short-term effects."
#            )
#        },
#        {
#            "objectKey": "file3.txt",
#            "content": (
#                "The article reports that inflation has risen by 3.2% over the past quarter. "
#                "Economists warn that this trend may continue if interest rates are not adjusted. "
#                "Consumer groups note the rising cost of living as evidence. "
#                "Some banks argue that this increase is within expected limits."
#            )
#        },
#        {
#            "objectKey": "file4.txt",
#            "content": (
#                "NASA confirms that the satellite has successfully entered orbit. "
#                "Scientists observe minor deviations in trajectory, which are under investigation. "
#                "Independent analysts suggest that the mission's success will advance space research. "
#                "The agency releases images and data for public review."
#            )
#        },
#        {
#            "objectKey": "file5.txt",
#            "content": (
#                "According to several reports, the technology startup achieved record revenue last year. "
#                "Investors praise the management team for strategic decisions. "
#                "Competitors question the sustainability of such growth. "
#                "Financial analysts provide a cautious outlook for the upcoming fiscal year."
#            )
#        }
#    ]
#
#    for f in test_files:
#        json_line = json.dumps(f)
#        result = test_main(json_line)
#
#        if isinstance(result, Exception):
#            print("Error processing")
#        else:
#            print(json.dumps(result, indent=2))
#
#
# def test_main(json_line):
#    nlp_model = spacy.load("spancat_v1.0")
#    nlp_model.add_pipe("sentencizer", before="spancat")
#    json_line = json_line.strip()
#
#    try:
#        obj = json.loads(json_line)
#        object_key = obj.get("objectKey", "")
#        content = obj.get("content", "")
#        file_name = obj.get("fileName", "")
#        result = run_spancat(object_key, content, file_name, nlp_model)
#        return result
#
#     # handle stderr here
#    except Exception as e:
#        return e


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

    # count spans per type
    counts = Counter(span["type"] for span in spans if span["type"] in types)
    total = sum(counts.values())

    if total == 0:
        return 0.0  # no spans, no score

    # get proportions for each type
    proportions = [counts[t] / total for t in types]

    # the more even proportion, the lower variance, the higher score
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

        # compute claim score
        file_data["claimScore"] = compute_claim_score(file_data["claimSpans"])

    return file_data


def process_bodies(bodies, nlp_model):
    results = []
    for k, file in bodies.items():
        meta = file.get('meta')
        body = file.get('body_b64')
        try:
            missing_padding = len(body) % 4
            if missing_padding:
                bytes_body = base64.b64decode(body + b'==')
            else:
                bytes_body = base64.b64decode(body)
        except Exception as e:
            meta_dict = dict(meta)
            err_obj = {
                "error": "testing",
                # "error": f"b64 decode error on file {str(meta_dict.get('objectKey', ''))}: {str(e)}"
            }
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue
        try:
            temp = ''
            # try utf-8
            temp = bytes_body.decode('utf-8')
        except Exception:
            try:
                # if that fails, try latin-1
                temp = bytes_body.decode('latin-1')
            except Exception as e:
                meta_dict = dict(meta)
                err_obj = {
                    "error": "testing",
                    # "error": f"utf-8/latin-1 decode error on file {str(meta_dict.get('objectKey', ''))}: {str(e)}"
                }
                print(json.dumps(err_obj), file=sys.stderr, flush=True)
                continue
        text_body = temp
        try:
            out_meta = dict(meta)
            out_meta['error'] = ""
            object_key = out_meta.get("objectKey", "")
            file_name = out_meta.get("fileName", "")
            result = run_spancat(object_key, text_body, file_name, nlp_model)
            results.append(result)

        except Exception as e:
            err_obj = {"error": f"nlp processing error: {str(e)}"}
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue
    return results


def main():
    bodies = dict()
    nlp_model = spacy.load("spancat_v1.0")
    nlp_model.add_pipe("sentencizer", before="spancat")

    for meta_line in sys.stdin:
        meta_line = meta_line.rstrip("\n")
        if not meta_line:
            continue
        try:
            meta = json.loads(meta_line)
        except Exception as e:

            msg = str(e)
            err_obj = {
                "error": f"invalid json metadata: {msg}",
                "objectKey": ""
            }
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue
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

        # save body for processing
        body_b64 = b''.join(chunks)
        key = meta.get('objectKey')
        bodies[key] = {"meta": meta, "body_b64": body_b64}

    # process all bodies
    results = process_bodies(bodies, nlp_model)

    # write all results
    for result in results:
        try:
            print(json.dumps(result), file=sys.stdout, flush=True)
        except Exception as e:
            err_obj = {"error": str(e)}
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
#
#    # decode/process bodies, and write back
#    for meta, body in bodies:
#        try:
#            missing_padding = len(body) % 4
#            if missing_padding:
#                bytes_body = base64.b64decode(body + b'==')
#            else:
#                bytes_body = base64.b64decode(body)
#        except Exception as e:
#            meta_dict = dict(meta)
#            err_obj = {
#                "error": "testing",
#                # "error": f"b64 decode error on file {str(meta_dict.get('objectKey', ''))}: {str(e)}"
#            }
#            print(json.dumps(err_obj), file=sys.stderr, flush=True)
#            continue
#        try:
#            temp = ''
#            # try utf-8
#            temp = bytes_body.decode('utf-8')
#        except Exception:
#            try:
#                # if that fails, try latin-1
#                temp = bytes_body.decode('latin-1')
#            except Exception as e:
#                meta_dict = dict(meta)
#                err_obj = {
#                    "error": "testing",
#                    # "error": f"utf-8/latin-1 decode error on file {str(meta_dict.get('objectKey', ''))}: {str(e)}"
#                }
#                print(json.dumps(err_obj), file=sys.stderr, flush=True)
#                continue
#        text_body = temp
#        try:
#            out_meta = dict(meta)
#            out_meta['error'] = ""
#            object_key = out_meta.get("objectKey", "")
#            file_name = out_meta.get("fileName", "")
#            result = run_spancat(object_key, text_body, file_name, nlp_model)
#            print(json.dumps(result), file=sys.stdout, flush=True)
#
#        except Exception as e:
#            err_obj = {"error": f"nlp processing error: {str(e)}"}
#            print(json.dumps(err_obj), file=sys.stderr, flush=True)
#            continue


if __name__ == "__main__":
    main()
