package main // import "github.com/cisocrgroup/segregs"

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/antchfx/xmlquery"
	"github.com/cisocrgroup/segregs/poly"
	_ "github.com/hhrutter/tiff"
)

var args = struct {
	padding int
}{}

func usage(prog string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-padding=n] XML IMG OUT-BASE\nOptions:\n", prog)
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	flag.IntVar(&args.padding, "padding", 0, "set padding for region images")
	flag.Usage = usage(os.Args[0])
	flag.Parse()
	if len(flag.Args()) != 3 {
		usage(os.Args[0])()
	}
	run(flag.Args()[0], flag.Args()[1], flag.Args()[2], args.padding)
}

func run(xmlName, imgName, outBase string, padding int) {
	// read image
	in, err := os.Open(imgName)
	chk(err)
	defer in.Close()
	img, _, err := image.Decode(in)
	chk(err)
	// img := readImg(in)
	for _, r := range regions(xmlName) {
		r.write(img, outBase, padding)
	}
}

func regions(name string) []region {
	in, err := os.Open(name)
	chk(err)
	defer in.Close()
	xml, err := xmlquery.Parse(in)
	chk(err)
	rs := xmlquery.Find(xml, "//*[local-name()='TextRegion']")
	var ret []region
	for _, r := range rs {
		// read region polygon
		ps := xmlquery.Find(r, "//*[local-name()='Point']")
		polygon, err := poly.NewFromPoints(ps)
		chk(err)
		// read id, type, language
		id, err := attr(r, "id")
		chk(err)
		typ, _ := attr(r, "type")
		lang, _ := attr(r, "primaryLanguage")
		// read unicode
		textnode := xmlquery.FindOne(r, "//*[local-name()='Unicode']")
		// append new region
		ret = append(ret, region{
			Polygon:         polygon,
			ID:              id,
			Type:            typ,
			PrimaryLanguage: lang,
			Text:            textnode.InnerText(),
		})
	}
	return ret
}

func attr(node *xmlquery.Node, key string) (string, error) {
	for _, attr := range node.Attr {
		if attr.Name.Local == key {
			return attr.Value, nil
		}
	}
	return "", fmt.Errorf("node %s: no attribute: %s", node.Data, key)
}

func chk(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

type region struct {
	Polygon                                     poly.Polygon
	ID, Type, PrimaryLanguage, Text, Image, Dir string
}

func (r region) write(img image.Image, outBase string, padding int) {
	// copy sub image
	rect := r.Polygon.BoundingRectangle()
	newRect := addPadding(rect, img.Bounds().Max, padding)
	newImg := image.NewRGBA(newRect)
	draw.Draw(newImg, newImg.Bounds(), img, newRect.Min, draw.Src)

	// Mask off pixels outside of the polygon.  Since newImg
	// retains the bounds of the original sub image, we do not
	// need to adjust for the new x- and y-coordinates.
	for x := newImg.Bounds().Min.X; x < newImg.Bounds().Max.X; x++ {
		for y := newImg.Bounds().Min.Y; y < newImg.Bounds().Max.Y; y++ {
			if !r.Polygon.Inside(image.Pt(x, y)) {
				newImg.Set(x, y, color.White)
			}
		}
	}

	// write region png and json files
	r.Dir = fmt.Sprintf("%s_%s", outBase, r.ID)
	r.Image = r.Dir + ".png"
	pout, err := os.Create(r.Image)
	chk(err)
	defer func() { chk(pout.Close()) }()
	// chk(png.Encode(pout, subImg))
	chk(png.Encode(pout, newImg))
	jout, err := os.Create(r.Dir + ".json")
	chk(err)
	defer func() { chk(jout.Close()) }()
	chk(json.NewEncoder(jout).Encode(r))
}

func addPadding(rect image.Rectangle, max image.Point, padding int) image.Rectangle {
	minCap := func(a int) int {
		if a < 0 {
			return 0
		}
		return a
	}
	maxCap := func(a, b int) int {
		if a > b {
			return b
		}
		return a
	}
	rect.Min.X = minCap(rect.Min.X - padding)
	rect.Min.Y = minCap(rect.Min.Y - padding)
	rect.Max.X = maxCap(rect.Max.X+padding, max.X)
	rect.Max.Y = maxCap(rect.Max.Y+padding, max.Y)
	return rect
}
