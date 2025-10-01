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


  # For more sophisticated claim extraction, especially when dealing with varied and nuanced claims, 
  # you might train a custom machine learning model. This would involve:

        #     Annotating a dataset: Manually label text segments as claims or non-claims.
        #     Feature engineering: Extract features from spaCy's Doc objects (POS tags, dependency relations, entity types, etc.) to train your model.
        #     Training a classifier: Use a classification algorithm (e.g., Support Vector Machines, Logistic Regression, or deep learning models) to learn to identify claims.
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


text = """
OpenAI researchers argue that robust safety protocols are essential for advanced AI systems. The European Commission has proposed regulations to ensure ethical AI development. Google maintains that transparency in AI decision-making is critical for public trust. Some experts claim that autonomous weapons powered by AI pose significant risks to global security. The Partnership on AI suggests that collaboration between industry and academia can improve AI safety standards. Microsoft asserts that bias mitigation in AI models should be a top priority. The Future of Life Institute warns that unchecked AI development could lead to unintended consequences. IBM demonstrates that explainable AI can help users understand complex model outputs. Stanford University researchers contend that ethical guidelines must evolve alongside AI capabilities. The Alan Turing Institute recommends regular audits of AI systems to detect and correct harmful behaviors.

Meta affirms that user privacy must be protected in AI-driven platforms. The United Nations calls for international cooperation to address AI safety challenges. Some ethicists maintain that AI should never be used for mass surveillance. DeepMind illustrates that reinforcement learning agents can be aligned with human values through careful reward design. The Center for Humane Technology advocates for responsible AI deployment in social media. Tesla claims that self-driving cars require rigorous safety validation before widespread adoption. The AI Now Institute emphasizes the importance of public input in shaping AI policy. Researchers at MIT propose that AI ethics education should be integrated into computer science curricula. The World Economic Forum encourages governments to invest in AI safety research. Amazon suggests that fairness in AI-powered hiring tools is achievable with diverse training data.

Harvard scholars argue that AI systems must be held accountable for their decisions. The OECD recommends standardized reporting for AI incidents and failures. Some technologists contend that open-source AI frameworks can foster safer innovation. The National Institute of Standards and Technology (NIST) demonstrates that robust testing environments are vital for AI reliability. The Ethics Advisory Board at Google insists that stakeholder engagement is necessary for ethical AI governance. The Royal Society claims that interdisciplinary research can address complex AI safety issues. The Mozilla Foundation calls for greater transparency in AI algorithms used online. Some policy makers maintain that AI should be subject to strict liability laws. The Institute of Electrical and Electronics Engineers (IEEE) proposes global standards for AI safety and ethics. The Responsible AI Consortium advocates for continuous monitoring of deployed AI systems.

The AI Ethics Lab suggests that scenario analysis can help anticipate potential risks. Some legal experts argue that AI-generated content should be clearly labeled. The Global Partnership on Artificial Intelligence (GPAI) recommends sharing best practices for AI safety across borders. The Stanford Human-Centered AI initiative claims that participatory design can reduce ethical risks in AI products. The Oxford Internet Institute maintains that public awareness campaigns are needed to educate users about AI safety. The Center for Security and Emerging Technology (CSET) asserts that national security strategies must account for AI vulnerabilities. The Carnegie Mellon University Robotics Institute demonstrates that simulation-based testing can uncover hidden flaws in AI systems.
"""
# testing alternative approach ----------------------------------------------------------------------

# for regex testing
# def join_words_with_pipe(strings):
#     return "|".join(" ".join(s.split()) for s in strings)

# def print_orgs(): 
#     doc = nlp(pull_one_file())
#     orgs = [ent.text for ent in doc.ents if ent.label_=="ORG"]
#     print(*orgs, sep='\n')

# def print_claim_orgs(): 
#     # processed = pull_one_file().replace()
#     doc = nlp(pull_one_file())
#     pattern = join_words_with_pipe(claim_verb_terms)
#     claim_orgs = set()
#     claim_sents = set()
#     for ent in doc.ents: 
#         if ent.label_ != "ORG": 
#             continue
#         if re.search(pattern, ent.sent.text): 
#             claim_orgs.add(ent.text)
#             claim_sents.add(ent.sent.text)
#     print(*claim_orgs, sep='\n')
#     print(*claim_sents, sep='\n')
#     # print("\nCLAIM ORGS ----------------------------------")
#     # pprint.pprint(claim_orgs)
#     # print("\nCLAIM SENTENCES -----------------------------")
#     # pprint.pprint(claim_sents)

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
            # claim_phrase.append(a)
            # print(claim_phrase)
            claim_verb = a.lemma_

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
            # Maybe modify to get multiple claims within one sentence? 
            claim_subtree = " ".join([j.text for j in a.subtree if j.i > a.i and not j.is_punct])
            
            # then also add the direct object(s) of that claim verb, 
            # as long as the original token is in the same subtree as 
            # the direct object
            # (need to find a way to supercede non dobj dependencies and retrieve whole propositions?)
            # claim_phrase.extend([j for j in a.children if j.dep_ == "dobj" and source in a.subtree])
        
            break
    
    # expand out verb phrase to get modifiers of the direct object
    # for tok in claim_phrase: 
    #     for i in tok.children: 
    #         if i.dep_ == "amod": 
    #             claim_phrase.append(i)
             
    # new_list = sorted(claim_phrase, key=lambda x: x.i)
    return claim_verb, claim_subtree, strength_phrase



def test_preprocess(): 
    print(preprocess_text(pull_one_file()))
    # print(preprocess_text(text))

def test_one_file(file=""): 

    # if file not supplied
    if not file: file = preprocess_text(pull_one_file())
    doc = nlp(file)
    claims = []
    relations = dict()
    explicit_claims = []
    sources = ["ORG", "PERSON", "GPE"]
    for ent in doc.ents: 
        if ent.label_ in sources:
            claim_verb, claim_contents, strength_phrase = source_to_claim(ent.root)
            if claim_verb and claim_contents:
                explicit_claims.append((ent.text, claim_verb, claim_contents, strength_phrase))
 
    pprint.pprint(explicit_claims)

def see_relations(text):
    doc = nlp(text)
    displacy.serve(doc, style="dep", auto_select_port=True)

def test_all_files(): 
    for i, doc in enumerate(nlp.pipe(pull_all_files())):

        claims = []
        relations = dict()
        explicit_claims = []
        sources = ["ORG", "PERSON", "GPE"]
        strength_phrase_count = 0 # testing
        for ent in doc.ents: 
            if ent.label_ in sources: 
                claim_verb, claim_contents, strength_phrase = source_to_claim(ent.root)
                if strength_phrase: strength_phrase_count += 1 # testing 
                if claim_verb and claim_contents :
                    explicit_claims.append((ent.text, claim_verb, claim_contents, strength_phrase))
        print("Strength phrase count: ", strength_phrase_count)
        pprint.pprint(explicit_claims)

def main():
    # ------------------------------------------------------------------------------------------------
    # Where we're at: 
    # ------------------
    # - CURRENT STRUCTURE: [source] [claim_verb] [claim_contents] | [claim_strength]
    # - was able to use dependency parsing to retrieve source, claim verb (defined by our claim_verb_terms), claim contents, and 
    #   SOMETIMES claim strength. 
    # - the rules are very loose and this wont work 100% of the time
    # - for most sentences that make a single claim and their source is an ORG, PERSON, or GPE, 
    #   this algorithm will work at retrieving at least a triplet of data describing the claim
    # - if an adverb is used to describe the claim verb and it ends in ly, it will be added to claim strength.
    #   - obviously this isnt perfect as there are other words like "only" that also fit this description 
    #   - and im sure there are adverbs that dont end in ly that i should be looking for 
    # - if multiple claims are made in a sentence, then only the first claim verb is parsed as the claim verb, but the remainder
    #   of the claim contents should end up in the claim_contents 
    # - does not recognized pronouns, would have to use coreferencing in order to accomplish this 
    # - will probably want to give a score to the claim strength, need to figure out how to do this


    # see_relations("Miller however suggests that the user should be cautious with AI.")
    # see_relations("The UK relunctantly and with great hesitation suggests sanctions on AI.")
    # see_relations("The UK states that it is reluctant with AI and it also claims that it is destructive")
    test_one_file("The UK states that it is reluctant with AI and it also claims that it is destructive")
    test_one_file("The UK states that it is reluctant with AI. It also claims that it is destructive")
    see_relations("The UK states that it is reluctant with AI and it also claims that it is destructive")
if __name__ == "__main__": 
    main() 