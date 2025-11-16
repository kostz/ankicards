# Ankicards

Tool leverages LLM capability to extract German verbs out of schoolbook scans to lately create Anki cards for a private German language study.

Schoolbook contains of irregular German verbs for levels A1 and A2, which are extracted, translated into English and Russian and enriched with a complete example sentences.

## Data Fields

Each verb record contains below data fields

- Infinitive
- Present form + example
- Praeteritum form + example

Both Infinitive form and Examples have translation into English and Russian

## Usage

### Initialize variables

Use your beloved AI subscription to create below file to source it it

```sh
#!/bin/bash

export LLM_API_KEY="<LLM API Key>"
export LLM_BASE_URL="<LLM API Url>"
export LLM_MODEL="<LLM Model>"
```

### Extract initial data from pages

```shell
make extract-verbs
```

### Enrich data with examples

```shell
make add-verb-examples
```

### Create Anki Cards (Not ready yet)

```shell
make create-anki-cards
```
