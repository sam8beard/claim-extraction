import spacy 
nlp = spacy.load("en_core_web_trf")

def main():
    doc = nlp(
            "Artificial Intelligence (AI) is revolutionizing the healthcare industry by enhancing diagnostic accuracy, personalizing treatment plans, and streamlining administrative tasks. "
            "IBM Watson Health and Google Health are leading this transformation. "
            "Machine learning algorithms analyze vast amounts of medical data from institutions like Mayo Clinic and Johns Hopkins Hospital to identify patterns that may be missed by human clinicians. "
            "For instance, AI-driven imaging tools such as Aidoc Radiology can detect early signs of diseases such as lung cancer and Alzheimer’s disease, enabling timely interventions. "
            "In 2025, these technologies are expected to be integrated into over 60 percent of hospitals worldwide."
    )
    print("Sentences: \n")
    for sent in doc.sents: 
        print(sent, "\n")
    print("Token details: \n")
    for token in doc: 
        print(token.text, token.lemma_, token.pos_, "\n")
    print("Named entities: \n")
    for ent in doc.ents: 
        print(ent.text, ent.label_, "\n")
        
    # Add Artifical Intelligence and AI as entities using PhraseMatcher 
    
if __name__ == "__main__": 
    main() 