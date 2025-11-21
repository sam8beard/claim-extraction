import os
import io
import pymupdf
import sys
import json
import base64
import traceback

BODY_DELIMETER = b"--END-BODY--\n"
BUF_SIZE = 4096


def main():
    '''

    Driver for the script

    '''
    for meta_line in sys.stdin:
        meta_line = meta_line.rstrip("\n")
        if not meta_line:
            continue

        # parse json metadata
        try:
            meta = json.loads(meta_line)
        except Exception as e:
            err_obj = {
                "error": f"invalid json metadata: {str(e)}"
            }
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue
        # read b64 body until sentinel
        try:
            chunks = []
            while True:
                chunk = sys.stdin.buffer.read(BUF_SIZE)
                if not chunk:
                    # EOF reached unexpectedly
                    break
                # get position of sentinel
                idx = chunk.find(BODY_DELIMETER)
                # sentinel is found
                if idx >= 0:
                    # read up until position of sentinel and exit read loop
                    chunks.append(chunk[:idx])
                    break
                # otherwise, add full chunk
                chunks.append(chunk)
            # get b64 string
            body_b64 = b"".join(chunks)
            # decode b64 to bytes
            pdf_bytes = base64.b64decode(body_b64)
        except Exception as e:
            err_obj = {"error": f"read/decode body error: {str(e)}"}
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue
        # convert pdf bytes to text
        try:
            converted = convert_to_txt(pdf_bytes)
            # encode converted bytes
            converted_b64 = base64.b64encode(converted)
        except Exception as e:
            tb = traceback.format_exc()
            err_obj = {"error": f"convert error: {str(e)}", "trace": tb}
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue
        # prepare metadata for return
        try:
            out_meta = dict(meta)
            # make sure body field exists
            out_meta['body'] = ""
            out_json = json.dumps(out_meta)
        except Exception as e:
            err_obj = {"error": f"metadata build error: {str(e)}"}
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue
        # write metadata line
        print(out_json, file=sys.stdout, flush=True)
        # stream converted bytes to buffer, then the sentinel
        try:
            reader = io.BytesIO(converted_b64)
            while True:
                chunk = reader.read(BUF_SIZE)
                if not chunk:
                    break
                sys.stdout.buffer.write(chunk)
            sys.stdout.buffer.write(BODY_DELIMETER)
            sys.stdout.buffer.flush()
        except Exception as e:
            err_obj = {"error": f"write output error: {str(e)}"}
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue
#
#    for line in sys.stdin:
#        # print(json.dumps)
#        result = build_data(line)
#        '''
#        Will need to switch sys.stdout.write()
#        with print(). For some reason, I couldn't get
#        write() to work, maybe it has something to do with me
#        making the output unbuffered? Make sure to confirm
#        that print() will work for all edge cases, will have to
#        seek solution for printing out to stderr.
#        NOTE:
#        Use print() if sys.std isnt working
#        '''
#        if isinstance(result, Exception):
#            exec_type, exec_value, trace = sys.exc_info()
#            tb_msg = "".join(traceback.format_tb(trace))
#            exc_msg = f"{tb_msg}"
#            # exec_info = [exec_type.__name__, exec_value, tb_msg]
#            # exception_msg = "\n".join(exec_info)
#            exception_json = {"error": "firing"}
#            exception_json = json.dumps(exception_json)
#            # sys.stderr.write(exception_json + "\n")
#            # sys.stderr.flush()
#            print(exception_json, file=sys.stderr, flush=True)
#        else:
#            # sys.stdout.write(result + "\n")
#            # sys.stdout.flush()
#            print(result, file=sys.stdout, flush=True)


# def build_data(line):
#    '''
#    Build data to stream
#
#    Argument:
#        line: a line from stdin
#
#    Returns:
#        an encoded json object
#
#    '''
#    try:
#        payload = json.loads(line)
#
#        # get values
#        # need to use loop to avoid omitted fields
#        key_list = ['title', 'objectKey', 'body',
#                    'url', 'error', 'originalKey']
#        output = dict.fromkeys(key_list)
#        for key, value in payload.items():
#            # if the key is body or objectKey, we need to process them
#            # otherwise, keep same value
#            if key == "body":
#                decoded_body = base64.b64decode(value)
#                processed_body = convert_to_txt(decoded_body)
#                # rencode to b64 and then into utf-8
#                rencoded_body = base64.b64encode(
#                    processed_body).decode('utf-8')
#                output[key] = rencoded_body
#            elif key == "objectKey":
#                new_key = get_new_object_key(value)
#                output[key] = new_key
#                # we know there will initially be no errors coming in
#            elif key == "error":
#                output[key] = "none"
#            else:
#                output[key] = value
#        # return json object
#        output = json.dumps(output)
#        return output
#    except Exception as e:
#        return e


def convert_to_txt(body):
    '''
    Extract text from a pdf file.

    Argument:
        body: A .pdf byte stream
        body: An object name corresponding to the reader

    Return:
        A bytes file stream representing the converted file
    '''

    processed_file = None

    try:
        file_buf = io.BytesIO(body)
    except Exception as e:
        raise RuntimeError(f"could not create buffer for pdf: {e}")

    # attempt to open pdf with file contents
    try:
        doc = pymupdf.open("pdf", file_buf)

        # clean up unrecognized xref objects
        for xref in range(1, doc.xref_length()):
            try:
                _ = doc.xref_object(xref)
            except Exception:
                # replace unrecognized xref object with an empty object
                doc.update_object(xref, "<<>>")

    except Exception as e:
        raise RuntimeError(f"open pdf failed: {e}")

    # extract text
    out = ""
    for page_number, page in enumerate(doc, start=1):
        try:
            out += page.get_text() + "\f"
        except Exception as e:
            raise RuntimeError(f"could not get text from page: {e}")
    try:
        # encode text
        processed_file = out.encode('utf-8')

    except Exception as e:
        raise RuntimeError(f"could not encode text bytes: {e}")

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
