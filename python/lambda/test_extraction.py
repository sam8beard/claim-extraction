# Test file for extracting text using pymupdf 
import pymupdf as pypdf

doc = pypdf.open("test-files/basic-text.pdf") # open the doc 

out = open("output.txt", "wb")
for page in doc: 
    text = page.get_text().encode("utf-8")
    out.write(text)
    out.write(bytes((12,)))
out.close()
