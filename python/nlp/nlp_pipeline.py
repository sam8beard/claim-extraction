import spacy 
from spacy.matcher import PhraseMatcher
from spacy.tokens import Span
from spacy.pipeline import EntityRuler
from utils.pull_text import pull_s3_files
nlp = spacy.load("en_core_web_trf")

#  custom NER entities and labels
artint_labels = [
    {"label": "ARTINT", "pattern": "Artificial Intelligence"},
    {"label": "ARTINT", "pattern": "A.I."},
    {"label": "ARTINT", "pattern": "AI"}
]

def main():


   

 


    # add Artifical Intelligence and AI (maybe others) as entities using Entity Ruler
    ruler = EntityRuler(nlp)

    # 

    # Will probably make helper function for adding custom entities 
    # and then build doc after

    #
    
       
    
    ruler = nlp.add_pipe("entity_ruler", before="ner")
    # ruler.add_patterns(artint_labels)


    # print(ai_patterns)
    # matcher.add("ARTINT", ai_patterns)
    # matches = matcher(doc)

    # new_ents = []

    # for match_id, start, end in matches: 
    #     new_ents.append(doc[start:end])
       
    # doc.set_ents(new_ents)
    
    #  # sample data 
    # doc = nlp(
    #         "Artificial Intelligence (AI) is revolutionizing the healthcare industry by enhancing diagnostic accuracy, personalizing treatment plans, and streamlining administrative tasks. "
    #         "IBM Watson Health and Google Health are leading this transformation. "
    #         "Machine learning algorithms analyze vast amounts of medical data from institutions like Mayo Clinic and Johns Hopkins Hospital to identify patterns that may be missed by human clinicians. "
    #         "For instance, AI-driven imaging tools such as Aidoc Radiology can detect early signs of diseases such as lung cancer and Alzheimer's disease, enabling timely interventions. "
    #         "In 2025, these technologies are expected to be integrated into over 60 percent of hospitals worldwide."
    # )
    # checking for ARTINT entities 
    # for ent in doc.ents: 
    #     print(ent.text, ent.label_)

    # validate tokenization, segmentation, and named entity recognition
    # print("Sentences: \n")
    # for sent in doc.sents: 
    #     print(sent, "\n")
    # print("Token details: \n")
    # for token in doc: 
    #     print(token.text, token.lemma_, token.pos_, "\n")
    # print("Named entities: \n")
    # for ent in doc.ents: 
    #     print(ent.text, ent.label_, "\n")
    # print(nlp.get_pipe("ner").labels)

    
    # pull files, convert to text, convert to string list
    
    for doc in nlp.pipe(pull_s3_files()):

        # validate tokenization, segmentation, and named entity recognition
        # print("Sentences: \n")
        # for sent in doc.sents: 
            # print(sent, "\n")
        # print("Token details: \n")
        # for token in doc: 
        #     print(token.text, token.lemma_, token.pos_, "\n")
        print("Named entities: \n")
        for ent in doc.ents: 
            print(ent.text, ent.label_, "\n")


        # print(doc)

    # pass all strings into pipeline and create a doc list 

if __name__ == "__main__": 
    main() 