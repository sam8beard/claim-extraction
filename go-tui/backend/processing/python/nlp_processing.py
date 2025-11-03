import sys
import json
import re
import unicodedata
import spacy
import numpy as np
from collections import Counter

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
        return 0.0 # no spans, no score

    # get proportions for each type 
    proportions = [counts[t] / total for t in types]

    # the more even proportion, the lower variance, the higher score
    mean = sum(proportions) / len(types) 
    variance = sum((p - mean) **2 for p in proportions) / len(types)

    score = 1 - variance 
    return round(score, 3) 

def run_spancat(object_key, raw_text, nlp_model): 
    text = preprocess_text(raw_text)
    file_data = {
        "objectKey": object_key,
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

def test(): 
    test_files = [
 {
        "objectKey": "file1.txt",
        "content": (
            "Dr. Smith claims that the vaccine is highly effective. "
            "However, recent studies suggest otherwise. "
            "The WHO confirms that further testing is required. "
            "Experts debate the methodology used in these studies. "
            "Ultimately, the evidence remains inconclusive."
        )
    },
    {
        "objectKey": "file2.txt",
        "content": (
            "Alice asserts that the new policy will reduce emissions. "
            "Bob counters that the economic impact will be severe. "
            "The government releases official figures supporting Alice's claim. "
            "Environmental groups respond positively, emphasizing long-term benefits. "
            "Analysts remain skeptical about the short-term effects."
        )
    },
    {
        "objectKey": "file3.txt",
        "content": (
            "The article reports that inflation has risen by 3.2% over the past quarter. "
            "Economists warn that this trend may continue if interest rates are not adjusted. "
            "Consumer groups note the rising cost of living as evidence. "
            "Some banks argue that this increase is within expected limits."
        )
    },
    {
        "objectKey": "file4.txt",
        "content": (
            "NASA confirms that the satellite has successfully entered orbit. "
            "Scientists observe minor deviations in trajectory, which are under investigation. "
            "Independent analysts suggest that the mission's success will advance space research. "
            "The agency releases images and data for public review."
        )
    },
    {
        "objectKey": "file5.txt",
        "content": (
            "According to several reports, the technology startup achieved record revenue last year. "
            "Investors praise the management team for strategic decisions. "
            "Competitors question the sustainability of such growth. "
            "Financial analysts provide a cautious outlook for the upcoming fiscal year."
        )
    }
    ]

    for f in test_files: 
        json_line = json.dumps(f)
        result = test_main(json_line)

        if isinstance(result, Exception): 
            print("Error processing")
        else: 
            print(json.dumps(result, indent=2))

def test_main(json_line): 
    nlp_model = spacy.load("spancat_v1.0")
    nlp_model.add_pipe("sentencizer", before="spancat")
    json_line = json_line.strip()
    
    try: 
        obj = json.loads(json_line)
        object_key = obj.get("objectKey", "")
        content = obj.get("content", "")
        result = run_spancat(object_key, content, nlp_model)
        return result

     # handle stderr here 
    except Exception as e:
        return e

def main(): 
    nlp_model = spacy.load("spancat_v1.0")
    nlp_model.add_pipe("sentencizer", before="spancat")

    for line in sys.stdin: 
        line = line.strip()
        if not line:
            continue
        try: 
            obj = json.loads(line)
            object_key = obj.get("objectKey", "")
            content = obj.get("content", "")
            result = run_spancat(object_key, content, nlp_model)
            print(json.dumps(result), flush=True)

        # handle stderr here 
        except:
            print("error")
if __name__ == "__main__":
    test()
