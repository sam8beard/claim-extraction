from utils.pull_text import pull_s3_files

def get_training(): 
    
    for text in pull_s3_files(): 
        print(text)
        print("\n-------------------------\n")
get_training()