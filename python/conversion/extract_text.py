import boto3, os, json, urllib.parse, io, pymupdf, logging
import psycopg2


def convert_to_txt(file_reader, file_name=None): 
    ''' 
    Extract text from a pdf file.

    Argument: 
        file_reader: A .pdf file reader
        pg_client: A Postgres client

    Return: 
        A bytes file stream representing the converted file
    '''
    # configure logger
    logger = logging.getLogger() 
    logger.setLevel(logging.INFO)

    processed_file = None
    
    try: 
        file_buf = io.BytesIO(file_reader)
    except Exception as e:
        processed_file = e

    # attempt to open pdf with file contents
    try: 
        doc = pymupdf.open("pdf", file_buf)

        # clean up unrecognized xref objects
        for xref in range(1, doc.xref_length()):
            try: 
                _ = doc.xref_object(xref)
            except Exception as e:
                # replace unrecognized xref object with an empty object
                # logger.error(f"\nREPLACING:\t{obj}\n")
                doc.update_object(xref, "<<>>")


    except Exception as e: 
        logger.error(f"\nFailed to open PDF {file_name}: {e}\n")
        processed_file = e
    
    # extract text
    out = ""
    for page_number, page in enumerate(doc, start=1): 
        try: 
            out += page.get_text() + "\f"
        except Exception as e: 
            logger.error(f"\nFailed to extract text from page {page_number} of {file_name}: {e}\n")
            continue
    try: 
        # encode text
        processed_file = out.encode('utf-8')

    except Exception as e: 
        processed_file = e

    assert type(processed_file) == bytes, processed_file.with_traceback()

    return processed_file

def get_new_object_key(raw_name): 
    ''' 
    Creates a new object key for a converted file

    Argument: 
        raw_name: The name of the original object
         
    Returns:
        The new object name

    '''
    assert isinstance(raw_name, str), "Object name must be a string"

    new_name = ""

    try: 
        parent_path = "raw"
        relative_path = os.path.relpath(raw_name, start=parent_path)
        root, ext = os.path.splitext(relative_path)
        new_name = "processed/" + root + ".txt"

    except Exception as e: 
        raw_name = e

    return new_name
