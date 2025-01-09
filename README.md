# ngramcounter

A program to count n-grams or bytes of files.

## Introduction

When analyzing classic encryptions it is often necessary to count the frequencies of bytes or [n-grams](https://en.wikipedia.org/wiki/N-gram) of texts.

An n-gram is a sequence of n characters.
A 1-gram is one character.
A 2-gram is two adjacent characters.
And so on.

This utility counts bytes or n-grams in files and writes the counts into CSV files.
"CSV" stands for 'character separated value'.

Only letters and numbers are counted.
Specifically all characters that are in the [Unicode](https://home.unicode.org/) letter or number [categories](https://en.wikipedia.org/wiki/Unicode_character_property#General_Category).

### Modes

The files can be analyzed in overlapping or sequential mode.
Overlapping mode means that a window slides over the text, which is shifted one character at a time.
Sequential mode means that the n-grams are analyzed one after the other.

E.g., when the text `always` is analyzed for 2-grams the modes would produce the following results:

| Mode          | Results                      |
|---------------|------------------------------|
| `overlapping` | `al`, `lw`, `wa`, `ay`, `ys` |
| `sequential`  | `al`, `wa`, `ys`             |

Overlapping mode handles all file lengths, while sequential mode requires that the number of characters in the file is a multiple of the n-gram size. 

Overlapping mode is the default.

### Encodings

Text files are just a stream of bytes with no inherent meaning.
There has to be a mapping between bytes and characters.
I.e., each character is represented by a mapping into a sequence of bytes.
This is called a [character encoding](https://en.wikipedia.org/wiki/Character_encoding).
This encoding needs to be specified.
The list of valid file encodings is given when the program is started with no arguments or the `help` option..

## Call

The program is called like this:

```
ngramcounter [-size <count>] [-encoding <encoding>] [-separator <char>] [-sequential] [files...]
```

### Options

The options have the following meaning:

| Option       | Meaning                                                                           |
|--------------|-----------------------------------------------------------------------------------|
| `size`       | Number of characters in an n-gram".                                               |
| `encoding`   | Character encoding of the source file. Can be any of the list below.              |
| `separator`  | Character used for separating fields in the CSV output. Can be either `,` or `;`. |
| `sequential` | Read n-grams sequentially.                                                        |
| `files`      | List of file names whose contents are to be counted.                              |
| `help`       | Print usage and exit.                                                             |

For every file in the file list a file with the name `<filebasename>_<ext>.csv` is written.
I.e. the file name is appended changed so that the period of the extension becomes an underscore and is then appended with the `.csv` extension.

The default for `separator` is `;`.
This separator implies that the decimal separator is a `,`.
If `separator` is specified as `,` the decimal separator is `.`.

If `sequential` is not specified, the files are analyzed in overlapping mode.

If no argument is specified a usage message is written.
This usage message contains a list of all supported encodings.

#### Encoding

The `encoding` option can have a lot of values that represent all the encodings that Go supports.
They full list is printed, when the program is called with the `help` option. 
Here are the most important ones:

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

`utf16` is a synonym for `utf16le`.

On Windows systems files are normally `Windows 1252`-encoded.
Windows also uses `UTF-16LE` encoding.
Some files may be `UTF-8`-encoded.

Linux systems normally use `UTF-8`.

There is no way to know what the encoding of a file is.
It has to be specified by the user.

However, there are a few exceptions to this rule.
The encoding is known if a file begins with a "[byte-order mark](https://en.wikipedia.org/wiki/Byte_order_mark)".
There exist three known byte-order marks, namely for `UTF-8`, `UTF-16BE` and `UTF-16LE`.
A byte-order mark is not mandatory for files encoded in one of those encodings.
It may or may not be present.
If it is present, the encoding is known.
The program uses the encoding of the byte-order mark, if the file begins with one.

### Output

The resulting output file has three columns:

| Name        | Meaning                                                              |
|-------------|----------------------------------------------------------------------|
| Byte/N-gram | Byte or n-gram.                                                      |
| Count       | Number of the times the byte or n-gram if found in the file.         |
| Share       | Share of the byte or n-gram of the total number of bytes or n-grams. |

### Return codes

The possible return codes are the following:

| Code | Meaning                   |
|------|---------------------------|
| `0`  | Successful processing     |
| `1`  | Error in the command line |
| `2`  | Error while processing    |

## Contact

Frank Schwab ([Mail](mailto:github.sfdhi@slmails.com "Mail"))

## License

This source code is published under the [Apache License V2](https://www.apache.org/licenses/LICENSE-2.0.txt).
