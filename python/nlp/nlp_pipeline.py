import spacy 
from spacy.matcher import PhraseMatcher
from spacy.matcher import Matcher
from spacy.matcher import DependencyMatcher
from spacy.tokens import Span, SpanGroup
from spacy.pipeline import EntityRuler
from spacy.util import filter_spans
from utils.pull_text import pull_s3_files

# load model
nlp = spacy.load("en_core_web_trf")

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


def main():

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
    claim_strength = []



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

    # possibly change and process all text at once
    for doc in nlp.pipe(pull_s3_files()):
        # ****** must make doc per span to add ******
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
        except Exception as e: 
            print(e)
            return 
    

if __name__ == "__main__": 
 
    main() 