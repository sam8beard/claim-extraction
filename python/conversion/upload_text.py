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

    successful_uploads = dict()

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

            if result: successful_uploads[name] = True

        except Exception as e: 
            successful_uploads[name] = False

    
    return successful_uploads
        
def log_extraction_state(client, uploads): 
    ''' 
    Changes all Postgres documents.text_extracted columns to True
    for files with successfully extracted text

    Argument: 
        client: A PsycoPG client
        uploads: A dictionary of files mapped to a boolean representing 
                 their extraction state
        
    '''


