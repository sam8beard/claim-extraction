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
        "dispute", "question", "refute", "counter", "reject", "critique", "oppose", "problematize", "believe"
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

longtext = (
    "The report by the UN argues that climate change is accelerating.\n"
    "According to John Smith, AI will revolutionize energy production by 2030.\n"
    "Google clearly believes its new algorithm improves search accuracy.\n"
    "The French government states that renewable energy will be the primary power source.\n"
    "Researchers claim that AI-generated art could replace traditional methods entirely.\n"
    "The World Health Organization warns that antibiotic resistance is a growing threat.\n"
    "Elena Torres argues that economic reform will reduce inequality.\n"
    "Experts suggest that this policy will likely improve public safety.\n"
    "Tesla claims its new battery design is more efficient than competitors.\n"
    "The US Department of Energy insists that fusion power will soon be viable.\n"
    "According to the BBC, global food shortages will increase in the next decade.\n"
    "Microsoft believes its latest cloud technology enhances data security.\n"
    "The European Union claims that new regulations will reduce carbon emissions.\n"
    "Dr. Amanda Lee argues that mental health awareness will improve community wellbeing.\n"
    "The World Bank insists that infrastructure investment will boost economic growth.\n"
    "Greenpeace warns that deforestation rates are dangerously high.\n"
    "Harvard researchers claim that sleep quality affects cognitive performance significantly.\n"
    "The Japanese government states that electric vehicles will dominate the market by 2040.\n"
    "According to Dr. Robert Chen, AI ethics will become a central concern for policymakers.\n"
    "Facebook argues that its new privacy features improve user trust.\n"
    "NASA reports that asteroid mining could become feasible within 50 years.\n"
    "According to Reuters, interest rates will likely rise next year.\n"
    "The CDC insists that vaccination campaigns reduce disease outbreaks.\n"
    "According to CNN, climate migration will increase in the coming years.\n"
    "The IMF claims that debt restructuring will stabilize developing economies.\n"
    "Dr. Sarah Patel argues that renewable energy can meet global demand by 2050.\n"
    "Apple claims its new chip design improves efficiency dramatically.\n"
    "The WHO warns that antibiotic misuse is accelerating resistance.\n"
    "According to Bloomberg, quantum computing will transform finance.\n"
    "The German government states that hydrogen energy will play a key role in decarbonization.\n"
    "According to the New York Times, AI bias remains a major challenge.\n"
    "Amazon believes its logistics improvements will cut delivery times.\n"
    "The UNDP insists that education access reduces poverty.\n"
    "Dr. Miguel Alvarez claims that urban green spaces improve mental health.\n"
    "According to Al Jazeera, global water shortages will worsen.\n"
    "The World Economic Forum states that automation will reshape labor markets.\n"
    "According to Wired, cybersecurity threats will evolve rapidly.\n"
    "The Chinese government claims that its space program will achieve a moon base.\n"
    "Dr. Emily Chen argues that AI-driven healthcare can improve diagnosis accuracy.\n"
    "The International Energy Agency warns that fossil fuel dependence persists.\n"
    "According to Fox News, renewable energy adoption will accelerate.\n"
    "The UK government insists that AI regulation will improve innovation.\n"
    "According to NPR, social media influences public opinion strongly.\n"
    "The Canadian government claims that clean technology will create jobs.\n"
    "Dr. Mark Reynolds argues that AI literacy is essential for future workforces.\n"
    "According to The Guardian, climate resilience is key to sustainable development.\n"
    "The World Food Programme warns that hunger crises are worsening globally.\n"
    "According to Financial Times, global trade patterns are shifting.\n"
    "The Brazilian government insists that conservation efforts will protect biodiversity.\n"
    "According to Vox, misinformation spreads faster than truth.\n"
    "The African Union claims that regional cooperation will strengthen economies.\n"
)
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
        # print("Ancestor: ", a.text)
        if a.lemma_ in claim_verb_terms and a.pos_ == "VERB": 
    
            claim_verb = a.lemma_

            # checks for adverb modifier that indicates degree of claim (not perfect)
            advmod = [j for j in a.children if j.dep_ == "advmod" and "ly" in j.suffix_]
            
            prep = [j for j in a.children if j.dep_ == "prep"]

            if advmod: 
                advmod = advmod.pop()
                # print(advmod.suffix_) # testing
                # construct strength
                advmod_string = advmod.text

                # strength_phrase = advmod_string + " " + " ".join([j.text for j in advmod.children])
                strength_phrase = advmod_string
                # print(strength_phrase, advmod.sentiment) # testing
        
            # should we even consider this case???? (prep [...] noun)
            # if prep: 
            #     prep = prep.pop()
            #     strength_phrase = " ".join([j.text for j in prep.subtree if j.i > prep.i and j.pos_ != "PROPN"])
            #     print(strength_phrase)

            # all descendants to the right of the root/claim verb that 
            # Maybe modify to get multiple claims within one sentence? 
            claim_subtree = " ".join([j.text for j in a.subtree if j.i > a.i and not j.is_punct])
            
            break
            # then also add the direct object(s) of that claim verb, 
            # as long as the original token is in the same subtree as 
            # the direct object

        # if verb is root of sentence but not the claim verb
        # there will not be a claim verb in this instance
        elif a.sent.root == a and a.pos_ == "VERB" and a.lemma_ not in claim_verb_terms: 
            # nsubj = [j for j in a.children if j.dep_ == "nsubj"].pop()
            # claim_subtree = nsubj.text + " " + " ".join([j.text for j in a.subtree if j.i > a.i and not j.is_punct])
            # claim_subtree = " ".join([j.text for j in a.sent[source.i + 1:]])
            claim_subtree = " ".join([j.text for j in a.sent if j.i > source.i and not j.is_punct])
            # print("Tracing: ", [j.text for j in a.sent])
            # print("All tokens after source: ", [j.text for j in a.sent if j.i > source.i and not j.is_punct])
            # print("Source: ", source.text) # testing
            # print("Root verb identified: ", a.text) # testing
            # print("Testing: ", claim_subtree) # testing

            break
    
             
    return claim_verb, claim_subtree, strength_phrase



def test_preprocess(): 
    print(preprocess_text(pull_one_file()))
    # print(preprocess_text(text))

def test_one_file(file=""): 

    # if file not supplied
    if not file: file = preprocess_text(pull_one_file())
    # print(file)
    doc = nlp(file)
    claims = []
    relations = dict()
    explicit_claims = []
    sources = ["ORG", "PERSON", "GPE", "NORP", "LAW"]
    strength_phrase_count = 0 # testing
    processed_sents_count = 0 # testing
    processed_sents = [] # testing
    for ent in doc.ents: 
        # print(ent.text, ent.label_) # testing
        # print(ent.root.text) # testing
        if ent.label_ in sources:
            claim_verb, claim_contents, strength_phrase = source_to_claim(ent.root)
            if strength_phrase: strength_phrase_count += 1 # testing
            if claim_contents:
                processed_sents_count += 1 # testing
                processed_sents.append(ent.root.sent.text) # testing
                # print(f"\nProcessed sentence {processed_sents_count}: {ent.root.sent.text}") # testing 
                # print(ent.text)
                explicit_claims.append((ent.text, claim_verb, claim_contents, strength_phrase))
 
    pprint.pprint(explicit_claims)
    # pprint.pprint(processed_sents)
    print(f"Number of claims identified: {len(explicit_claims)}" )

def test_all_files(): 
    for i, doc in enumerate(nlp.pipe(pull_all_files())):

        claims = []
        relations = dict()
        explicit_claims = []
        sources = ["ORG", "PERSON", "GPE", "NORP", "LAW"]
        strength_phrase_count = 0 # testing
        processed_sents_count = 0 # testing
        processed_sents = [] # testing
        for ent in doc.ents: 
            if ent.label_ in sources: 
                claim_verb, claim_contents, strength_phrase = source_to_claim(ent.root)
                if strength_phrase: strength_phrase_count += 1 # testing 
                if claim_contents :
                    processed_sents_count += 1 # testing
                    processed_sents.append(ent.root.sent.text)
                    explicit_claims.append((ent.text, claim_verb, claim_contents, strength_phrase))
        print("Strength phrase count: ", strength_phrase_count) # testing
        pprint.pprint(explicit_claims)
        pprint.pprint(processed_sents)
        
def see_relations(text):
    doc = nlp(text)
    displacy.serve(doc, style="dep", auto_select_port=True)

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
    # test_one_file("The UK states that it is reluctant with AI and it also claims that it is destructive")
    # test_one_file("The UK states that it is reluctant with AI. It also claims that it is destructive")
    # see_relations("The UK states that it is reluctant with AI and it also claims that it is destructive")
    test_one_file(longtext)
    
    # see_relations("According to John Smith, AI will revolutionize energy production by 2030.")
    # see_relations("According to John Smith, AI will revolutionize energy production by 2030. According to the BBC, global food shortages will increase in the next decade. According to Dr. Robert Chen, AI ethics will become a central concern for policymakers.  According to Reuters, interest rates will likely rise next year. According to CNN, climate migration will increase in the coming years. According to Bloomberg, quantum computing will transform finance. According to the New York Times, AI bias remains a major challenge.  According to Wired, cybersecurity threats will evolve rapidly. According to Fox News, renewable energy adoption will accelerate. According to NPR, social media influences public opinion strongly. According to The Guardian, climate resilience is key to sustainable development. According to Financial Times, global trade patterns are shifting. According to Vox, misinformation spreads faster than truth.")
    # test_one_file("In an article from CNN, it was stated that AI is great. CNN reported that AI is not good. According to John Smith, AI will revolutionize energy production by 2030. According to the BBC, global food shortages will increase in the next decade. According to Dr. Robert Chen, AI ethics will become a central concern for policymakers.  According to Reuters, interest rates will likely rise next year. According to CNN, climate migration will increase in the coming years. According to Bloomberg, quantum computing will transform finance. According to the New York Times, AI bias remains a major challenge.  According to Wired, cybersecurity threats will evolve rapidly. According to Fox News, renewable energy adoption will accelerate. According to NPR, social media influences public opinion strongly. According to The Guardian, climate resilience is key to sustainable development. According to Financial Times, global trade patterns are shifting. According to Vox, misinformation spreads faster than truth.")
    # see_relations("The French government states that renewable energy will be the primary power source. The Japanese government states that electric vehicles will dominate the market by 2040.")
    # test_one_file("The French government states that renewable energy will be the primary power source. The Japanese government states that electric vehicles will dominate the market by 2040.")
if __name__ == "__main__": 
    main() 