import boto3, os, urllib.parse, io, pymupdf, logging
import psycopg2
import sys, json, base64

def main(): 
    '''
    Driver for the script

    '''

    for line in sys.stdin:
        # _ = sys.stdin.read()
            # print(json.dumps)
        try: 
            raw_data = build_data(line)
            json_data = json.dumps(raw_data)
            sys.stdin.write(raw_data)
            sys.stdin.write("\n")
            sys.stdin.flush()
        except Exception as e:
            error_msg = {
                "error": "error processing file"
            }
            json_data = json.dumps(raw_data)
            sys.stderr.write(json_data)
            sys.stderr.write("\n")
            sys.stderr.flush()


        # sys.stdout.flush()
        # sys.stderr.write("Reached convert_to_txt()\n")
        # try: 
        #     data = sys.
        #     output = build_data(line) + "\n"
        #     sys.stdout.write(output)
        #     sys.stdout.flush()
        # except Exception as e: 
        #     error_output = { 
        #     "error": str(e),
        #     }
        #     message = json.dumps(error_output)
        #     sys.stderr.write(message + "\n")
        #     sys.stderr.flush()
    
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
        # decoder = json.JSONDecoder()
        payload = json.loads(line)
        # json.load(line)

        # get values
        title = payload['title']
        body = payload['body']
        old_key = payload['objectKey']
        url = payload['url']

        # process file
        converted_body = str(base64.b64decode(body))
        processed = convert_to_txt(converted_body)
        
        # get new object key
        new_key = get_new_object_key(old_key)

        output = {
            "title": str(title),
            "objectKey": str(new_key),
            "url": str(url), 
            "body": base64.b64encode(processed).encode('utf-8'), # binary bytes encoded 
        }
        # return json object
        return output
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