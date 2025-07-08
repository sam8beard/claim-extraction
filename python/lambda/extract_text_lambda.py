import boto3, os, json, urllib.parse
from dotenv import load_dotenv

# load env vars
load_dotenv()

# connect to s3 client
# s3 = boto3.client('s3') # COMMENTED OUT FOR NOW WHILE TESTING

# FOR TESTING: load sample event for parsing

def main(): 
    test_file = open('sample_event.json', 'r')
    data = json.load(test_file)
    handler(data, None)

def handler(event, context): 
    # get an object from the event
    bucket = event['Records'][0]['s3']['bucket']['name']
    key = urllib.parse.unquote_plus(event['Records'][0]['s3']['object']['key'], encoding='utf-8')

    print(bucket, key)

if __name__ == '__main__': 
    main()