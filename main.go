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
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/antchfx/xmlquery"
	"github.com/cisocrgroup/segregs/poly"
	_ "github.com/hhrutter/tiff"
)

func usage(prog string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-padding=n] XML IMG OUT-BASE\nOptions:\n", prog)
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	padding := flag.Int("padding", 0, "set padding for region images")
	flag.Usage = usage(os.Args[0])
	flag.Parse()
	if len(flag.Args()) != 3 {
		usage(os.Args[0])()
	}
	run(flag.Args()[0], flag.Args()[1], flag.Args()[2], *padding)
}

func run(xmlName, imgName, outBase string, padding int) {
	// Read the iamge once.
	in, err := os.Open(imgName)
	chk(err)
	defer in.Close()
	img, _, err := image.Decode(in)
	chk(err)
	var wg sync.WaitGroup
	rs := regions(xmlName)
	wg.Add(len(rs))
	for _, r := range rs {
		go func(r region) {
			defer wg.Done()
			r.write(img, outBase, padding)
		}(r)
	}
	wg.Wait()
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
		// Read the region's polygon and inner text.
		ps := xmlquery.Find(r, "//*[local-name()='Point']")
		polygon, err := newPolygonFromPoints(ps)
		chk(err)
		textnode := xmlquery.FindOne(r, "//*[local-name()='Unicode']")
		if textnode == nil { // Skip regions with missing Unicode node.
			continue
		}
		reg := region{
			"Coordinates": polygon,
			"Text":        textnode.InnerText(),
		}
		for _, attr := range r.Attr {
			reg[attr.Name.Local] = attr.Value
		}
		ret = append(ret, reg)
	}
	return ret
}

func newPolygonFromPoints(points []*xmlquery.Node) (poly.Polygon, error) {
	attrAsInt := func(node *xmlquery.Node, key string) (int, bool) {
		for _, attr := range node.Attr {
			if attr.Name.Local != key {
				continue
			}
			val, err := strconv.Atoi(attr.Value)
			if err != nil {
				return 0, false
			}
			return val, true
		}
		return 0, false
	}
	var ret poly.Polygon
	for _, point := range points {
		x, xok := attrAsInt(point, "x")
		y, yok := attrAsInt(point, "y")
		if xok && yok {
			ret = append(ret, image.Point{X: x, Y: y})
		}
	}
	return ret, nil
}

func chk(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

type region map[string]interface{}

func (r region) write(img image.Image, outBase string, padding int) {
	// Copy the subregion from the base image.
	coords := r["Coordinates"].(poly.Polygon)
	rect := coords.BoundingRectangle()
	newRect := addPadding(rect, img.Bounds().Max, padding)
	newImg := image.NewRGBA(newRect)
	draw.Draw(newImg, newImg.Bounds(), img, newRect.Min, draw.Src)

	// Mask off pixels outside of the polygon.  Since newImg
	// retains the bounds of the original sub image, we do not
	// need to adjust for the new x- and y-coordinates.
	for x := newImg.Bounds().Min.X; x < newImg.Bounds().Max.X; x++ {
		for y := newImg.Bounds().Min.Y; y < newImg.Bounds().Max.Y; y++ {
			if !coords.Inside(image.Pt(x, y)) {
				newImg.Set(x, y, color.White)
			}
		}
	}

	// Write region png, json and gt.txt files.
	r["Dir"] = fmt.Sprintf("%s_%s", outBase, r["id"].(string))
	r["Image"] = r["Dir"].(string) + ".png"
	pout, err := os.Create(r["Image"].(string))
	chk(err)
	defer func() { chk(pout.Close()) }()
	chk(png.Encode(pout, newImg))
	jout, err := os.Create(r["Dir"].(string) + ".json")
	chk(err)
	defer func() { chk(jout.Close()) }()
	chk(json.NewEncoder(jout).Encode(r))
	gtout := r["Dir"].(string) + ".gt.txt"
	chk(ioutil.WriteFile(gtout, []byte(r["Text"].(string)+"\n"), 0666))
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
