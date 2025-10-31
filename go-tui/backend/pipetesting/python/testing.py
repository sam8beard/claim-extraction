import pymupdf
import os, urllib.parse, io, logging
import sys, json, base64


def main(): 

    for line in sys.stdin: 
        do_something(line)





def do_something(line):
    sys.stdout.write(line) 
    sys.stdout.flush()







if __name__=="__main__":
    main()
