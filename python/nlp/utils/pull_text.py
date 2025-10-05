import boto3, json 
from pathlib import Path
import unicodedata
import re
# connect to s3 client 
s3 = boto3.client('s3')

def pull_all_files(): 
    
    # open file and read file with keys
    file_path = Path(__file__).parent.parent / "training/s3-keys.json"
 
    with open (file_path, "r") as file: 
        keys = json.load(file)

    
    for object_key in keys:
        file_object = s3.get_object(Bucket="claim-pipeline-docstore", Key=object_key)
        # read from object body and decode 
        file_contents = file_object['Body'].read().decode('utf-8')

        # provide texts one at a time for streaming 
        # yield preprocess_text(file_contents)
        yield preprocess_text(file_contents)

# for faster testing
def pull_one_file():
    # open file and read file with keys
    file_path = Path(__file__).parent.parent / "training/s3-keys.json"
 
    with open (file_path, "r") as file: 
        keys = json.load(file)

    object_key = keys[0]
    
    file_object = s3.get_object(Bucket="claim-pipeline-docstore", Key=object_key)
    # read from object body and decode 
    file_contents = file_object['Body'].read().decode('utf-8')

    return preprocess_text(file_contents)

# normalize text 
def preprocess_text(text):
    text = unicodedata.normalize("NFC", text)
    text = "".join(c for c in text if c.isprintable())
    text = re.sub(r"\s+", " ", text)
    text = re.sub(r"([\-=_])\1{2,}", "", text)
    text = re.sub(r"[\x0c\x0b]", "", text)
    text = text.strip()
    return text 

# def preprocess_text(text_list): 
#     for text in text_list: 
#         text = unicodedata.normalize("NFC", text)
#         text = "".join(c for c in text if c.isprintable())
#         text = re.sub(r"\s+", " ", text)
#         text = text.strip()
#         yield text 