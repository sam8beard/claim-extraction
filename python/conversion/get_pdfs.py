from minio import Minio
from minio.error import S3Error

def get_pdf_files(client): 

    '''
    Downloads all pdf files from MinIO bucket

    Argument: 
        client: A MinIO client

    Return: 
        A dictionary of object names mapped to their contents
    '''
    assert isinstance(client, Minio), "Must be a properly initialized MinIO client"

    bucket = "claim-pipeline-docstore"
    try: 

        # get all objects that exist under the raw/ directory    
        objects = client.list_objects(bucket_name=bucket, prefix="raw/", recursive=True)

        files = dict()
        for obj in objects: 
            response = None
            try: 
                name = obj.object_name
                response = client.get_object(
                    bucket_name=bucket, 
                    object_name=name
                )
                
            finally: 
                if response:
                    content = response.read()
                    files[name] = content
                    response.close()
                    response.release_conn()
    except Exception as e: 
        print(e)
    return files