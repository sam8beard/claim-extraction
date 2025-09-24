import boto3, json 
from pathlib import Path

# connect to s3 client 
s3 = boto3.client('s3')

def pull_s3_files(): 
    
    # open file and read file with keys
    file_path = Path(__file__).parent.parent / "nlp/training/s3-keys.json"
    with open (file_path, "r") as file: 
        keys = json.load(file)

    
    for object_key in keys:
        file_object = s3.get_object(Bucket="claim-pipeline-docstore", Key=object_key)
        # read from object body and decode 
        file_contents = file_object['Body'].read().decode('utf-8')

        # provide texts one at a time for streaming 
        yield file_contents