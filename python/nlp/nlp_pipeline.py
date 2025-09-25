import spacy 
from spacy.matcher import PhraseMatcher
from spacy.tokens import Span
from spacy.pipeline import EntityRuler
from utils.pull_text import pull_s3_files

# load model
nlp = spacy.load("en_core_web_trf")

# START WITH JUST THESE THREE ENTITIES, EXPAND OVER TIME
#  custom NER entities and labels

# 

# Eventually implement EntityRuler.initialize to read patterns from disk or memory
# https://spacy.io/api/entityruler
# other helpful functions in here too
#


# custom entities and labels for custom entity recognition 
custom_ents = [


# FIX! - remove other entities that also have AI in label
# ARTINT - artificial intelligence in any way it is mentioned (recognized by ner as ORG)
    # need to edit to avoid other entities that have AI in name
    # maybe check if its already labeled ORG? i think this only works "AI"
    # it seems that only mentions of ai that are labeled as ORG 
    # are not part of any kind of institution name

    # need to prevent override this causes when a default ner entity uses AI in its label
    # e.g. Institue for Artificial Intelligence, The AI Act, etc. 
    # maybe find a way to include both? would that be redundant? 
    {"label": "ARTINT", "pattern": "Artificial Intelligence"},
    {"label": "ARTINT", "pattern": "A.I."},
    {"label": "ARTINT", "pattern": "AI"},
    
# FIX! - not showing up at all
# CLAIM_VERB - any verb that indicates a claim
    # maybe categorize by strength? 
    # probably categorize strength using PhraseMatcher
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "argue"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "assert"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "claim"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "contend"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "propose"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "maintain"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "state"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "suggest"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "insist"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "hold"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "affirm"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "demonstrate"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "establish"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "establish"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "prove"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "show"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "validate"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "substantiate"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "confirm"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "illustrate"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "justify"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "interpret"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "evaluate"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "analyze"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "examine"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "assess"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "deduce"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "infer"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "theorize"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "posit"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "hypothesize"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "advocate"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "recommend"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "encourage"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "urge"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "emphasize"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "call for"}}, # might not work
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "champion"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "challenge"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "dispute"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "question"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "refute"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "counter"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "reject"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "critique"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "oppose"}},
    {"label": "CLAIM_VERB", "pattern": {"LEMMA": "problematize"}},


# FUNCTIONAL - identitifying some related technologies, probably removing other entities too
# TECH - any other technologies related to AI or the development of AI (adequate scope??)

    # will refine scope later 
    {"label": "TECH", "pattern": "Machine Learning"},
    {"label": "TECH", "pattern": "Deep Learning"},
    {"label": "TECH", "pattern": "Neural Networks"},
    {"label": "TECH", "pattern": "Transformers"},
    {"label": "TECH", "pattern": "Reinforcement Learning"},
    {"label": "TECH", "pattern": "Large Language Models"},
    {"label": "TECH", "pattern": "LLMs"},
    {"label": "TECH", "pattern": "Natural Language Processing"},
    {"label": "TECH", "pattern": "NLP"},
    {"label": "TECH", "pattern": "Computer Vision"},
    {"label": "TECH", "pattern": "Robotics"},
    {"label": "TECH", "pattern": "Autonomous Systems"}

# CLAIM_STRENGTH - strength of related claim verb


]
# possible claim structures: 
# [SOURCE] [CLAIM_STRENGTH] [CLAIM_VERB] [claim being made]
# [SOURCE] [CLAIM_VERB] [CLAIM_STRENGTH] [claim being made]

def main():


   

 


    # add Artifical Intelligence and AI (maybe others) as entities using Entity Ruler
    ruler = EntityRuler(nlp)

    # 

    # Will probably make helper function for adding custom entities 
    # and then build doc after

    #
    
       
    
    ruler = nlp.add_pipe("entity_ruler", after="ner")
    ruler.add_patterns(custom_ents)


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
    sentence_count = 0
    token_count = 0
    for doc in nlp.pipe(pull_s3_files()):
        sentence_count += len(list(doc.sents))
        token_count += len(doc)
        print("Current num of sentences: ", sentence_count)
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


    print("Sentence count: ", sentence_count)
    print("Token count: ", token_count)

if __name__ == "__main__": 
    main() 