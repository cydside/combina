# combina - Combinatorial Password and Phrase Generator

## 📝 Overview

**combina** is a powerful combinatorial generator written in Go that creates all possible combinations and permutations of characters or words. Originally developed in C ([combina-0.4.1_src.tar.gz](https://github.com/cydside/combina/blob/main/combina-0.4.1_src.tar.gz) added as source code), it has been completely rewritten in Go to offer enhanced functionality, performance, and usability.
The software implements combinatorial analysis principles to exhaustively generate all possible strings from a given set of elements (characters or words) with specified lengths.

## ✨ Key Features

### 🎯 Generation Types

    Permutations with repetition (-r): Elements can be reused multiple times (e.g., AAA, AAB, ABA, ABB, ...)

    Permutations without repetition (-d): Each element used only once (e.g., ABC, ACB, BAC, ...)

    Combinations with repetition (-m): Order does not matter, elements reusable (e.g., AAA, AAB, ABB, BBB, ...)

    Combinations without repetition (-c): Order does not matter, unique elements (e.g., ABC, ABD, ACD, ...)

### 🔤 Input Modes

    Character Mode: Use predefined or custom character sets

        Lowercase letters (-a)

        Uppercase letters (-A)

        Numbers (-n)

        Special characters (-s)

        Custom set (-user)

    Phrase Mode: Work with words instead of individual characters

        Splits phrases into words

        Supports custom separators for input and output

        Ideal for generating passphrases or word combinations

### 🚀 Performance and Concurrency

    Multi-threading support with configurable worker pools (-workers)

    Parallel generation to maximize CPU utilization

    Batch processing for handling large volumes of combinations

    Real-time progress monitoring (-verbose)

### 🔐 Output Formats

    Plain: Simple text output (default)

    MD5: MD5 hash of the combination

    SHA1: SHA1 hash of the combination

    HEX: Hexadecimal representation

### 🎨 Customization

    Prefix/Suffix: Add custom strings before and after each combination

    Progressive Length: Generate combinations of varying lengths (-p)

    Dual Separators: Independent separators for input (split) and output (join)

### 💡 Use Cases

    Security Testing: Dictionary generation for penetration testing

    Password Recovery: Custom wordlist creation

    Linguistics: Study of word combinations

    Education: Teaching combinatorial principles

    Name Generation: Systematic creation of usernames or identifiers

    Software Testing: Exhaustive input generation for testing

### 🛠️ Practical Examples
bash

### 6-character alphanumeric passwords
combina -a -n -k 6 -r

### 3-word passphrases from a sentence
combina -phrase "dog cat bird fish" -k 3 -r

### Combinations with MD5 hashes
combina -a -n -s -k 4 -r -md5

### With prefix/suffix and formatted output
combina -phrase "red,green,blue" -input-separator "," -output-separator "-" -k 2 -r -prefix "[" -suffix "]"

### Multi-threading for large volumes
combina -a -A -n -s -k 6 -r -workers 8 -verbose

### Progressive length from 4 to 6 characters
combina -a -n -k 6 -p 4 -r

### 📊 Statistics and Progress

In -verbose mode, the software displays:

    Total number of combinations to generate

    Real-time progress

    Generation speed (combinations/sec)

    Per-worker statistics

    Total elapsed time

### 🔧 System Requirements

    Go 1.21 or higher

    Operating System: Linux, macOS, Windows

    Memory: Variable based on combination size

    CPU: Multi-core support for parallel processing

### 📦 Installation
bash

## Direct installation
go install combina/cmd/combina@latest

## Or build from source
git clone [[github.com/cydside/combina]](https://github.com/cydside/combina)
cd combina
go build -o combina cmd/combina/main.go

### 📚 API Documentation

The software is structured as a modular Go package:

    combinacore: Core generation logic

    formats: Output format handlers

    Easily extensible with new formatters or generators

### 🤝 Contributing

The project welcomes contributions! Potential improvement areas:

    Optimized generation algorithms

    Graphical interfaces

### 📄 License

The MIT License (MIT) - Copyright (c) 2006-present Danilo CICERONE
