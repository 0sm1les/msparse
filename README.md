# msparse

msparse is a small command line utility that parses masscan ouput into a ip:port list.

## usage
```bash
# msparse
usage:
        msparse <input type> <input file> <output file>
input types:
        xml, json, list
example:
        msparse list masscan.txt filteredscan.txt
```