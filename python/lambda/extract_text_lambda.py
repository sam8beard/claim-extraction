import boto3, os, json, urllib.parse
from dotenv import load_dotenv

# load env vars
load_dotenv()

# connect to s3 client
s3 = boto3.client('s3') 

def main(): 
    test_file_download()

# testing boto3 file downloading 
def test_file_download(): 
    try:
        s3.download_file('claim-pipeline-docstore', 'test/2025-06-24T15:09:44-04:00_basic-text.pdf', 'output.txt')
    except Exception as e: 
        print(e)

# main handler for extraction lambda
def handler(event, context): 
    # get an object from the event
    bucket = event['Records'][0]['s3']['bucket']['name']
    key = urllib.parse.unquote_plus(event['Records'][0]['s3']['object']['key'], encoding='utf-8')

    print(bucket, key)


if __name__ == '__main__': 
    main()