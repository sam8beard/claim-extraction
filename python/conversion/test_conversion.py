from minio import Minio
from minio.error import S3Error
import psycopg2
from get_pdfs import get_pdf_files
from extract_text import convert_to_txt
import pprint
def main(): 
    '''

    Test conversion feature

    '''

    # create clients
    minio_client = Minio(
        endpoint="localhost:9000",
        access_key="muel",
        secret_key="password",
        secure=False
    )

    # load env vars before hand
    # eventually replace these with .env vars 
    # ps_client = psycopg2.connect(dbname="mydbname", user="myusername", password="mypassword")

    pdfs = get_pdf_files(minio_client)
    # for name, contents in pdfs.items():
    #     print(name)

    for name, contents in pdfs.items(): 
        convert_to_txt(contents, name)





if __name__ == "__main__": 
    main()