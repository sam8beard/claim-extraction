import pymupdf
import os, urllib.parse, io, logging
import sys, json, base64


def main(): 

    for line in sys.stdin: 
        data = do_something(line)
        json_object = json.dumps(data)
        print(json_object)
        sys.stdin.flush()




def do_something(line):
    data = json.loads(line) 
    # ACCOUNT FOR OMITTED FIELDS!!!
    new_name = f"Old name: {data['name']}, New name: Jerry"
    new_job = f"Old job: {data['job']}, New job: Plumber"
    new_data = {"name": new_name, "job": new_job}
    return new_data







if __name__=="__main__":
    main()
