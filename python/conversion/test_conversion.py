from minio import Minio
from minio.error import S3Error
import psycopg2
import extract_text
import get_pdfs
import upload_text
import os
import pprint
from dotenv import load_dotenv

def create_clients(): 
    '''
    Loads environment variables and creates MinIO client 
    and PostgreSQL client.

    Return: 
        - minio_client: MinIO client
        - pg_client: PostgreSQL client
    '''

    minio_client, pg_client = None, None

    try: 
        # load env vars
        load_dotenv()

        # create clients
        db_params = { 
            "host": os.getenv("DB_HOST"),
            "database": os.getenv("DB_NAME"),
            "user": os.getenv("DB_USERNAME"), 
            "password": os.getenv("DB_PASSWORD"), 
            "port": os.getenv("DB_PORT")
        }
        minio_params = {
            "endpoint": os.getenv("MINIO_ENDPOINT"),
            "access_key": os.getenv("MINIO_USER"),
            "secret_key": os.getenv("MINIO_PASSWORD"),
            "secure": False
        }
       
        pg_client = psycopg2.connect(**db_params)
        minio_client = Minio(**minio_params)
       

    except Exception as e: 
        minio_client, ps_client = e, e
        traceback = e.with_traceback()
    
    minio_success = not isinstance(minio_client, Exception)
    pg_success = not isinstance(pg_client, Exception) 

    assert minio_success and pg_success, traceback

    return minio_client, pg_client
    
def main(): 
    '''

    Test conversion feature

    '''

    minio_client, pg_client = create_clients()

    pdfs = get_pdfs.get_pdf_files(minio_client)

    for name, content in pdfs.items(): 
        # extract text 
        text_content = extract_text.convert_to_txt(name, content)
        # create new object key
        new_name = extract_text.get_new_object_key(name)

        # upload file
        file_dict = {new_name: text_content}
        upload = upload_text.upload_text_files(minio_client, file_dict)
        
        # original object name and upload result
        extraction_state = (name, upload[new_name])
        
        # log extraction state in PG documents table 
        success = upload_text.log_extraction_state(pg_client, extraction_state)
        
        print(success)

        # if success: 
        #     print(success)


    pg_client.close()
    
    



if __name__ == "__main__": 
    main()