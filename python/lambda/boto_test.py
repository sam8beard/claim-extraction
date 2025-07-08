import boto3, os
from dotenv import load_dotenv

load_dotenv()

s3 = boto3.resource('s3')

for bucket in s3.buckets.all(): 
    print(bucket.name)