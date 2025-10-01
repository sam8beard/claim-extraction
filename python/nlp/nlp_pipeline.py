import spacy 
import re
import unicodedata
from spacy import displacy
import pprint
from spacy.matcher import PhraseMatcher
from spacy.matcher import Matcher
from spacy.matcher import DependencyMatcher
from spacy.tokens import Span, SpanGroup
from spacy.pipeline import EntityRuler
from spacy.util import filter_spans
from utils.pull_text import pull_all_files
from utils.pull_text import pull_one_file
from utils.pull_text import preprocess_text

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

text = """
OpenAI researchers argue that robust safety protocols are essential for advanced AI systems. The European Commission has proposed regulations to ensure ethical AI development. Google maintains that transparency in AI decision-making is critical for public trust. Some experts claim that autonomous weapons powered by AI pose significant risks to global security. The Partnership on AI suggests that collaboration between industry and academia can improve AI safety standards. Microsoft asserts that bias mitigation in AI models should be a top priority. The Future of Life Institute warns that unchecked AI development could lead to unintended consequences. IBM demonstrates that explainable AI can help users understand complex model outputs. Stanford University researchers contend that ethical guidelines must evolve alongside AI capabilities. The Alan Turing Institute recommends regular audits of AI systems to detect and correct harmful behaviors.

Meta affirms that user privacy must be protected in AI-driven platforms. The United Nations calls for international cooperation to address AI safety challenges. Some ethicists maintain that AI should never be used for mass surveillance. DeepMind illustrates that reinforcement learning agents can be aligned with human values through careful reward design. The Center for Humane Technology advocates for responsible AI deployment in social media. Tesla claims that self-driving cars require rigorous safety validation before widespread adoption. The AI Now Institute emphasizes the importance of public input in shaping AI policy. Researchers at MIT propose that AI ethics education should be integrated into computer science curricula. The World Economic Forum encourages governments to invest in AI safety research. Amazon suggests that fairness in AI-powered hiring tools is achievable with diverse training data.

Harvard scholars argue that AI systems must be held accountable for their decisions. The OECD recommends standardized reporting for AI incidents and failures. Some technologists contend that open-source AI frameworks can foster safer innovation. The National Institute of Standards and Technology (NIST) demonstrates that robust testing environments are vital for AI reliability. The Ethics Advisory Board at Google insists that stakeholder engagement is necessary for ethical AI governance. The Royal Society claims that interdisciplinary research can address complex AI safety issues. The Mozilla Foundation calls for greater transparency in AI algorithms used online. Some policy makers maintain that AI should be subject to strict liability laws. The Institute of Electrical and Electronics Engineers (IEEE) proposes global standards for AI safety and ethics. The Responsible AI Consortium advocates for continuous monitoring of deployed AI systems.

The AI Ethics Lab suggests that scenario analysis can help anticipate potential risks. Some legal experts argue that AI-generated content should be clearly labeled. The Global Partnership on Artificial Intelligence (GPAI) recommends sharing best practices for AI safety across borders. The Stanford Human-Centered AI initiative claims that participatory design can reduce ethical risks in AI products. The Oxford Internet Institute maintains that public awareness campaigns are needed to educate users about AI safety. The Center for Security and Emerging Technology (CSET) asserts that national security strategies must account for AI vulnerabilities. The Carnegie Mellon University Robotics Institute demonstrates that simulation-based testing can uncover hidden flaws in AI systems.
"""
# testing alternative approach ----------------------------------------------------------------------

# for regex testing
def join_words_with_pipe(strings):
    return "|".join(" ".join(s.split()) for s in strings)

def print_orgs(): 
    doc = nlp(pull_one_file())
    orgs = [ent.text for ent in doc.ents if ent.label_=="ORG"]
    print(*orgs, sep='\n')

def print_claim_orgs(): 
    # processed = pull_one_file().replace()
    doc = nlp(pull_one_file())
    pattern = join_words_with_pipe(claim_verb_terms)
    claim_orgs = set()
    claim_sents = set()
    for ent in doc.ents: 
        if ent.label_ != "ORG": 
            continue
        if re.search(pattern, ent.sent.text): 
            claim_orgs.add(ent.text)
            claim_sents.add(ent.sent.text)
    print(*claim_orgs, sep='\n')
    print(*claim_sents, sep='\n')
    # print("\nCLAIM ORGS ----------------------------------")
    # pprint.pprint(claim_orgs)
    # print("\nCLAIM SENTENCES -----------------------------")
    # pprint.pprint(claim_sents)

# a function that identifies claim verbs and direct objects that are grammatically linked to a source
def source_to_claim(source):
    claim_phrase = []
    claim_verb = ""
    claim_subtree = ""
    strength_phrase = ""
    # iterate through all ancestors of the source token
    for a in source.ancestors: 
        # when you get to a claim verb
        if a.lemma_ in claim_verb_terms and a.pos_ == "VERB": 
            # add verb to phrase list
            # print("\nAdding verb: ", a.text)
            claim_phrase.append(a)
            # print(claim_phrase)
            claim_verb = a.lemma_

            # check for adverb modifier on claim verb (should only be one)
            # advmod = [j for j in a.children if j.dep_ == "advmod" and j.suffix_ == "-ly"]

            # checks for adverb modifier that indicates degree of claim (not perfect)
            advmod = [j for j in a.children if j.dep_ == "advmod" and "ly" in j.suffix_]
            
            prep = [j for j in a.children if j.dep_ == "prep"]

            if advmod: 
                advmod = advmod.pop()
                print(advmod.suffix_)
                # construct strength
                advmod_string = advmod.text

                # strength_phrase = advmod_string + " " + " ".join([j.text for j in advmod.children])
                strength_phrase = advmod_string
                print(strength_phrase, advmod.sentiment)

            # should we even consider this case???? (prep [...] noun)
            # if prep: 
            #     prep = prep.pop()
            #     strength_phrase = " ".join([j.text for j in prep.subtree if j.i > prep.i and j.pos_ != "PROPN"])
            #     print(strength_phrase)

            # all descendants to the right of the root/claim verb that 
            claim_subtree = " ".join([j.text for j in a.subtree if j.i > a.i and not j.is_punct])
            
            # then also add the direct object(s) of that claim verb, 
            # as long as the original token is in the same subtree as 
            # the direct object
            # (need to find a way to supercede non dobj dependencies and retrieve whole propositions?)
            # print("Adding direct objects...")
            claim_phrase.extend([j for j in a.children if j.dep_ == "dobj" and source in a.subtree])
            # print(claim_phrase)
            # print("Verb subtree: ", list(a.subtree)) # gets every single word in a claim sentence

            # stop after the first verb
            # print("Verb children:" , list(a.children))
            # print("Right children of verb: ", [t.text for t in ])
        
            break
    
    # expand out verb phrase to get modifiers of the direct object
    for tok in claim_phrase: 
        for i in tok.children: 
            # print("Child of ", tok.text, ": ", i.text)
            if i.dep_ == "amod": 
                # print("Adding modifers to verb phrase...")
                claim_phrase.append(i)
                # print(claim_phrase)
                
            # if i.dep_ == "prep":
            #     # claim_phrase.extend([j for j in i.children if j.dep_ == "pobj" and source in  ])
            #     claim_phrase.extend([c for c in i.children if source in tok.subtree])
            #     print("Whats this look like: ", claim_phrase)
        # for i in tok.children: 
        #     if i.dep_ == "prep": 
        #         claim_phrase.extend([j for j in ])
    # sort tokens by position in original sentence
    # print(claim_phrase)
    new_list = sorted(claim_phrase, key=lambda x: x.i)
    return ''.join([i.text_with_ws for i in new_list]).strip(), claim_verb, claim_subtree, strength_phrase



def test_preprocess(): 
    print(preprocess_text(pull_one_file()))
    # print(preprocess_text(text))

def print_source_to_claim(): 
    # doc = nlp(pull_one_file())
    doc = nlp(preprocess_text(pull_one_file()))
    claims = []
    relations = dict()
    explicit_claims = []
    sources = ["ORG", "PERSON", "GPE"]
    for ent in doc.ents: 
        if ent.label_ in sources:
            relations[ent.text], claim_verb, claim_subtree, strength_phrase = source_to_claim(ent.root)
            if claim_verb and claim_subtree:
                explicit_claims.append((ent.text, claim_verb, claim_subtree, strength_phrase))
            # claims.append(source_to_claim(ent.root))
    # for token in doc: 
    #     if token.ent_type_ == "ORG": 
    #         relations[token.text] = source_to_claim(token)
    #         claims.append(source_to_claim(token))

    # for ent in doc.ents: 
    #     if ent
    # print(sorted(list(set(claims))))
    # see_relations(doc)
    # pprint.pprint(relations)
    pprint.pprint(explicit_claims)
    # print(relations)

def see_relations(text):
    doc = nlp(text)
    displacy.serve(doc, style="dep", auto_select_port=True)

def lets_see(): 
    for i, doc in enumerate(nlp.pipe(pull_all_files())):

    # doc = nlp("The UK with great strength suggests sanctions on AI."
    # doc = nlp("The UK relunctantly and with great hesitation suggest sanctions on AI.")
    # doc = nlp(text)
    # doc = nlp(pull_one_file())
        claims = []
        relations = dict()
        explicit_claims = []
        sources = ["ORG", "PERSON", "GPE"]
        strength_phrase_count = 0 # testing
        for ent in doc.ents: 
            if ent.label_ in sources: 
                relations[ent.text], claim_verb, claim_subtree, strength_phrase = source_to_claim(ent.root)
                if strength_phrase: strength_phrase_count += 1 # testing 
                if claim_verb and claim_subtree :
                    explicit_claims.append((ent.text, claim_verb, claim_subtree, strength_phrase))
        print("Strength phrase count: ", strength_phrase_count)
        pprint.pprint(explicit_claims)

# tester for one file ------------------------------------------------------------------------------------
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
        # print("\nARTINT Spans --------------------------------------------")
        # print(doc.spans["art_int"])

        claim_verb_id = nlp.vocab.strings["CLAIM_VERB"]
        claim_verb_spans = [Span(doc, start, end, label=claim_verb_id) for _, start, end in cv_matches]
        doc.spans["claim_verb"] = claim_verb_spans
        # print("\nCLAIM_VERB Spans --------------------------------------------")
        # print(doc.spans["claim_verb"])

        tech_id = nlp.vocab.strings["TECH"]
        tech_spans = [Span(doc, start, end, label=tech_id) for _, start, end in t_matches]
        doc.spans["tech"] = tech_spans
        # print("\nTECH Spans --------------------------------------------")
        # print(doc.spans["tech"])


       
        
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
        # ------------------------------------------------------------
        # iterate over tokens 
        # check 
        #   - lemma or pos for a matching claim verb 
        #   - dep for lables like nsubj, dobj, (research these)
        #   - use subtree to get related tokens for source and object
        #--------------------------------------------------------------
        # EXPLICIT CLAIM STRUCTURE: 
        # Dependencies USUALLY are 
        #   - root verb (claim_verb)
        #   - nsubj (source)
        #   - ccomp (clausal complement -> proposition/claim content)

        # structure triplet might look like: 
        # (source: nsubj, verb: claim_verb, proposition: ccomp)

        # REVISED APPROACH 2.0
        # from spacy.matcher import Matcher

        # matcher = Matcher(nlp.vocab)
        # # Example: Match "X announced Y" where X is a PROPN (Proper Noun) and Y is a NOUN or VERB
        # pattern = [{"POS": "PROPN"}, {"LEMMA": "announce"}, {"POS": {"IN": ["NOUN", "VERB"]}}]
        # matcher.add("CLAIM_ANNOUNCEMENT", [pattern])

        # matches = matcher(doc)
        # for match_id, start, end in matches:
        #     span = doc[start:end]
        #     print(f"Claim candidate: {span.text}")



        # Dependency Parsing for Deeper Understanding.
        # Analyze the dependency tree to understand the relationships between words and 
        # identify the subject and object of assertive verbs, which can help in pinpointing the core of a claim.
        # from spacy.matcher import Matcher

        # matcher = Matcher(nlp.vocab)
        # # Example: Match "X announced Y" where X is a PROPN (Proper Noun) and Y is a NOUN or VERB
        # pattern = [{"POS": "PROPN"}, {"LEMMA": "announce"}, {"POS": {"IN": ["NOUN", "VERB"]}}]
        # matcher.add("CLAIM_ANNOUNCEMENT", [pattern])

        # matches = matcher(doc)
        # for match_id, start, end in matches:
        #     span = doc[start:end]
        #     print(f"Claim candidate: {span.text}")
      
        # For more sophisticated claim extraction, especially when dealing with varied and nuanced claims, you might train a custom machine learning model. This would involve:

        #     Annotating a dataset: Manually label text segments as claims or non-claims.
        #     Feature engineering: Extract features from spaCy's Doc objects (POS tags, dependency relations, entity types, etc.) to train your model.
        #     Training a classifier: Use a classification algorithm (e.g., Support Vector Machines, Logistic Regression, or deep learning models) to learn to identify claims.

        for v, s in enumerate(doc.sents):
            # note: didn't adjust index num for empty sentences, change later if needed.

            # if sentence is not empty
            if not str(s).isspace():
                    
                print("\nSentence number: ", v, "-----------------------")
                print("Original sent: ", str(s))
                
                # convert list of claim verb spans to strings for equality checking 
                string_spans = [str(s) for s in list(doc.spans["claim_verb"])]

                for t in s:
                    # structured claim
                    source, claim_verb, prop = "", "", ""
                    claim_triplet = []

                    # if token is claim verb
                    if t.text in string_spans:
                        
                        print("\nClaim verb found: ", t.text)
                        children = [tok for tok in t.subtree]
                        print("\nChildren of claim verb: ", children)
                        claim_verb = t.text
                        
                        # iterate through subtree
                        for child in children: 
                            # find subject
                            if child.dep_ == "nsubj": 
                                # print("\nSubject found: ", child.text)
                                source = child.text

                            # find proposition
                            if child.dep_ == "ccomp" or child.dep_ == "dobj": 
                                # print("\nProposition found: ", child.text, child.dep_)
                                prop = child.text

                        claim_triplet.extend([source, claim_verb, prop])
                        # print("\nClaim triplet: ", claim_triplet)

            # breaks at 1000 sentences
            if v == 1000: 
                break
                
    except Exception as e: 
        print(e)
        return 
    
# tester for all files ------------------------------------------------------------------------------------
def test_all_files(): 
    
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
    # test_one_file()
    # test_approach_2()
    # print_orgs()
    # print_claim_orgs()
    # print_source_to_claim()
    # test_preprocess()
    # lets_see()
    # see_relations("The UK strongly and willingly suggests sanctions on AI.")
    lets_see()
    # see_relations("Miller however suggests that the user should be cautious with AI.")
    # see_relations("The UK relunctantly and with great hesitation suggests sanctions on AI.")

if __name__ == "__main__": 
    main() 