import os
import io
import pymupdf
import sys
import json
import base64
import traceback

BODY_DELIMETER = b'--END-BODY--\n'
BUF_SIZE = 4096


def process_bodies(bodies):
    results = dict()
    for key, file in bodies.items():
        body_b64 = file.get('body_b64')
        meta = file.get('meta')
        try:
            # decode b64 to bytes
            missing_padding = len(body_b64) % 4
            if missing_padding:

                pdf_bytes = base64.b64decode(body_b64 + b'==')
            else:
                pdf_bytes = base64.b64decode(body_b64)
        except Exception as e:
            err_obj = {
                "error": f"read/decode body error: {str(e)}",
                "originalKey": meta['objectKey']
            }
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
            # add new object key to payload
            old_key = out_meta['originalKey']
            try:
                new_key = get_new_object_key(old_key)
                out_meta['objectKey'] = new_key
                out_json = json.dumps(out_meta)
            except Exception as e:
                err_obj = {
                    "error": f"could not build new object key: {str(e)}"
                }
                print(json.dumps(err_obj), file=sys.stderr, flush=True)
                continue
        except Exception as e:
            err_obj = {"error": f"metadata build error: {str(e)}"}
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue

        results[key] = {"converted_b64": converted_b64, "out_json": out_json}

    return results


def main():
    '''

    Driver for the script

    '''
    for meta_line in sys.stdin:
        meta_line = meta_line.rstrip("\n")
        if not meta_line:
            continue

        try:
            meta = json.loads(meta_line)
        except Exception as e:
            err_obj = {"error": f"invalid json metadata: {str(e)}"}
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue

        # Read body
        try:
            leftover = b''
            chunks = []
            while True:
                chunk = sys.stdin.buffer.read(BUF_SIZE)
                if not chunk:
                    break
                combined = leftover + chunk
                idx = combined.find(BODY_DELIMETER)
                if idx >= 0:
                    chunks.append(combined[:idx])
                    leftover = combined[idx+len(BODY_DELIMETER):]
                    break
                else:
                    keep = combined[-(len(BODY_DELIMETER)-1):]
                    chunks.append(combined[:-len(keep)])
                    leftover = keep

            body_b64 = b''.join(chunks)
        except Exception as e:
            err_obj = {"error": f"error reading body: {str(e)}"}
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue

        # Process immediately
        try:
            missing_padding = len(body_b64) % 4
            if missing_padding:
                pdf_bytes = base64.b64decode(body_b64 + b'==')
            else:
                pdf_bytes = base64.b64decode(body_b64)

            converted = convert_to_txt(pdf_bytes)
            converted_b64 = base64.b64encode(converted)

            # Prepare output metadata
            out_meta = dict(meta)
            out_meta['body'] = ""
            old_key = out_meta.get('originalKey', meta.get('objectKey'))
            new_key = get_new_object_key(old_key)
            out_meta['objectKey'] = new_key
            out_json = json.dumps(out_meta)

            # Write output immediately
            print(out_json, file=sys.stdout, flush=True)
            sys.stdout.buffer.write(converted_b64)
            sys.stdout.buffer.write(BODY_DELIMETER)
            sys.stdout.buffer.flush()

        except Exception as e:
            tb = traceback.format_exc()
            err_obj = {
                "error": f"processing error: {str(e)}", "trace": tb, "originalKey": meta.get('objectKey')}
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue


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
        raise e

    return new_name


if __name__ == "__main__":
    main()
