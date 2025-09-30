import spacy 
from spacy.matcher import PhraseMatcher
from spacy.matcher import Matcher
from spacy.matcher import DependencyMatcher
from spacy.tokens import Span, SpanGroup
from spacy.pipeline import EntityRuler
from spacy.util import filter_spans
from utils.pull_text import pull_all_files
from utils.pull_text import pull_one_file
# load model
nlp = spacy.load("en_core_web_trf")

# --------------------------------------------------------------------------
# possible claim structures: 
# [SOURCE] [CLAIM_STRENGTH] [CLAIM_VERB] [claim being made]
# [SOURCE] [CLAIM_VERB] [CLAIM_STRENGTH] [claim being made]

# REVISED PLAN - 9/29
# instead of using entity ruler to override the entity rule set in ner, 
# add custom spans containing target terms and then use DependencyMatcher to 
# match syntactic structure of claims within our target domain 

# 1. add three span objects to doc.spans (for now): "artint", "tech", "claim_verbs"(?)
# 2. create pattern that matches our proposed claim structure 
#   - this is where most of the trial and error will be conducted 
#   - what are we looking for in a claim? 
#   - what should a claim look like? 
#   - should we have source (PERSON/ORG/pronoun), claim being made (object/complement)?
#   - are claim verbs even needed? 
#   - check if tech or artinit terms occur inside the object
# 
# 3. add the pattern(s) to DependencyMatcher 

# ----------------------------------------------------------------------------

# target terms for PhraseMatcher/Matcher
artint_terms = ["Artificial Intelligence", "AI", "A.I."]
claim_verb_terms = [
        "argue", "assert", "claim", "contend", "propose", "maintain", "state", "suggest", "insist", "hold",
        "affirm", "demonstrate", "establish", "prove", "show", "validate", "substantiate", "confirm", "illustrate",
        "justify", "interpret", "evaluate", "analyze", "examine", "assess", "deduce", "infer", "theorize", "posit",
        "hypothesize", "advocate", "recommend", "encourage", "urge", "emphasize", "call for", "champion", "challenge",
        "dispute", "question", "refute", "counter", "reject", "critique", "oppose", "problematize"
    ]
tech_terms = [
        "Machine Learning",
        "Deep Learning",
        "Neural Networks",
        "Transformers",
        "Reinforcement Learning",
        "Large Language Models",
        "LLMs",
        "Natural Language Processing",
        "NLP",
        "Computer Vision",
        "Robotics",
        "Autonomous Systems"
    ]

# add claim strength group???
# claim_strength = []

# create matchers for target terms
a_matcher = PhraseMatcher(nlp.vocab, attr="LOWER")
cv_matcher = Matcher(nlp.vocab)
t_matcher = PhraseMatcher(nlp.vocab, attr="LOWER")

# define patterns - must make doc for each individual term
artint_patterns = [nlp.make_doc(t) for t in artint_terms]
a_matcher.add("ARTINT", artint_patterns)
tech_patterns = [nlp.make_doc(t) for t in tech_terms]
t_matcher.add("TECH", tech_patterns)

# create list of lemma dictionaries for each claim verb
claim_verb_patterns = [[{"LEMMA": c}] for c in claim_verb_terms]
cv_matcher.add("CLAIM_VERB", claim_verb_patterns)


# tester for one file
def test_one_file(): 
    
    # process doc 
    doc = nlp(pull_one_file())

    try: 
        # get matches from doc 
        a_matches = a_matcher(doc)
        cv_matches = cv_matcher(doc)
        t_matches = t_matcher(doc)

        # create numeric ID for each span 
        # retrieve spans from doc 
        # add spans to corresponding group

        artint_id = nlp.vocab.strings["ARTINT"]
        artint_spans = [Span(doc, start, end, label=artint_id) for _, start, end in a_matches]
        doc.spans["art_int"] = artint_spans
        print("\nARTINT Spans --------------------------------------------")
        print(doc.spans["art_int"])

        claim_verb_id = nlp.vocab.strings["CLAIM_VERB"]
        claim_verb_spans = [Span(doc, start, end, label=claim_verb_id) for _, start, end in cv_matches]
        doc.spans["claim_verb"] = claim_verb_spans
        print("\nCLAIM_VERB Spans --------------------------------------------")
        print(doc.spans["claim_verb"])

        tech_id = nlp.vocab.strings["TECH"]
        tech_spans = [Span(doc, start, end, label=tech_id) for _, start, end in t_matches]
        doc.spans["tech"] = tech_spans
        print("\nTECH Spans --------------------------------------------")
        print(doc.spans["tech"])

        # PROCESS DOC......
        # -----------------------------------------
        # -----------------------------------------
        # -----------------------------------------
        # -----------------------------------------
        # -----------------------------------------
        # -----------------------------------------
        # -----------------------------------------
        # -----------------------------------------
        # for each sentence: 
        # 1. find claim verbs 
        # 2. look in dependency parse for the source of claim
        # 3. find object/complement 
        # 4. check if tech or ai terms occur inside the object 
        #       - this will qualify the claim
        tokens = []
       
        for v, s in enumerate(doc.sents):
            # note: didn't adjust index num for empty sentences, change later if needed.

            # if sentence is not empty
            if not str(s).isspace():
                    
                print("\nSentence number: ", v, "-----------------------")
                print("Original sent: ", str(s))
                print(str(s).isspace())
                # iterate over tokens 
                # check 
                #   - lemma or pos for a matching claim verb 
                #   - dep for lables like nsubj, dobj, (research these)
                #   - use subtree to get related tokens for source and object
                # tokens.append([t.text for t in s])
                claim_sent = []
                for t in s:
                    
                    # convert list of claim verb spans to strings for equality checking 
                    string_spans = [str(s) for s in list(doc.spans["claim_verb"])]

                    if t.dep_ == "nsubj" or t.dep_ == "dobj": 
                        print("Subject/Object found: ", t.text)
                        claim_sent.append(t.text)
                    # check for tokens that are claim verbs
                    if t.pos_ == "VERB" and t.lemma_ in string_spans:
                        print("CLAIM VERB FOUND: ", t.text)
                        print("Syntactic dependency relation: ", t.dep_)
                        claim_sent.append(t.text)
                    
                if claim_sent: print("Claim sent: ", claim_sent)
            
            # breaks at 1000 sentences
            if v == 1000: 
                break
                
    except Exception as e: 
        print(e)
        return 
    
# tester for all files
def test_all_files(): 
    # # create matchers for target terms
    # a_matcher = PhraseMatcher(nlp.vocab, attr="LOWER")
    # cv_matcher = Matcher(nlp.vocab)
    # t_matcher = PhraseMatcher(nlp.vocab, attr="LOWER")
    
    # # define patterns - must make doc for each individual term
    # artint_patterns = [nlp.make_doc(t) for t in artint_terms]
    # a_matcher.add("ARTINT", artint_patterns)
    # tech_patterns = [nlp.make_doc(t) for t in tech_terms]
    # t_matcher.add("TECH", tech_patterns)

    # # create list of lemma dictionaries for each claim verb
    # claim_verb_patterns = [[{"LEMMA": c}] for c in claim_verb_terms]
    # cv_matcher.add("CLAIM_VERB", claim_verb_patterns)

    # possibly change and process all text at once
    for i, doc in enumerate(nlp.pipe(pull_all_files())):
        # ****** must make doc per span to add ******
        print("Num of sents: ", len(list(doc.sents)))
        print("File number: ", i)
        try: 
            # get matches from doc 
            a_matches = a_matcher(doc)
            cv_matches = cv_matcher(doc)
            t_matches = t_matcher(doc)

            # create numeric ID for each span 
            # retrieve spans from doc 
            # add spans to corresponding group

            artint_id = nlp.vocab.strings["ARTINT"]
            artint_spans = [Span(doc, start, end, label=artint_id) for _, start, end in a_matches]
            doc.spans["art_int"] = artint_spans
            print("\nARTINT Spans --------------------------------------------")
            print(doc.spans["art_int"])

            claim_verb_id = nlp.vocab.strings["CLAIM_VERB"]
            claim_verb_spans = [Span(doc, start, end, label=claim_verb_id) for _, start, end in cv_matches]
            doc.spans["claim_verb"] = claim_verb_spans
            print("\nCLAIM_VERB Spans --------------------------------------------")
            print(doc.spans["claim_verb"])

            tech_id = nlp.vocab.strings["TECH"]
            tech_spans = [Span(doc, start, end, label=tech_id) for _, start, end in t_matches]
            doc.spans["tech"] = tech_spans
            print("\nTECH Spans --------------------------------------------")
            print(doc.spans["tech"])

            # PROCESS DOC......
            # -----------------------------------------
            # -----------------------------------------
            # -----------------------------------------
            # -----------------------------------------
            # -----------------------------------------
            # -----------------------------------------
            # -----------------------------------------
            # -----------------------------------------
            # for each sentence: 
            # 1. find claim verbs 
            # 2. look in dependency parse for the source of claim
            # 3. find object/complement 
            # 4. check if tech or ai terms occur inside the object 
            #       - this will qualify the claim
            tokens = []

            # testing 
            print("Span list type: ", type(doc.spans["claim_verb"]))
            print("Span list type: ", type(list(doc.spans["claim_verb"])))
            print("Span list item type: ", type(str(list(doc.spans["claim_verb"]).pop())))
            print(list(doc.spans["claim_verb"]))
            for v, s in enumerate(doc.sents):
                print("Sentence number: ", v)
                print("Original sent: ", s)
                # iterate over tokens 
                # check 
                #   - lemma or pos for a matching claim verb 
                #   - dep for lables like nsubj, dobj, (research these)
                #   - use subtree to get related tokens for source and object
                # tokens.append([t.text for t in s])
                claim_sent = []
                for t in s:
                   
                    # convert list of claim verb spans to strings for equality checking 
                    string_spans = [str(s) for s in list(doc.spans["claim_verb"])]

                    if t.dep_ == "nsubj" or t.dep_ == "dobj": 
                        print("Subject/Object found: ", t.text)
                        claim_sent.append(t.text)
                    # check for tokens that are claim verbs
                    if t.pos_ == "VERB" and t.lemma_ in string_spans:
                        print("CLAIM VERB FOUND: ", t.text)
                        print("Syntactic dependency relation: ", t.dep_)
                        claim_sent.append(t.text)
                    
                
                print("\nClaim sent: ", claim_sent)
                
                    # if token is a claim verb 
                    # if t.lemma_ in doc.spans["claim_verb"] or t.pos_ in doc.spans["claim_verb"]: 

                # breaks at 1000 sentences
                
                if v == 1000: 
                    break
            
            
            if i == 0: 
                break
            
        except Exception as e: 
            print(e)
            return 
        
def main():
    test_one_file()

if __name__ == "__main__": 
    main() 