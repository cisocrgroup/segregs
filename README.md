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
`segregs [-padding=p] XML IMG BASE`

Segment the the regions from XML and IMG and write the resulting image
files to BASE\_REGID.png with a padding of p pixels. The resulting json
file is written to BASE\_REGID.json.

## Installation
To install just type `go get github.com/cisocrgroup/segregs`.
