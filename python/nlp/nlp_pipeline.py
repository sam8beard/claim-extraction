from tqdm import tqdm
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
from utils.pull_text import pull_n_files
from utils.pull_text import preprocess_text
from spacy.language import Language, Doc
import logging 
import sys 


logging.basicConfig(level=logging.INFO, stream=sys.stderr)

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
        "dispute", "question", "refute", "counter", "reject", "critique", "oppose", "problematize", "believe", "warn"
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

claim_mod_terms = [
    "absolutely", "allegedly", "apparently", "arguably", "assuredly", "boldly", "clearly", "confidently", "conclusively",
    "definitely", "deeply", "dramatically", "evidently", "exactly", "explicitly", "expressly", "firmly", "forcefully",
    "frequently", "generally", "highly", "honestly", "importantly", "indeed", "indisputably", "ironically", "literally",
    "mostly", "naturally", "notably", "officially", "openly", "obviously", "often", "particularly", "persistently",
    "plainly", "positively", "potentially", "powerfully", "precisely", "probably", "profoundly", "purportedly", "quietly",
    "rarely", "repeatedly", "reportedly", "respectfully", "rigorously", "roughly", "seriously", "significantly",
    "solemnly", "specifically", "strongly", "supposedly", "surely", "tentatively", "thoroughly", "truly", "undoubtedly",
    "unquestionably", "usually", "verbally", "visibly", "widely", "willingly"
]


example_text = """
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

target_sources = ["ORG", "PERSON", "GPE", "NORP", "LAW"]
# testing alternative approach ----------------------------------------------------------------------

# initialize test doc
def initialize_doc(file="") -> Doc: 
    if not file: file = pull_one_file()
    doc = nlp(file)

    return doc

def initialize_all_docs(): 
    for doc in nlp(pull_all_files()): 
        yield doc

# def trim_whitespace(span): 
#     text = span.text
#     stripped_text = text.strip()
#     offset = len(text) - len(text.lstrip())
#     right_offset = len(text) - len(text.rstrip())
#     new_start_char = span.start_char + offset
#     new_end_char = new_start_char + len(stripped_text) # testing
#     return span.doc.char_span(new_start_char, span.end_char)


# uses modified get_tuples_spcat_2 to retrieve training data with char offsets
def get_training_list_spcat_2(): 
    nlp = spacy.load("en_core_web_trf")
    train_data = [] # add tuples here
    final_text = "" # add entire claim sentence here
    annotations = [] # add offset tuples here

    

    claim_count = 0
    sent_count = 0
    seen_sents = set()
    num = 10
    # NOTE: testing with pull_n_files
    for doc in nlp.pipe(tqdm(pull_all_files(), total=150, bar_format='{l_bar}{bar:20}{r_bar}', desc="Preparing training data..."), batch_size=8):
    # for doc in nlp.pipe(tqdm(pull_n_files(num), total=num, desc="Preparing training data..."), batch_size=8):

        for ent in doc.ents:
            
            if ent.label_ in target_sources:
                source_span = ent.root.sent

                sent_text = source_span.text.strip()

                if sent_text not in seen_sents: 
                    seen_sents.add(sent_text)

                    
                    
                    for entry in get_tuples_spcat_2(source_span): 
                        
                        if entry: 
                            # final_text = preprocess_text(source_span.text)
                            # NOTE: trying without preprocessing
                            final_text = sent_text
                            claim_count += 1
                            annotations = entry
                            data = ((final_text, annotations))
                            train_data.append(data)
    # for data in train_data: 
    #     print(f"\n\n{data}\n\n")
    return train_data

# gets spans based on char offsets 
def get_tuples_spcat_2(sent):
    source_start, source_end = 0, 0
    verb_start, verb_end = 0, 0
    content_start, content_end = 0, 0
    claim_mod_start, claim_mod_end = 0, 0
    claim_mod = ""
    

    for source in sent.ents: 

        # if source is in target entities, is capitalized, and is at least 2 chars long -> valid
        if source.label_ in target_sources and source.text[0].isupper() and len(source.text) >= 2: 
            
        

            # new start of entity offset 
            source_start = source.start_char - sent.start_char
            source_end = source.end_char - sent.start_char
            # source_start = ent_start
            # source_end = ent_end
          
            # NOTE: we have found a valid claim
            for a in source.root.ancestors: 
                if a.lemma_ in claim_verb_terms and a.pos_ == "VERB":
                    
                    verb = a
                  
                    verb_start, verb_end = get_new_token_offset(verb, sent)

                    # checks for adverb modifier that indicates degree of claim (not perfect but hopefully catches some)
                    # if the modifier is in the children of the claim verb, is an adverb, ends in -ly, and is directly before the verb


                    # check children of verb for claim mod candidates
                    advmod_candidates = [j for j in a.children if j.dep_ == "advmod" and j.text.lower() in claim_mod_terms]

                    # if nothing, look 2 tokens on either side of the verb
                    if not advmod_candidates: 
                        verb_idx = a.i
                        window = [j for j in sent if abs(j.i - verb_idx) <= 2]
                        for i, token in enumerate(window):
                            if token.text.lower() in claim_mod_terms and token.pos_ == "ADV":
                                span_tokens = [token]
                                for next_tok in window[i+1:]:
                                    if next_tok.text.lower() in claim_mod_terms and next_tok.pos_ == "ADV": 
                                        span_tokens.append(next_tok)
                                    else: 
                                        break
                                    advmod_candidates.append(span_tokens)

                    # select first candidate
                    if advmod_candidates: 
                        span = advmod_candidates[0]
                        if isinstance(span, list): 
                            claim_mod = " ".join([t.text for t in span])
                            claim_mod_start, claim_mod_end = get_new_token_offset(span[0], sent)[0], get_new_token_offset(span[-1], sent)[-1]
                        else: 
                            claim_mod = span.text
                            claim_mod_start, claim_mod_end = get_new_token_offset(span, sent)
                    else: 
                        claim_mod = ""
                        claim_mod_start, claim_mod_end = 0, 0


                    #     for j in sent: 
                    #         if abs(j.i - verb_idx) <= 2 and j.text.lower() in claim_mod_terms and j.pos_ == "ADV": 
                    #             advmod_candidates.append(j)
                    # if advmod_candidates: 
                    #     advmode = adv
                    # build claim_mod modifier

                    # if advmod: 
                    #     advmod = advmod.pop()
                    #     claim_mod = advmod
                    #     claim_mod_start, claim_mod_end = get_new_token_offset(claim_mod, sent)
                    #     claim_mod = advmod.text

                    #     # get content span and offset with claim_mod included
                    #     content_offsets = [get_new_token_offset(j, sent) for j in a.subtree if j.i > a.i and not j.is_punct]
                    # else: 

                    # get content span and offset
                    content_offsets = [get_new_token_offset(j, sent) for j in a.subtree if j.i > a.i and not j.is_punct]
                    # testing_content_text = " ".join([j.text for j in a.subtree if j.i > a.i])
                 
                    # testing_tuple = [source.text, verb.text, testing_content_text, claim_mod]
                   
                    if len(content_offsets) >= 2: 
                        content_start, content_end = content_offsets[0][0], content_offsets[-1][-1]
                        
                        # overlaps = [flag for flag in list((
                        #             (overlap(source_start, source_end, content_start, content_end)), 
                        #             (overlap(content_start, content_end, verb_start, verb_end))))
                        #             if flag == True]
                        
                        
                        # if not overlaps: 
                        
                        # new_sent = sent.text
                        
                        # starts = [new_sent[source_start], new_sent[verb_start], new_sent[content_start]]
                        
                        # ends = [new_sent[source_end], new_sent[verb_end], new_sent[content_end - 1]]
                        invalid_span = []
                        
                        if claim_mod: 
                            starts = [source_start, verb_start, content_start, claim_mod_start]
                            ends = [source_end, verb_end, content_end, claim_mod_end]
                            for start, end in zip(starts, ends): 
                                if not is_valid_span(start, end):
                                    invalid_span.append(is_valid_span(start, end))
                            if not invalid_span: 
                                yield {
                                    "spans": {
                                        "sc": [
                                            (source_start, source_end, "SOURCE"),
                                            (verb_start, verb_end, "CLAIM_VERB"),
                                            (content_start, content_end, "CLAIM_CONTENTS"),
                                            (claim_mod_start, claim_mod_end, "CLAIM_MOD")
                                        ]
                                    }
                                }
                        else: 
                            starts = [source_start, verb_start, content_start]
                            ends = [source_end, verb_end, content_end]
                            for start, end in zip(starts, ends): 
                                if not is_valid_span(start, end):
                                    invalid_span.append(is_valid_span(start, end))
                            if not invalid_span: 
                                yield {
                                    "spans": {
                                        "sc": [
                                            (source_start, source_end, "SOURCE"),
                                            (verb_start, verb_end, "CLAIM_VERB"),
                                            (content_start, content_end, "CLAIM_CONTENTS")
                                        ]
                                    }
                                }

# retrieves training data with token offsets
def get_training_list_spcat(): 
    train_data = [] # add tuples here
    final_text = "" # add entire claim sentence here
    annotations = [] # add offset tuples here

    # doc = initialize_doc(data)
    
    # doc = initialize_doc(example_text)
    # doc = initialize_doc(
    #     'Several initiatives – such as AI4All and the AI Now Institute – explicitly '
    #     'advocate for fair, diverse, equitable, and non-discriminatory inclusion in '
    #     'AI at all stages, with a focus on support for under-represented groups.'
    # )
    # doc = initialize_doc()
    claim_count = 0
    sent_count = 0
    seen_sents = []
    # for doc in nlp.pipe(pull_all_files()): 
    # testing
    for doc in nlp.pipe(pull_n_files(5)):
        sent_count += len(list(doc.sents))
        
        for ent in doc.ents:
            
            if ent.label_ in target_sources:
                source_span = ent.root.sent
            
                if source_span not in seen_sents: 
                    seen_sents.append(source_span)

                    
                    
                    for entry in get_tuples_spcat(source_span): 
                        
                        if entry: 
                            final_text = preprocess_text(source_span.text)
                            claim_count += 1
                            annotations = entry
                            data = ((final_text, annotations))
                            train_data.append(data)
    for data in train_data: 
        print(f"\n\n{data}\n\n")
    return train_data

def get_training_list_ner(): 
    train_data = [] # add tuples here
    final_text = "" # add entire claim sentence here
    annotations = [] # add offset tuples here

    # doc = initialize_doc(data)
    
    # doc = initialize_doc(example_text)
    # doc = initialize_doc(
    #     'Several initiatives – such as AI4All and the AI Now Institute – explicitly '
    #     'advocate for fair, diverse, equitable, and non-discriminatory inclusion in '
    #     'AI at all stages, with a focus on support for under-represented groups.'
    # )
    # doc = initialize_doc()
    claim_count = 0
    sent_count = 0
    seen_sents = []
    for doc in nlp.pipe(pull_n_files(10)): 
        sent_count += len(list(doc.sents))
        
        for ent in doc.ents:
            
            if ent.label_ in target_sources:
                source_span = ent.root.sent
            
                if source_span not in seen_sents: 
                    seen_sents.append(source_span)

                    
                    
                    for entry in get_tuples_ner(source_span): 
                        
                        if entry: 
                            final_text = preprocess_text(source_span.text)
                            claim_count += 1
                            annotations = {"entities": entry}
                            data = ((final_text, annotations))
                            train_data.append(data)
    return train_data
    

# get token offset from local sentence
def get_new_token_offset(token, sent): 
    rel_start = token.idx - sent.start_char
    rel_end = rel_start + len(token.text)
    return rel_start, rel_end


# NEED TO REFACTOR TO GET STARTING AND ENDING INDICES INSTEAD OF OFFSETS 


# !!!!!!!!!!!!!
# In spaCy, if a Span consists of only one token, the start attribute of that Span will be 
# the index of that singular token, and the end attribute will be the index of the token after it.
# !!!!!!!!!!!!!


# gets spans based on token offsets
def get_tuples_spcat(sent): 
    source_start, source_end = 0, 0
    verb_start, verb_end = 0, 0
    content_start, content_end = 0, 0
    claim_mod_start, claim_mod_end = 0, 0
    claim_mod = ""
    
    # testing - making new doc from sentence arg 
    doc = nlp(sent.text)
    # for source in sent.ents: 
    for source in doc.ents:
        if source.label_ in target_sources: 
            
    
            # new start of entity offset 
            source_start = source.start
            source_end = source.end
            # source_start = ent_start
            # source_end = ent_end
          
            # NOTE: we have found a valid claim
            for a in source.root.ancestors: 
                if a.lemma_ in claim_verb_terms and a.pos_ == "VERB":
                    
                    verb = a
                  
                    verb_start = verb.i
                    verb_end = verb.i + 1

                    # checks for adverb modifier that indicates degree of claim (not perfect but hopefully catches some)
                    # if the modifier is in the children of the claim verb, is an adverb, ends in -ly, and is directly before the verb
                    advmod = [j for j in a.children if j.dep_ == "advmod" and j.text in claim_mod_terms and j.i == a.i - 1]
                
                    # build claim_mod modifier
                    if advmod: 
                        advmod = advmod.pop()
                        claim_mod = advmod
                        claim_mod_start = claim_mod.i
                        claim_mod_end = claim_mod.i + 1
                        claim_mod = advmod.text

                        # get content span and offset with claim_mod included
                        content_i = [j for j in a.subtree if j.i > a.i and j.i > advmod.i and not j.is_punct]
                    else: 
                        # get content span and offset
                        content_i = [j for j in a.subtree if j.i > a.i and not j.is_punct]
                    testing_content_text = " ".join([j.text for j in a.subtree if j.i > a.i])
                 
                    testing_tuple = [source.text, verb.text, testing_content_text, claim_mod]
                   
                    if len(content_i) >= 2: 
                        content_start = content_i[0].i
                        
                        content_end = content_i[-2].i 
                        # content_start, content_end = content_offsets[0][0], content_offsets[-1][-1]
                        
                        # overlaps = [flag for flag in list((
                        #             (overlap(source_start, source_end, content_start, content_end)), 
                        #             (overlap(content_start, content_end, verb_start, verb_end))))
                        #             if flag == True]
                        
                        # if not overlaps: 
                    
                        new_sent = sent.text
                        
                        # starts = [new_sent[source_start], new_sent[verb_start], new_sent[content_start]]
                        
                        # ends = [new_sent[source_end], new_sent[verb_end], new_sent[content_end - 1]]
                        if claim_mod: 
                            
                            yield {
                                "spans": {
                                    "sc": [
                                        (source_start, source_end, "SOURCE"),
                                        (verb_start, verb_end, "CLAIM_VERB"),
                                        (content_start, content_end, "CLAIM_CONTENTS"),
                                        (claim_mod_start, claim_mod_end, "CLAIM_MOD")
                                    ]
                                }
                            }
                        else: 
                            yield {
                                "spans": {
                                    "sc": [
                                        (source_start, source_end, "SOURCE"),
                                        (verb_start, verb_end, "CLAIM_VERB"),
                                        (content_start, content_end, "CLAIM_CONTENTS")
                                    ]
                                }
                            }




# check to see if span is valid (start < end)
def is_valid_span(start, end): 
    return start < end

# get tuples from a target sentence
# sentence: a span that contains the ent, sources: a list of target sources
def get_tuples_ner(sent):
    source_start, source_end = 0, 0
    verb_start, verb_end = 0, 0
    content_start, content_end = 0, 0
    claim_mod_start, claim_mod_end = 0, 0
    claim_mod = ""
    

    for source in sent.ents: 

        if source.label_ in target_sources: 
            
            # NOTE: dont think i need these
            # new start of sentence offset 
            sent_start = sent.end_char - sent.end_char 
            sent_end = sent.end_char - sent.start_char

            # new start of entity offset 
            source_start = source.start_char - sent.start_char
            source_end = source.end_char - sent.start_char
            # source_start = ent_start
            # source_end = ent_end
          
            # NOTE: we have found a valid claim
            for a in source.root.ancestors: 
                if a.lemma_ in claim_verb_terms and a.pos_ == "VERB":
                    
                    verb = a
                  
                    verb_start, verb_end = get_new_token_offset(verb, sent)

                    # checks for adverb modifier that indicates degree of claim (not perfect but hopefully catches some)
                    # if the modifier is in the children of the claim verb, is an adverb, ends in -ly, and is directly before the verb
                    advmod = [j for j in a.children if j.dep_ == "advmod" and "ly" in j.suffix_ and j.i == a.i - 1]
                
                    # build claim_mod modifier
                    if advmod: 
                        advmod = advmod.pop()
                        claim_mod = advmod
                        claim_mod_start, claim_mod_end = get_new_token_offset(claim_mod, sent)
                        claim_mod = advmod.text

                        # get content span and offset with claim_mod included
                        content_offsets = [get_new_token_offset(j, sent) for j in a.subtree if j.i > a.i and j.i > advmod.i and not j.is_punct]
                    else: 
                        # get content span and offset
                        content_offsets = [get_new_token_offset(j, sent) for j in a.subtree if j.i > a.i and not j.is_punct]
                    testing_content_text = " ".join([j.text for j in a.subtree if j.i > a.i])
                 
                    testing_tuple = [source.text, verb.text, testing_content_text, claim_mod]
                   
                    if len(content_offsets) >= 2: 
                        content_start, content_end = content_offsets[0][0], content_offsets[-1][-1]
                        
                        overlaps = [flag for flag in list((
                                    (overlap(source_start, source_end, content_start, content_end)), 
                                    (overlap(content_start, content_end, verb_start, verb_end))))
                                    if flag == True]
                        
                        if not overlaps: 
                        
                            new_sent = sent.text
                          
                            starts = [new_sent[source_start], new_sent[verb_start], new_sent[content_start]]
                            
                            ends = [new_sent[source_end], new_sent[verb_end], new_sent[content_end - 1]]
                            if claim_mod: 
                                
                                yield [
                                    (source_start, source_end, "SOURCE"),
                                    (verb_start, verb_end, "CLAIM_VERB"),
                                    (content_start, content_end, "CLAIM_CONTENTS"),
                                    (claim_mod_start, claim_mod_end, "CLAIM_MOD")
                                ]
                            else: 
                                yield [
                                    (source_start, source_end, "SOURCE"),
                                    (verb_start, verb_end, "CLAIM_VERB"),
                                    (content_start, content_end, "CLAIM_CONTENTS"),
                                ]

# check if entities overlap   
def overlap(start1, end1, start2, end2): 
    return end1 >= start2 and end2 >= start1

# # NOTE: don think i need this
# # a function that identifies claim verbs and direct objects that are grammatically linked to a source
# def source_to_claim(source): 
#     claim_phrase = []
#     claim_verb = ""
#     claim_subtree = ""
#     strength_phrase = ""
#     # iterate through all ancestors of the source token
#     for a in source.ancestors: 
#         # when you get to a claim verb
#         if a.lemma_ in claim_verb_terms and a.pos_ == "VERB": 
            
#             # get claim verb
#             claim_verb = a.lemma_

#             # checks for adverb modifier that indicates degree of claim (not perfect)
#             advmod = [j for j in a.children if j.dep_ == "advmod" and "ly" in j.suffix_]
            
#             # this does nothing right now 
#             prep = [j for j in a.children if j.dep_ == "prep"]

#             # build claim_mod modifier
#             if advmod: 
#                 advmod = advmod.pop()
#                 strength_phrase = advmod.text
    
#             # all descendants to the right of the root/claim verb that 
#             # Maybe modify to get multiple claims within one sentence? 
#             claim_subtree = " ".join([j.text for j in a.subtree if j.i > a.i and not j.is_punct])
            
#             break
    
#     return claim_verb, claim_subtree, strength_phrase


# def test_one_file(file=""): 

#     # if file not supplied
#     if not file: file = preprocess_text(pull_one_file())
#     doc = nlp(file)
#     claims = []
#     relations = dict()
#     explicit_claims = []
#     sources = ["ORG", "PERSON", "GPE", "NORP", "LAW"]
#     strength_phrase_count = 0 # testing
#     processed_sents_count = 0 # testing
#     processed_sents = [] # testing
#     for ent in doc.ents: 
#         # print(ent.text, ent.label_) # testing
#         # print(ent.root.text) # testing
#         if ent.label_ in sources:
#             claim_verb, claim_contents, strength_phrase = source_to_claim(ent.root)
#             if strength_phrase: strength_phrase_count += 1 # testing
#             if claim_contents:
#                 processed_sents_count += 1 # testing
#                 processed_sents.append(ent.root.sent.text) # testing
#                 # print(f"\nProcessed sentence {processed_sents_count}: {ent.root.sent.text}") # testing 
#                 # print(ent.text)
#                 explicit_claims.append((ent.text, claim_verb, claim_contents, strength_phrase))
 
#     pprint.pprint(explicit_claims)
#     print(f"Number of claims identified: {len(explicit_claims)}" )

# def test_all_files(): 
#     for i, doc in enumerate(nlp.pipe(pull_all_files())):

#         claims = []
#         relations = dict()
#         explicit_claims = []
#         sources = ["ORG", "PERSON", "GPE", "NORP", "LAW"]
#         strength_phrase_count = 0 # testing
#         processed_sents_count = 0 # testing
#         processed_sents = [] # testing
#         for ent in doc.ents: 
#             if ent.label_ in sources: 
#                 claim_verb, claim_contents, strength_phrase = source_to_claim(ent.root)
#                 if strength_phrase: strength_phrase_count += 1 # testing 
#                 if claim_contents :
#                     processed_sents_count += 1 # testing
#                     processed_sents.append(ent.root.sent.text)
#                     explicit_claims.append((ent.text, claim_verb, claim_contents, strength_phrase))
#         print("Strength phrase count: ", strength_phrase_count) # testing
#         pprint.pprint(explicit_claims)
#         pprint.pprint(processed_sents)
        
def see_relations(text):
    doc = nlp(text)
    displacy.serve(doc, style="dep", auto_select_port=True)


def main():
    # ------------------------------------------------------------------------------------------------
    # Where we're at: 
    # ------------------
    # - CURRENT STRUCTURE: [source] [claim_verb] [claim_contents] | [claim_strength]
    # - was able to use dependency parsing to retrieve source, claim verb (defined by our claim_verb_terms), claim contents, and 
    #   SOMETIMES claim claim_mod. 
    # - the rules are very loose and this wont work 100% of the time
    # - for most sentences that make a single claim and their source is an ORG, PERSON, or GPE, 
    #   this algorithm will work at retrieving at least a triplet of data describing the claim
    # - if an adverb is used to describe the claim verb and it ends in ly, it will be added to claim claim_mod.
    #   - obviously this isnt perfect as there are other words like "only" that also fit this description 
    #   - and im sure there are adverbs that dont end in ly that i should be looking for 
    # - if multiple claims are made in a sentence, then only the first claim verb is parsed as the claim verb, but the remainder
    #   of the claim contents should end up in the claim_contents 
    # - does not recognized pronouns, would have to use coreferencing in order to accomplish this 
    # - will probably want to give a score to the claim claim_mod, need to figure out how to do this


    # see_relations("Miller however suggests that the user should be cautious with AI.")
    # see_relations("The UK relunctantly and with great hesitation suggests sanctions on AI.")
    # see_relations("The UK states that it is reluctant with AI and it also claims that it is destructive")
    # test_one_file("The UK states that it is reluctant with AI and it also claims that it is destructive")
    # test_one_file("The UK states that it is reluctant with AI. It also claims that it is destructive")
    # see_relations("The UK states that it is reluctant with AI and it also claims that it is destructive")
    # test_one_file()
    # see_relations("""Google’s translate system can suffer from gender bias by making sentences taken from the
    #                 U.S. Bureau of Labor Statistics into a dozen languages that are gender neutral, including Yoruba,
    #                 Hungarian, and Chinese, translating them into English, and showing that Google Translate shows favoritism toward males for stereotypical fields such as STEM jobs.""")
#     see_relations('Several initiatives – such as AI4All and the AI Now Institute – explicitly '
#   'advocate for fair, diverse, equitable, and non-discriminatory inclusion in '
#   'AI at all stages, with a focus on support for under-represented groups.')
    # see_relations("According to John Smith, AI will revolutionize energy production by 2030.")
    # see_relations("According to John Smith, AI will revolutionize energy production by 2030. According to the BBC, global food shortages will increase in the next decade. According to Dr. Robert Chen, AI ethics will become a central concern for policymakers.  According to Reuters, interest rates will likely rise next year. According to CNN, climate migration will increase in the coming years. According to Bloomberg, quantum computing will transform finance. According to the New York Times, AI bias remains a major challenge.  According to Wired, cybersecurity threats will evolve rapidly. According to Fox News, renewable energy adoption will accelerate. According to NPR, social media influences public opinion strongly. According to The Guardian, climate resilience is key to sustainable development. According to Financial Times, global trade patterns are shifting. According to Vox, misinformation spreads faster than truth.")
    # test_one_file("In an article from CNN, it was stated that AI is great. CNN reported that AI is not good. According to John Smith, AI will revolutionize energy production by 2030. According to the BBC, global food shortages will increase in the next decade. According to Dr. Robert Chen, AI ethics will become a central concern for policymakers.  According to Reuters, interest rates will likely rise next year. According to CNN, climate migration will increase in the coming years. According to Bloomberg, quantum computing will transform finance. According to the New York Times, AI bias remains a major challenge.  According to Wired, cybersecurity threats will evolve rapidly. According to Fox News, renewable energy adoption will accelerate. According to NPR, social media influences public opinion strongly. According to The Guardian, climate resilience is key to sustainable development. According to Financial Times, global trade patterns are shifting. According to Vox, misinformation spreads faster than truth.")
    # see_relations("The French government states that renewable energy will be the primary power source. The Japanese government states that electric vehicles will dominate the market by 2040.")
    # test_one_file("The French government states that renewable energy will be the primary power source. The Japanese government states that electric vehicles will dominate the market by 2040.")
    # get_training_list_ner()
    get_training_list_spcat()
    
if __name__ == "__main__": 
    main() 