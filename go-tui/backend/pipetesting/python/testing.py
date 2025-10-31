import pymupdf
import os, urllib.parse, io, logging
import sys, json, base64


def main(): 

    for line in sys.stdin: 
        data = process(line)
        json_object = json.dumps(data)
        print(json_object)
        sys.stdin.flush()




def process(line):
    data = json.loads(line) 
    new_data = ""
    # ACCOUNT FOR OMITTED FIELDS!!!
    for key, value in data.items(): 
        if key == "name": 
            new_name = value.upper()
            new_data = {"name": new_name, "id": data["id"] }
            return new_data







if __name__=="__main__":
    main()
