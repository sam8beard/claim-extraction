import boto3, os, urllib.parse, io, pymupdf, logging
import psycopg2
import sys, json, base64

def main(): 
    '''
    Driver for the script

    '''
    for line in sys.stdin:
        try: 
            sys.stdout.write(build_data(line) + "\n") 
            sys.stdout.flush()
        except Exception as e: 
            error_output = { 
            "error": str(e),
            }
            message = json.dumps(error_output)
            sys.stderr.write(message + "\n")
            sys.stderr.flush()
    
def build_data(line): 
    '''
    Build data to stream

    Arguemnt:
        line: a line from stdin

    Returns:
        an encoded json object

    '''
    try:
        # read from stdin
        payload = json.loads(line)
        data = payload['data']

        # get values
        title = data['title']
        body = data['body']
        old_key = data['objectKey']
        url = data['url']

        # process file
        processed = convert_to_txt(base64.b64decode(body))
        
        # get new object key
        new_key = get_new_object_key(old_key)

        output = {
            "title": title,
            "objectKey": new_key,
            "url": url,
            "body": base64.b64encode(processed).decode('utf-8'),
        }
        # encode object
        return json.dumps(output)
    except Exception as e: 
        
        return  e

def convert_to_txt(body): 
    ''' 
    Extract text from a pdf file.

    Argument: 
        body: A .pdf byte stream
        body: An object name corresponding to the reader

    Return: 
        A bytes file stream representing the converted file
    '''
    
    # configure logger
    logger = logging.getLogger() 
    logger.setLevel(logging.INFO)

    processed_file = None
    
    try: 
        file_buf = io.BytesIO(body)
    except Exception as e:
        return e
    
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
        return e
    
    # extract text
    out = ""
    for page_number, page in enumerate(doc, start=1): 
        try: 
            out += page.get_text() + "\f"
        except Exception as e: 
            return e
    try: 
        # encode text
        processed_file = out.encode('utf-8')

    except Exception as e: 
        return e

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

    new_name = ""

    try: 
        parent_path = "raw"
        relative_path = os.path.relpath(raw_name, start=parent_path)
        root, ext = os.path.splitext(relative_path)
        new_name = "processed/" + root + ".txt"

    except Exception as e: 
        return e

    return new_name

if __name__ == "__main__": 
    main()