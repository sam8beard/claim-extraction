import os, urllib.parse, io, pymupdf, logging
import sys, json, base64
import traceback
def main(): 
    '''

    Driver for the script

    '''
    '''
    New logic: 

    read until newline to get metadata, then read until sentinel for body 
    '''
    BODY_DELIMITER = b"--END-BODY--\n" while True:
        meta_line = sys.stdin.readline()
        if not meta_line: 
            break # EOF
        # build metadata json
        metadata = build_data(meta_line)

        # process body chunk by chunk
        chunks = []
        while True: 
            chunk = sys.stdin.buffer.read(4096) 
            if not chunk: 
                break # EOF
            # append delimeter
            if BODY_DELIMETER in chunk: 
                idx = chunk.find(BODY_DELIMETER)
                chunks.append(chunk[:idx])
                break
            chunks.append(chunk)
        body_bytes = b"".join(chunks)
        converted_bytes = convert_to_txt(b64.base64decode(body_bytes))

        if isinstance(metadata, Exception) or isinstance(converted_bytes, Exception):
        # write bytes to stderr
            exception_json = {"error": "Unable to process file"}
            print(json.dumps(exception_json), file=sys.stderr, flush=True)
        else:
        # write bytes to stdout
        # write metadata
            print(metadata, file=sys.stdout, flush=True)
            # write body in chunks with delimeter at end
            converted_bytes = b64.base64encode(converted_bytes)
            reader = io.BytesIO(converted_bytes)
            # write chunk by chunk and append write delimeter at end
            while True: 
                chunk = reader.read(4096) 
                if not chunk: 
                    break # EOF
                sys.stdout.buffer.write(chunk)
            sys.stdout.buffer.write(BODY_DELIMETER)
            sys.stdout.buffer.flush()

            ''' 

            Test with this in extractor.go:

            JVBERi0xLjQKJeLjz9MKMSAwIG9iago8PC9UeXBlIC9DYXRhbG9n
            L1BhZ2VzIDIgMCBSCj4+CmVuZG9iagoyIDAgb2JqCjw8L1R5cGUg
            L1BhZ2VzL0tpZHNbMyAwIFJdL0NvdW50IDEKPj4KZW5kb2JqCjMg
            MCBvYmoKPDwvVHlwZSAvUGFnZS9NZWRpYUJveCBbMCAwIDYxMiA3
            OTJdL0NvbnRlbnRzIDQgMCBSCi9SZXNvdXJjZXMgPDwvRm9udCA8
            PC9GMSA1IDAgUj4+Pj4+Pj4KZW5kb2JqCjQgMCBvYmoKPDwvTGVu
            Z3RoIDQ5Pj4Kc3RyZWFtCkJUIC9GMSAxMiBUZgoKSGVsbG8sIFdv
            cmxkIQplbmRzdHJlYW0KZW5kb2JqCjUgMCBvYmoKPDwvVHlwZSAv
            Rm9udC9TdWJ0eXBlIC9UeXBlMQpOYW1lIC9GMQpCYXNlRm9udCAv
            SGVsdmV0aWNhPj4KZW5kb2JqCnhyZWYKMCA2CjAwMDAwMDAwMDAg
            NjU1MzUgZiAKMDAwMDAwMDAxMCAwMDAwMCBuIAowMDAwMDAwMDYx
            IDAwMDAwIG4gCjAwMDAwMDAxMjAgMDAwMDAgbiAKMDAwMDAwMDM0
            MCAwMDAwMCBuIAp0cmFpbGVyCjw8L1NpemUgNi9Sb290IDEgMCBS
            L0luZm8gNiAwIFI+PgpzdGFydHhyZWYKNDA0CiUlRU9G
            '''

    #    if isinstance(result, Exception) or isinstance(converted_bytes, Exception):
    #        # write to stderr
    #    else: 
    #        # write to stdout 
    #        print(metadata, file=sys.stdout, flush=True)
    #        converted_chunks = []
    #        while True: 
    #            chunk = converted_bytes.read(4096)
    #            if not chunk: 
    #                break # EOF
    #            if BODY_SENTINEL in chunk 
    
    # for line in sys.stdin:
    #         # print(json.dumps)
    #     
    #     result = build_data(line)
    #     '''
    #     Will need to switch sys.stdout.write() 
    #     with print(). For some reason, I couldn't get 
    #     write() to work, maybe it has something to do with me 
    #     making the output unbuffered? Make sure to confirm 
    #     that print() will work for all edge cases, will have to 
    #     seek solution for printing out to stderr. 
    #     
    #     NOTE: 
    #     Use print() if sys.std isnt working 
    #     '''
    #     if isinstance(result, Exception): 
    #         exec_type, exec_value, trace = sys.exc_info()
    #         tb_msg = "".join(traceback.format_tb(trace))
    #         exc_msg = f"{tb_msg}"
    #         # exec_info = [exec_type.__name__, exec_value, tb_msg]
    #         # exception_msg = "\n".join(exec_info)
    #         exception_json = {"error": "firing"}
    #         exception_json = json.dumps(exception_json)
    #         # sys.stderr.write(exception_json + "\n")
    #         # sys.stderr.flush()
    #         print(exception_json, file=sys.stderr, flush=True)
    #     else:
    #         # sys.stdout.write(result + "\n") 
    #         # sys.stdout.flush()
    #         print(result, file=sys.stdout, flush=True)
       
def build_data(line): 
    '''
    Build data to stream

    Arguemnt:
        line: a line from stdin

    Returns:
        an encoded json object

    '''
    try:
        payload = json.loads(line)

        # get values
        # need to use loop to avoid omitted fields 
        key_list = ['title', 'objectKey', 'body', 'url', 'error', 'originalKey']
        output = dict.fromkeys(key_list)
        for key, value in payload.items(): 
            # if the key is body or objectKey, we need to process them
            # otherwise, keep same value 
            if key == "body":
                decoded_body = base64.b64decode(value)
                processed_body = convert_to_txt(decoded_body)
                # rencode to b64 and then into utf-8
                rencoded_body = base64.b64encode(processed_body).decode('utf-8')
                output[key] = rencoded_body
            elif key == "objectKey": 
                new_key = get_new_object_key(value)
                output[key] = new_key
                # we know there will initially be no errors coming in
            elif key == "error": 
                output[key] = "none" 
            else: 
                output[key] = value
        # return json object
        output = json.dumps(output)
        return output
    except Exception as e: 
       return e 

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
