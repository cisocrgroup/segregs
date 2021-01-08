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
  -lines segment line regions with according .gt.txt files
  -padding int set padding for region images
  -workers int set number of worker threads (default #cpus)
```

Segment the the regions from XML and IMG and write the resulting image
files to `OUT/REGID.png` with a padding of p pixels. The resulting
json file is written to `OUT_REGID.json`.  If `-lines` is given, each
line snippet is written to `OUT/LINENUM.png` and its text is written
to `OUT/LINENUM.gt.txt` instead.

## Installation
To install just type `go get github.com/cisocrgroup/segregs`.
