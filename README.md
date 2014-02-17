chicken
=========

A set of tools for writing programming languages in go. Current work is on a PEG parser which will be eventually used to generate itself.

The PEG Parser:
======

NOTE: The current implementation is missing several necessary features to be able to support actual languages.

### Install:
    go get github.com/Logiraptor/chicken


### Defining your language:
The language is expressed as a parsing expression grammar. Rules are expressed like so:

    rule <- partA partB
    partA <- 'a'
    partB <- ~'\\d+'

partA above is a string literal.  
partB above is defined to recognize a regular expression denoted with a `~` before the quoted regexp.

### Planned:
The following have yet to be implemented.

	rule <- partA*
	rule <- partA+
	rule <- partA?
	rule <- partA / partB