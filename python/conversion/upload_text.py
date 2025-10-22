from minio import Minio
from minio.error import S3Error
import io
import psycopg2

def upload_text_files(client, files):

    '''
    Uploads text files to MinIO bucket under the 'processed/'
    directory 

    Argument: 
        client: A MinIO client
        files: A dictionary of object names mapped to contents
    
    Return: 
        A dictionary of object names mapped to a boolean representing 
        their upload success
    '''

    assert isinstance(client, Minio), "Must be a properly initialized MinIO client"
    assert len(files) > 0, "Dictionary must have at least one key-value pair"

    uploads = dict()

    bucket_name = "claim-pipeline-docstore"
    success = False

    for name, content in files.items(): 
        file_reader = io.BytesIO(content)
        content_length = len(content)
        try: 
            result = client.put_object(
                bucket_name=bucket_name, 
                object_name=name,
                data=file_reader,
                length=content_length
            )

            if result: uploads[name] = True

        except Exception as e: 
            uploads[name] = False

    return uploads
        
def log_extraction_state(client, uploads): 
    ''' 
    Changes all Postgres documents.text_extracted columns to True
    for files with successfully extracted text

    Argument: 
        client: A PostgreSQL client
        uploads: A tuple containing the raw object name and a boolean representing 
                 its extraction state
        
    '''

    # example of query needed to get extracted keys 
    # SELECT s3_key FROM documents WHERE text_extracted=true ORDER BY uploaded_at

    query = """
        UPDATE documents
        SET text_extracted = %s
        WHERE s3_key = %s
    """
    try: 
        with client.cursor() as cur: 
            
            success = False
            
            name, extracted = uploads
            # print(uploads)
            if extracted: 
                # THIS NEEDS TO BE THE RAW NAME, NOT PROCESSED
                cur.execute(query, (extracted, name))
                success = f"\nModified documents.text_extracted for: \n{name}\n"
            
                client.commit()
            else: 
                success = f"\nCould not modify row:\t{name}\n{extracted}\n"
    except Exception as e: 
        print(f"\nDB error: {e}\n")
        sucess = f"\nDB error: {e}\n"
        client.rollback()

   
        

    return success