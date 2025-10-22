
from minio import Minio
from minio.error import S3Error
import psycopg2
from dotenv import load_dotenv
import os
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
    