import boto3, os, json, urllib.parse, io, pymupdf
from dotenv import load_dotenv

# load env vars
load_dotenv()

# connect to s3 client
s3 = boto3.client('s3') 

def main(): 
    test_file_download()

# testing boto3 file downloading || THIS WORKS
def test_file_download():
    try:
        # write downloaded file to buffer
        buf = io.BytesIO()
        s3.download_fileobj('claim-pipeline-docstore', 'test/2025-06-24T15:09:44-04:00_basic-text.pdf', buf)

        # extract text
        doc = pymupdf.open('pdf', buf)
        out = ""
        print()

        for page in doc: 
            # retrieve text from a page in pdf, add form feed char
            text = page.get_text().encode('utf-8')
            out += str(text) + "\f"
        
        # encode output string once more and wrap in buf
        text_bytes = out.encode('utf-8')
        buf_out = io.BytesIO(text_bytes)

        # create object key for upload # TO DO — examine naming convention in go cli tool
        

        # upload file
        s3.upload_fileobj('claim-pipeline-docstore', '')
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