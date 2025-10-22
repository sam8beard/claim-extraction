from minio import Minio
from minio.error import S3Error
import psycopg2
from utils import extract_text
from utils import get_pdfs
from utils import upload_text
from utils import create_clients
import os
import pprint
from dotenv import load_dotenv

def main(): 
    '''

    Test conversion feature

    '''

    minio_client, pg_client = create_clients.create_clients()

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



    pg_client.close()
    
    



if __name__ == "__main__": 
    main()