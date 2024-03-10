# ngramcounter

A program to count n-grams or bytes of files.

## Introduction

When analyzing classic encryptions it is often necessary to count the frequencies of bytes or [n-grams](https://en.wikipedia.org/wiki/N-gram) of texts.

An n-gram is a sequence of n characters.
A 1-gram is one character.
A 2-gram is two adjacent characters.
And so on.

This utility counts bytes or n-grams in files and writes them into CSV files.
The CSV stands for 'character separated value'

## Call

The program is called like this:

```
ngramcounter [-ngram count] [-encoding encoding] [-separator char] [files...]
```

The parts have the following meaning:

| Part        | Meaning                                                                          |
|-------------|----------------------------------------------------------------------------------|
| `count`     | Number of characters in an n-gram".                                              |
| `encoding`  | Character encoding. Can be any of the list below                                 |
| `separator` | Character used for separating fields in the CSV output. Can be either ',' or ';' |
| `files`     | List of file names count.                                                        |

For every file in the file list a file with the name `filename.csv` is written.
I.e. the file name is appended with `.csv`.

The possible return codes are the following:

| Code | Meaning                   |
|------|---------------------------|
| `0`  | Successful processing     |
| `1`  | Error in the command line |
| `2`  | Error while processing    |

The resulting output file has three columns:

| Name        | Meaning                                                              |
|-------------|----------------------------------------------------------------------|
| Byte/N-gram | Byte or n-gram.                                                      |
| Count       | Number of the times the byte or n-gram if found in the file.         |
| Share       | Share of the byte or n-gram of the total number of bytes or n-grams. |

## Contact

Frank Schwab ([Mail](mailto:github.sfdhi@slmails.com "Mail"))

## License

This source code is published under the [Apache License V2](https://www.apache.org/licenses/LICENSE-2.0.txt).
