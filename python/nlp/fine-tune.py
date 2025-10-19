from spacy.pipeline.spancat import DEFAULT_SPANCAT_MODEL
from pathlib import Path
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
from nlp_pipeline import get_training_list_ner, get_training_list_spcat, claim_verb_terms, claim_mod_terms, get_training_list_spcat_2
from spacy.training.example import Example
import random
from spacy.training import offsets_to_biluo_tags
from spacy.lookups import load_lookups
from collections import Counter
from spacy.pipeline.spancat import Suggester
from thinc.types import Ragged, Ints1d
from typing import Optional, Iterable, cast
from thinc.api import get_current_ops, Ops
from spacy import registrations
from thinc.api import Config
from spacy.lang.en import English
from spacy.util import load_config
from tqdm import tqdm 

# configure logger for testing
logging.basicConfig(level=logging.INFO, stream=sys.stderr)

# labels for training
labels = ["SOURCE", "CLAIM_VERB", "CLAIM_CONTENTS", "CLAIM_MOD"]

# def print_label_count(): 
#     training_data = get_training_list_spcat()

#     counts = Counter(a[2] for _, ann in training_data for a in ann["spans"]["sc"])
#     max_count = max(counts.values())

#     grouped = {label: [] for label in counts}
#     for ex in training_data: 
#         _, ann = ex
#         labels = [a[2] for a in ann["spans"]["sc"]]
#         if labels: 
#             grouped[labels[0]].append(ex)
#     balanced = []
#     for label, examples in grouped.items(): 
#         if not examples: 
#             continue
#         needed = max_count - len(examples)
#         balanced.extend(examples + random.choices(examples, k=max(0, needed)))
#     print(counts)
#     random.shuffle(balanced)
#     new_counts = Counter(a[2] for _, ann in balanced for a in ann["spans"]["sc"])
#     print(balanced)
    # print(new_counts)
    # training_data = get_training_list_spcat()
    # label_counts = Counter() 
    # balanced_data = []
    # for text, annots_dict in training_data:
    #     # print(text)
    #     # print(annots_dict)
    #     spans = annots_dict["spans"]["sc"]
    #     for i, a in enumerate(spans):
    #         # print(a)
    #         if "CLAIM_STRENGTH" in a or "CLAIM_CONTENTS" in a:
    #             tuples = a * 3
    #             # print(tuples)
    #             spans.append(tuples)
    #             # balanced_data.append((text, spans.append(tuples)))
    #         else: 
    #             spans.append(a)
    #             # balanced_data.append((text, spans.extend(a * 3)))
    #             # balanced_data.append((text, spans.extend(a * 3)))
    #         # label = a[2]
    #         # print(a[2])
    #         # if a[2] == "CLAIM_CONTENTS" or "CLAIM_STRENGTH":
           
    #         #     balanced_data.append((text, spans.extend(a * 3)))
    #         # else: 
    #         #     balanced_data.append((text, spans))
    #         # # print(a)
    #     balanced_data.append(text, spans)
    #         # label_counts[a[2]] += 1
    # print(balanced_data)
    # print(label_counts)
    # examples = [annots['spans']['sc'] for _, annots in training_data]
    # for i, group in enumerate(examples):
    #     for j, tup in enumerate(group): 
    #         labels = [label for label in tup]
    #         # print(f"\nLabel num: {j}\t{labels[2]}\n")
    #         print(labels)
    #         if "CLAIM_CONTENTS" in labels or "SOURCE" in labels:
    #             print("Here:", labels)
    #             print(f"Contents found : {tup}" )

    #     # label_counts.update([a["label"] for a in annots])
    # print(label_counts)

# def balance_training_data(training_data):
#     balanced_data = []
#     for example in training_data: 
#         labels = [a["label"] for a in example[1]]
#         if "CLAIM_CONTENTS" in labels or "CLAIM_STRENGTH" in labels: 
#             balanced_data.extend([example] * 3)
#         else: 
#             balanced_data.append(example) 
#     return balanced_data


def fine_tune_ner(): 
    # load and configure model
    nlp = spacy.load("en_core_web_sm")
    nlp.initialize()
    ner = nlp.get_pipe("ner")
    for label in labels: nlp.add_pipe(label)

    # convert training data to Examples
    examples = []
    training_data = get_training_list_ner()

    # adjust as needed 
    other_pipes = [pipe for pipe in nlp.pipe_names if pipe != "ner"]
    with nlp.disable_pipes(*other_pipes):
        optimizer = nlp.resume_training()
        epochs = 30
        for itn in range(epochs): 
            random.shuffle(examples)
            losses = {}
            batches = spacy.util.minibatch(training_data, size=32)
            for batch in batches: 
                examples = []
                for text, annots in batch: 
                    # logging.info(f"Current annots: {annots}")
                    doc = nlp.make_doc(text)
                    example = Example.from_dict(doc, annots)
                    examples.append(example)
                nlp.update(examples, drop=0.3, losses=losses)
            logging.info(f"Iteration: {itn + 1}, Losses: {losses}")
    logging.info(f"Number of examples processed: {len(training_data)}")

    nlp.to_disk('ner_v1.0')

# will be used as the suggester function
@spacy.registry.misc("claim_suggester.v1") 
def build_claim_suggester() -> Suggester: 
    def claim_suggester(
        docs: Iterable[Doc], *, ops: Optional[Ops] = None
    ) -> Ragged : 
        if ops is None: 
            ops = get_current_ops()
        target_sources = ["ORG", "PERSON", "GPE", "NORP", "LAW"]
        claim_verb_terms = [
            "argue", "assert", "claim", "contend", "propose", "maintain", "state", "suggest", "insist", "hold",
            "affirm", "demonstrate", "establish", "prove", "show", "validate", "substantiate", "confirm", "illustrate",
            "justify", "interpret", "evaluate", "analyze", "examine", "assess", "deduce", "infer", "theorize", "posit",
            "hypothesize", "advocate", "recommend", "encourage", "urge", "emphasize", "call for", "champion", "challenge",
            "dispute", "question", "refute", "counter", "reject", "critique", "oppose", "problematize", "believe", "warn"
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
        spans = []
        lengths = []
        for doc in docs: 

            if doc.has_annotation("DEP"):
                logging.info("Firing Here")
                claims = []
                
                # begin looking for claims based off of existence of valid source
                for ent in doc.ents:
                    logging.info("Firing here")
                   
                    source_tuple = None
                    verb_tuple = None
                    mod_tuple = None
                    content_tuple = None
                    if ent.label_ in target_sources: 
                        source_tuple = (ent.start, ent.end)
                        logging.info(source_tuple)
                        # logging.info(source_tuple) # not firing
                        for a in ent.root.ancestors:
                            # need to reference claim_verb_terms
                            logging.info(a.lemma_)
                            logging.info(a.text)
                            if a.lemma_ in claim_verb_terms and a.pos_ == "VERB":
                                logging.info("Firing on claim_verb_terms")
                                # possible claims verb 
                                a_start = a.i
                                a_end = a.i + 1
                                verb_tuple = [a_start, a_end]
                                # possible claims strength
                                claim_mod = [j for j in a.children if j.dep_ == "advmod" and j.text in claim_mod_terms]
                                
                                if claim_mod: # claim_strength found
                                    claim_mod = claim_mod[0]
                                    mod_start = claim_mod.i
                                    mod_end = claim_mod.i + 1
                                    mod_tuple = [mod_start, mod_end]
                                    # scan for content after verb and advmod
                                    content = [j for j in a.subtree if j.i > a.i and j.i > claim_mod.i and not j.is_punct]
                                    if len(content) >= 2:
                                        content_start = content[0].i
                                        content_end = content[-2].i
                                        if content_start < content_end: 
                                            content_tuple = [content_start, content_end]

                                else:  # no claims strength
                                    # scan for content after verb
                                    content = [j for j in a.subtree if j.i > a.i and not j.is_punct]
                                    if len(content) >= 2:
                                        content_start = content[0].i
                                        content_end = content[-2].i
                                        if content_start < content_end: 
                                            content_tuple = [content_start, content_end]

                                if source_tuple and verb_tuple and mod_tuple and content_tuple: 
                                    claims.append(source_tuple)
                                    claims.append(verb_tuple)
                                    claims.append(mod_tuple)
                                    claims.append(content_tuple)

                                elif source_tuple and verb_tuple and content_tuple and not mod_tuple: 
                                    claims.append(source_tuple)
                                    claims.append(verb_tuple)
                                    claims.append(content_tuple)
                claims = ops.asarray(claims, dtype="i")
                if claims.shape[0] > 0: 
                        spans.append(claims)
                        lengths.append(claims.shape[0])
                else: 
                    logging.info("Does not have DEP")
                    lengths.append(0)
            else: 
                lengths.append(0)
        lengths_array = cast(Ints1d, ops.asarray(lengths, dtype="i"))
        if len(spans) > 0: 
            output = Ragged(ops.xp.vstack(spans), lengths_array)
        else: 
            output = Ragged(ops.xp.zeros((0,0), dtype="i"), lengths_array)
        assert output.dataXd.ndim == 2
    
        return output
    
    return claim_suggester


def fine_tune_spcat(): 
    # # load and configure model
    # config = Config().from_disk("config.cfg")
    # nlp = spacy.Language.from_config(config)
    # nlp.from_config(config_path="config.cfg")
    # nlp = spacy.load("./config.cfg")
    # config_path = Path("./config.cfg")
    # config = load_config(config_path)
    # config = config.to_bytes
    # nlp = spacy.load("en_core_web_sm")
    # nlp = spacy.load(config)

    # add components to frozen and annotating lists for custom suggester
    training_config = { 
        "training": { 
            "frozen_components": ["tok2vec","tagger","parser","attribute_ruler", "ner"],
            "annotating_components": ["tok2vec","tagger","parser","attribute_ruler", "ner"]
        }
    }

    nlp = spacy.blank('en', config=training_config)
    base_nlp = spacy.load('en_core_web_sm')
    
    # add components to blank model sourced from en_core_web_sm
    for name in ["tok2vec", "tagger", "parser", "attribute_ruler", "ner"]:
        nlp.add_pipe( 
            name, 
            source = base_nlp
        )
    
    
        # comp_name = name
        # logging.info("firing")
        # logging.info(base_nlp.get_pipe(name).model)
        # comp = base_nlp.get_pipe(name)
        # nlp.add_pipe(base_nlp.get_pipe(comp).model, name=comp)
    

    # nlp = spacy.load("en_core_web_sm")
    # logging.info(nlp.component_names)

    # PATH FORWARD
    # have to figure out how to build config so the the frozen and annot components used in training
    # span cat are sourced from en_core_web_sm 


    # nlp.
    # config = { 
    #     "nlp": 
    # }

    spancat_config = {
            "spans_key": "sc",
            "suggester": {"@misc": "spacy.ngram_range_suggester.v1", "min_size": 1, "max_size": 20},
            "threshold": 0.8,
            "max_positive": None, 
            # "frozen_components": ["tok2vec", "tagger", "parser", "attribute_ruler", "ner"],
            # "annotating_components": ["tok2vec", "tagger", "parser", "attribute_ruler", "ner"],
            "model": DEFAULT_SPANCAT_MODEL,
            
            # "training": {
            #     "frozen_components": [
            #         "tok2vec","tagger","parser","attribute_ruler", "ner"
            #     ],
            #     "annotating_components": [
            #         "tok2vec","tagger","parser","attribute_ruler", "ner"
            #     ]
            # } 


    }

    
    nlp.add_pipe("spancat", config=spancat_config)
    spancat = nlp.get_pipe('spancat')
    # nlp.analyze_pipes(pretty=True)
    pprint.pprint(nlp.config)
    # logging.info(nlp.analyze_pipes(pretty=True))
    # logging.info(spancat.getmembers())
    # nlp.add_pipe("training")
    # logging.info(nlp.training)
    return

    for label in labels: spancat.add_label(label)

    pipe_exceptions = ["spancat"]
    unaffected_pipes = [pipe for pipe in nlp.pipe_names if pipe not in pipe_exceptions]
    
    nlp.initialize()
    sgd = nlp.create_optimizer()

    # convert training data to Examples
    # training_data = balance_training_data(get_training_list_spcat())
    
    # NOTE: trying other training list function that returns char offsets
    #
    #       get_training_list_spcat returns TOKEN offsets, but Example.from_dict()
    #       expects CHAR offsets
    # training_data = get_training_list_spcat()

    training_data = get_training_list_spcat_2()
    examples = []
    
    for text, annots in training_data: 
        doc = nlp.make_doc(text)
        logging.info(len(doc))
        logging.info(text)
        logging.info(annots)

        # NOTE: issue here with token indices not aligning
        examples.append(Example.from_dict(doc, annots))

    with nlp.disable_pipes(*unaffected_pipes): 
        for itn in tqdm(range(40)):
            random.shuffle(examples)
            losses = {}
            batches = spacy.util.minibatch(examples, size=spacy.util.compounding(4.0, 32.0, 1.001))

            for batch in batches: 
                nlp.update(batch, losses=losses, drop= 0.1, sgd=sgd)
            logging.info(f"Iteration: {itn + 1}, Losses: {losses}")
    #     for text, annots in training_data: 
    #         # with nlp.select_pipes(disable="spancat"):
    #         # print(text)
    #         # print(annots)
    #         doc = nlp.make_doc(text)
    #         examples.append(Example.from_dict(doc, annots))
    #         # # logging.info(text)
    #         # # logging.info(annots) 
    #         # # logging.info(len(doc))
    #         # # doc.spans["sc"] = annots["spans"]
    #         # # logging.info(annots["spans"]["sc"])
    #         # spans_to_add = []
    #         # for span in annots["spans"]["sc"]:
    #         #     (start, end, label) = span
    #         #     # logging.info(f"{start}, {end}, {label}") 
    #         #     new_span = Span(doc, start, end, label)
    #         #     spans_to_add.append(new_span)
    #         #     # logging.info(f"New span: {new_span}")
    #         #     # spans_to_add.append(new_span)
    #         # # logging.info(spans_to_add)
    #         # # logging.info(doc.spans["sc"])
    #         # # for token in doc: 
    #         # #     logging.info(f"{token.text}, {token.i}")
    #         # logging.info(f"Spans object: {doc.spans}")
    #         # doc.spans["sc"] = spans_to_add
    #         # span_dict = {"spans": {"sc": spans_to_add}}
    #         # logging.info(span_dict)
    #         # logging.info(f"Spans key object: {doc.spans['sc']}")
    #         # # print(annots)
    #         # # examples.append(Example.from_dict(doc, {}))
    #         # new_doc = nlp.make_doc(text)
    #         # examples.append(Example.from_dict(new_doc, annots))

    # logging.info(f"Examples after processing: {examples[0:5]}")
    # # initialize model after adding labels and before training 
 
    # optimizer = nlp.initialize()
   
    # epochs = 30
    # for itn in range(epochs): 
    #     random.shuffle(examples)
    #     losses = {}
    #     batches = spacy.util.minibatch(examples, size=4)
    #     # for batch in batches: 
    #     #     nlp.update(batch, drop=0.2, losses=losses)
    #     for batch in batches: 

    #         # logging.info(f"Text: {example.text}")
    #         spancat.update(batch, drop=0.2, losses=losses, sgd=optimizer)
    #     logging.info(f"Iteration: {itn + 1}, Losses: {losses}")
    # logging.info(f"Number of examples processed: {len(training_data)}")

    nlp.to_disk('spancat_v1.0')

# 
def test_suggester():
    nlp = spacy.load("en_core_web_sm")
    # ner = nlp.get_pipe("ner") 
    # nlp.add_pipe("spancat")
    
    # print(nlp.pipe_names)
    suggester = build_claim_suggester()
    docs = []
    for text in pull_n_files(2):
        p_doc = nlp(text)
        for sent in p_doc: 
            doc = nlp.make_doc(sent.text, deps=["nsubj", "aux", "ROOT", "prep", "pcomp"])
            logging.info(doc.deps)
            docs.append(doc)
            # logging.info(doc)
        # docs.append(doc)
    # spancat = nlp.get_pipe("spancat")
    # spancat.initialize(docs)
    # logging.info(docs)
    result = suggester(docs)
    
    # logging.info(f"Results: {result}")

# testing custom ner model
def see_results_ner(): 
    nlp_updated = spacy.load("ner_v1.0")
    num = 40
    for file in pull_n_files(num): 
        doc = nlp_updated(file)
        results = [(ent.label_, ent.text) for ent in doc.ents]
    
        for result in results: 
            pprint.pprint(f"{result}\n\n")

# testing custom spancat model
def see_results_spcat(): 
    nlp_updated = spacy.load("spancat_v1.0")
    num = 5
    results = []
    label_counts = {"SOURCE": 0, "CLAIM_VERB": 0, "CLAIM_CONTENTS": 0, "CLAIM_MOD": 0}
    for file in pull_n_files(num):
        doc = nlp_updated(file)
        # for span in doc.spans['sc']: 
        #     if span.label_ == "SOURCE": 
        #         logging.info(f"\n\nSource Found: {span.text}")
        #         source_count += 1
        #     elif span.label_ == "CLAIM_VERB"
        logging.info("Doc Results -------------------------------\n")
        spans = doc.spans['sc']
        for span, confidence in zip(spans, spans.attrs["scores"]):
            label_counts[span.label_] += 1
            logging.info(f"{span.label_} | {confidence} | {span.text}")
        # for span in doc.spans['sc']: 
        #     if span.label_ == "SOURCE": 
        #         logging.info(f"\n\nSource Found: {span.text}")
        #         source_count += 1

            # logging.info(f"\n\nSpan: \n{span}\nSpan Label: {span.label_}\n")
        # logging.info(f"\n\nSOURCE count: {source_count}")
        
        print(f"\n\n")
    pprint.pprint(label_counts)
        # results = [(ent.label_, ent.text) for ent in doc.ents]

    # for result in results: 
    #     pprint.pprint(f"{result}\n\n")
    # logging.info("Firing")

# test_suggester()
# see_results()
# debug()
# fine_tune_ner()
# see_results()
# get_training_list_spcat_2() 
fine_tune_spcat()
# see_results_spcat()
# print_label_count()
# see_results_spcat()
# align_offsets_to_tokens()

# def align_offsets_to_tokens():
#     train_data = get_training_list()
#     aligned = []
#     for text, annots in train_data:
#         doc = nlp.make_doc(text)
#         entities = []
#         for start, end, label in annots.get("entities", []):
#             span = doc.char_span(start, end, label=label, alignment_mode="contract")
#             if span is None:
#                 print(f"Cannot align entity '{text[start:end]}' in: {text}")
#             else:
#                 entities.append((span.start_char, span.end_char, label))
#         aligned.append((text, {"entities": entities}))
#     # logging.info(f"{aligned}")
#     return aligned

# def debug(): 
#     claim_count = 0
#     valid_count = 0
#     train_data = get_training_list_ner() 
#     for text, annots in train_data:
#         claim_count += 1 
#         doc = nlp.make_doc(text)
#         entities = annots.get("entities", [])
#         tags = offsets_to_biluo_tags(doc, entities)
        
#         if tags is None: 
#             logging.info(f"Alignment error")
#         else:
#             logging.info(f"Tags: {tags}") 
#             valid_count += 1
#     logging.info(f"Claims: {claim_count}")
#     logging.info(f"Valid claims: {valid_count}")



# for next model, try 40 epochs, drop=0.35, optimizer.learn_rate = 0.0005

# losses over 30 epochs with 829 examples, drop = 0.3, batchsize = 32
# INFO:root:Iteration: 1, Losses: {'ner': np.float32(6035.229)}
# INFO:root:Iteration: 2, Losses: {'ner': np.float32(5059.371)}
# INFO:root:Iteration: 3, Losses: {'ner': np.float32(4288.8633)}
# INFO:root:Iteration: 4, Losses: {'ner': np.float32(3696.9487)}
# INFO:root:Iteration: 5, Losses: {'ner': np.float32(2620.0908)}
# INFO:root:Iteration: 6, Losses: {'ner': np.float32(1959.4371)}
# INFO:root:Iteration: 7, Losses: {'ner': np.float32(1724.4547)}
# INFO:root:Iteration: 8, Losses: {'ner': np.float32(1459.8206)}
# INFO:root:Iteration: 9, Losses: {'ner': np.float32(1427.7465)}
# INFO:root:Iteration: 10, Losses: {'ner': np.float32(1294.2133)}
# INFO:root:Iteration: 11, Losses: {'ner': np.float32(1219.9303)}
# INFO:root:Iteration: 12, Losses: {'ner': np.float32(1157.6885)}
# INFO:root:Iteration: 13, Losses: {'ner': np.float32(1112.0709)}
# INFO:root:Iteration: 14, Losses: {'ner': np.float32(1036.1287)}
# INFO:root:Iteration: 15, Losses: {'ner': np.float32(1034.9926)}
# INFO:root:Iteration: 16, Losses: {'ner': np.float32(1011.6324)}
# INFO:root:Iteration: 17, Losses: {'ner': np.float32(915.64374)}
# INFO:root:Iteration: 18, Losses: {'ner': np.float32(921.95074)}
# INFO:root:Iteration: 19, Losses: {'ner': np.float32(907.56305)}
# INFO:root:Iteration: 20, Losses: {'ner': np.float32(835.0009)}
# INFO:root:Iteration: 21, Losses: {'ner': np.float32(815.7493)}
# INFO:root:Iteration: 22, Losses: {'ner': np.float32(822.19446)}
# INFO:root:Iteration: 23, Losses: {'ner': np.float32(808.78516)}
# INFO:root:Iteration: 24, Losses: {'ner': np.float32(769.9993)}
# INFO:root:Iteration: 25, Losses: {'ner': np.float32(760.7867)}
# INFO:root:Iteration: 26, Losses: {'ner': np.float32(767.87103)}
# INFO:root:Iteration: 27, Losses: {'ner': np.float32(757.646)}
# INFO:root:Iteration: 28, Losses: {'ner': np.float32(749.51526)}
# INFO:root:Iteration: 29, Losses: {'ner': np.float32(719.44934)}
# INFO:root:Iteration: 30, Losses: {'ner': np.float32(762.30396)}
# INFO:root:Number of examples processed: 829