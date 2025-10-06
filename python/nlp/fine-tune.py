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
from nlp_pipeline import get_training_list
from spacy.training.example import Example
import random
from spacy.training import offsets_to_biluo_tags

logging.basicConfig(level=logging.INFO, stream=sys.stderr)

# load model 
nlp = spacy.load("en_core_web_sm")
# def get_training_data(): 
#     training_data = get_training_list()
#     return training_data
# get ner component 
ner = nlp.get_pipe("ner")

# add labels to ner
labels = ["SOURCE", "CLAIM_VERB", "CLAIM_CONTENTS", "CLAIM_STRENGTH"]
for label in labels: 
    ner.add_label(label)

def fine_tune(): 

    

    # convert training data to Examples
    examples = []
    training_data = get_training_list()
    # training_data = align_offsets_to_tokens()
    # for text, annots in training_data: 
    #     # logging.info(f"Text once passed to fine tune:\n {text}")
    #     for start, end, label in annots["entities"]:
    #         logging.info(f"{start, end, label}")
    #         span = text[start:end]
    #         if not span.strip():
    #             logging.info(f"Empty span in: {text[start:end]!r}")
    #         if end > len(text):
    #             logging.info(f"Span out of range: {text[start:end]!r}")
    #         if span not in text:
    #             logging.info(f"Span text mismatch: '{span}' not found in '{text}'")
    #         # doc = nlp.make_doc(text)
    #     # example = Example.from_dict(doc, annots)
    #     # examples.append(example)


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

def align_offsets_to_tokens():
    train_data = get_training_list()
    aligned = []
    for text, annots in train_data:
        doc = nlp.make_doc(text)
        entities = []
        for start, end, label in annots.get("entities", []):
            span = doc.char_span(start, end, label=label, alignment_mode="contract")
            if span is None:
                print(f"Cannot align entity '{text[start:end]}' in: {text}")
            else:
                entities.append((span.start_char, span.end_char, label))
        aligned.append((text, {"entities": entities}))
    # logging.info(f"{aligned}")
    return aligned

def debug(): 
    claim_count = 0
    valid_count = 0
    train_data = get_training_list() 
    for text, annots in train_data:
        claim_count += 1 
        doc = nlp.make_doc(text)
        entities = annots.get("entities", [])
        tags = offsets_to_biluo_tags(doc, entities)
        
        if tags is None: 
            logging.info(f"Alignment error")
        else:
            logging.info(f"Tags: {tags}") 
            valid_count += 1
    logging.info(f"Claims: {claim_count}")
    logging.info(f"Valid claims: {valid_count}")

def see_results(): 
    nlp_updated = spacy.load("ner_v1.0")
    # doc = nlp_updated(pull_one_file())
    num = 5
    for file in pull_n_files(num): 
        doc = nlp_updated(file)
        results = [(ent.label_, ent.text) for ent in doc.ents]
    
        for result in results: 
            pprint.pprint(f"{result}\n\n")
        

# see_results()
# debug()
# fine_tune()
see_results()
# align_offsets_to_tokens()

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