import boto3, os, json, urllib.parse, io, pymupdf, logging

# configure logger
logger = logging.getLogger() 
logger.setLevel(logging.INFO)

# connect to s3 client
s3 = boto3.client('s3') 

# main handler for extraction lambda
def process_handler(event, context): 

    # get an object from the event
    bucket_name = event['Records'][0]['s3']['bucket']['name']
    object_key = urllib.parse.unquote_plus(event['Records'][0]['s3']['object']['key'], encoding='utf-8')

    # ERROR CHECKING: check if object key and bucket name are properly received 
    logger.info(f"Processing file: {object_key} in bucket: {bucket_name}")


    # now we have the bucket and object key, lets retrieve the object 
    resp = s3.get_object(Bucket=bucket_name, Key=object_key)

    # ERROR CHECKING: check if file is properly downloaded 
    file_size = resp['ContentLength']
    logger.info(f"Downloaded file of size: {file_size} bytes")

    # wrap file body in buffer
    file_body_buf = io.BytesIO(resp['Body'].read())

    try: 
        doc = pymupdf.open("pdf", file_body_buf)
    except Exception as e: 
        logger.error(f"Failed to open PDF {object_key}: {e}")
        return
    
    # extract text from pdf
    out = ""
    for page_number, page in enumerate(doc, start=1): 
        try: 
            out += page.get_text() + "\f"
        except Exception as e: 
            logger.error(f"Failed to extract text from page {page_number} of {object_key}: {e}")
            return
        
    # ERROR CHECKING: log length of extracted text 
    logger.info(f"Extracted text length: {len(out)} characters")

    # encode output string
    processed_file = out.encode('utf-8')

    # create new object key
    parent_path = "raw"
    relative_path = os.path.relpath(object_key, start=parent_path)
    root, ext = os.path.splitext(relative_path)
    new_object_key = "processed/" + root + ".txt"

    # ERROR CHECKING: check to see if new object key is properly constructed 
    logger.info(f"Uploading processed text to: {new_object_key}")

    try: 
        # add processed file to bucket 
        s3.put_object(Body=processed_file, Bucket=bucket_name, Key=new_object_key)
    except Exception as e: 
        logger.error(f"Failed to upload processed file: {e}")
        return 