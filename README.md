```
   _____  ___    ____   ___   ___    ____    _____
  / ___/ / _ \  / _  \ /  _\ / _ \  / _  \  / ___/
 (__  ) /  __/ / /_/ //  /  /  __/ / /_/ / (__  )
/____/  \____\ \__  / \_/   \____\ \_   / /____/
              ___/ /              ___/ /
              \___/               \___/
```

# segregs
Segment regions (from
[PageXML](http://www.primaresearch.org/publications/ICPR2010_Pletschacher_PAGE)-formatted
files) into seperate image files.

## Usage
```
segregs [Options] XML IMG OUT
Options:
  -lines segment line regions
  -padding int set padding for region image snippets
  -workers int set number of worker threads (default #cpus)
```

Segment the the regions from XML and IMG and write the resulting image
files to `OUT/REGID.png`. The resulting json file is written to
`OUT_REGID.json` and the text file is written to `OUT/REGID.gt.txt`.
If `-lines` is given, each image line snippet is written to
`OUT/LINENUM.png` and its text is written to `OUT/LINENUM.gt.txt`
instead.  The `-padding` option sets the number of padding pixels
around the written image snippets, the `-workers` sets the number of
concurrent workers.

## Installation
To install just type `go get github.com/cisocrgroup/segregs`.
