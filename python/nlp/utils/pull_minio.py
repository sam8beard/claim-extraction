from minio import Minio
from minio.error import S3Error
import requests
import json
import os
from dotenv import load_dotenv
import re
import unicodedata

def get_keys():
    ''' 
    Sends a GET request to Go API to retrieve 
    properly extracted keys. 

    Return: 
        - keys: list of key strings
    '''
    url = "http://localhost:60000/documents/extracted"

    try: 

        resp = requests.get(url)
        if resp.status_code == 200: 
            data = resp.json()
            if data: 
                keys = data['keys']
            else: 
                return "Response is empty"
            # print(data)
            
    except Exception as e: 
        return e
    
    return keys

def pull_n_files(keys, n=1): 
    '''
    Parses keys and downloads files from MinIO. 
    Returns the files as strings.

    Argument: 
        - keys: file keys
        - n: number of files to pull
    Yield: 
        - file: file string
    '''
    
    # load env vars
    load_dotenv()

    # create MinIO client
    minio_params = {
            "endpoint": os.getenv("MINIO_ENDPOINT"),
            "access_key": os.getenv("MINIO_USER"),
            "secret_key": os.getenv("MINIO_PASSWORD"),
            "secure": False
        }
    
    client = Minio(**minio_params)
    bucket = "claim-pipeline-docstore"

    found_count = 0
    for i in range(min(n, len(keys))): 
        name = keys[i]
        print(name)
        # NOTE: i think theres an issue with the .txt keys 
        #       and the way i named them
        response = None
        try:
            response = client.get_object(
                bucket_name=bucket,
                object_name=name
            )
            file = response.read().decode('utf-8')
            response.close()
            response.release_conn()
            # print(f"\nFOUND\n")
            found_count += 1
            yield preprocess_text(file)

        except Exception as e: 
            print(f"\nFile does not exist in bucket: {name}\n")
            continue
    # print(len(keys))
    # print(found_count)
def preprocess_text(text):
    ''' 
    Preprocess text for model training

    Argument: 
        - text: the text to normalize
    Return: 
        - text: the normalized text 
    '''
    text = unicodedata.normalize("NFC", text)
    text = "".join(c for c in text if c.isprintable())
    text = re.sub(r"\s+", " ", text)
    # text = re.sub(r"([\-=_])\1{2,}", "", text)
    # text = re.sub(r"[\x0c\x0b]", "", text)
    text = text.strip()
    return text 

keys = get_keys()
# for key in keys: 
#     print(f"\n{key}\n")
for file in pull_n_files(keys, 2): 
    print(file[3000:])
    
