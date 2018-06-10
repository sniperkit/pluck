
# Common options

Many of the command line programs that are part of _datatools_ use command line
options.  These have largely been standardized. Here's a list of common options
and their general roles.

## Standard Options

options | description
:--- |:---
-h, -help | display help
-v, -version | display version
-l, -license | display license
-i, -input | read from this input file, to read from standard input use `-i -`
-o, -output | write to this output file
-quiet | suppress error messages

## Common Options

options | description
:--- |:---
-nl,-newline | if true add a trailing newline, if false suppress it
-d, -delimiter | for those commands that support setting delimiters
-od, -output-delimiter | for those commands that support output delimiters
-s, -start | for start times and indexe values
-e, -end | for end times and indexe values
-E, -expression | for expressions used by programs, e.g. filter expressions



