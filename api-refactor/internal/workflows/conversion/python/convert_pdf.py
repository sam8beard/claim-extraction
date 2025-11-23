import os
import io
import pymupdf
import sys
import json
import base64
import traceback

BODY_DELIMETER = b'--END-BODY--\n'
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
                "error": f"invalid json metadata: {str(e)}",
                "originalKey": meta['objectKey']
            }
            print(json.dumps(err_obj), file=sys.stderr, flush=True)
            continue
        # read b64 body until sentinel
        try:
            leftover = b''
            chunks = []
            while True:
                chunk = sys.stdin.buffer.read(BUF_SIZE)
                if not chunk:
                    # EOF reached unexpectedly
                    break
                # get position of sentinel
                combined = leftover + chunk
                idx = combined.find(BODY_DELIMETER)
                # sentinel is found
                if idx >= 0:
                    # read up until position of sentinel
                    chunks.append(combined[:idx])
                    # consume the rest of the sentinel
                    leftover = combined[idx+len(BODY_DELIMETER):]
                    break

                # say the delimeter is of length N

                # if the ENTIRE sentinel is not found,
                # we still need to keep the last N - 1
                # bytes of the chunk

                # why N - 1 bytes?

                # because this is the max
                # amount of bytes of the delimeter that could be
                # read without actually detecting the entire
                # delimeter

                # consider:
                    # sometexthere--END-BODY-- -> N bytes long -> detected
                    # sometexthere--END-BODY-  -> N - 1 bytes long -> not detected

                    # we save "--END-BODY-"

                    # we prepend on to the next chunk read which contains
                    # the remaining bytes of the delimeter

                    # for example:
                    #   "--END-BODY-" + "-[json meta data]nextfilebodyhere..."
                    #   new chunk = --END-BODY--[json meta data]nextfilebodyhere..."
                    #
                    #   now when we check this chunk, a detection is triggered
                    #
                    #   we then read up to the start of the sentinel
                    #   and append it to our chunks
                    #   (which in this case, is nothing. because the sentinel
                    #    is detected to have started at the first character,
                    #    idx would = 0, and the call to combined[:idx] would return
                    #    and empty list. so nothing gets appended to our total chunks."

                    # in this case, that remaining dash is what we need to trigger a
                    # detection and move forward with the body processing.
                else:
                    # keep the remaining bytes for overlap
                    keep = combined[-(len(BODY_DELIMETER)-1):]
                    chunks.append(combined[:-len(keep)])
                    leftover = keep
            # get b64 string
            body_b64 = b''.join(chunks)
            # decode b64 to bytes
            pdf_bytes = base64.b64decode(body_b64)
        except Exception as e:
            err_obj = {
                "error": f"read/decode body error: {str(e)}",
                "originalKey": meta['objectKey']
            }
            print(json.dumps(err_obj), file=sys.stderr, flush=True)

            # do we also need to find a way to read the rest
            # of stdin after we throw an error in this block????

            # possible issues:

            # characters not in b64 alphabet are being passed in
            # this could cause decode to throw an error

            # properly padded base64 strings have a length that is
            # a multiple of four. so we can check the lengh is valid,
            # and if not, add padding ourselves

            # either way, our base64 string is probably corrupted in some way

            # we know that if this exception happens,
            # it happens because of a call to b64decode
            # which means, either an EOF was hit or the delimeter was found

            # but since the last chunk read is only up until the start
            # of the delimeter, that means the rest of the delimeter would
            # possibly be left in stdin.

            # without assuming any flushing is happening on the
            # go side of things, ( which im pretty sure it isnt,
            # you can't flush stdin ) we need to read the rest of
            # the delimeter and discard it so the pipe is cleared
            # for the next file

            # could be possible the delimeter is being split across two chunks,
            # but this is probably not the reason error is being thrown everytime

            # this should take care of getting rid of the delimeter
            while True:
                chunk = sys.stdin.buffer.read(BUF_SIZE)
                idx = chunk.find(BODY_DELIMETER)
                if idx >= 0:
                    break
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
