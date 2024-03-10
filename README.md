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

For n-gram counting the [character encoding](https://en.wikipedia.org/wiki/Character_encoding) of the source file must be known.
It can any of the following values:

| Name        | Meaning                                                          |
|-------------|------------------------------------------------------------------|
| `cp437`     | [IBM Code Page 437](https://en.wikipedia.org/wiki/Code_page_437) |
| `cp850`     | [IBM Code Page 850](https://en.wikipedia.org/wiki/Code_page_850) |
| `cp852`     | [IBM Code Page 852](https://en.wikipedia.org/wiki/Code_page_852) |
| `iso88591`  | [ISO 8859-1](https://en.wikipedia.org/wiki/ISO/IEC_8859-1)       |
| `iso885915` | [ISO 8859-15](https://en.wikipedia.org/wiki/ISO/IEC_8859-15)     |
| `utf16be`   | [UTF-16BE](https://en.wikipedia.org/wiki/UTF-16)                 |
| `utf16le`   | [UTF-16LE](https://en.wikipedia.org/wiki/UTF-16)                 |
| `utf8`      | [UTF-8](https://en.wikipedia.org/wiki/UTF-8)                     |
| `win1250`   | [Windows 1250](https://en.wikipedia.org/wiki/Windows-1250)       |
| `win1252`   | [Windows 1252](https://en.wikipedia.org/wiki/Windows-1252)       |

`utf16` can be used as a synonym for `utf16be`.

On Windows systems files a normally `Windows 1252` encoded.
Linux systems normally use `UTF-8`.

## Contact

Frank Schwab ([Mail](mailto:github.sfdhi@slmails.com "Mail"))

## License

This source code is published under the [Apache License V2](https://www.apache.org/licenses/LICENSE-2.0.txt).
