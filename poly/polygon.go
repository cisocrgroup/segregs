package poly // import "github.com/cisocrgroup/segregs/poly"

import (
	"fmt"
	"image"
	"strconv"
	"strings"
)

// Polygon defines a list of closed points that define a polygon.
type Polygon []image.Point

// New creates a new polygon from a whitespace separated list of comma
// separated coordinates e.g. `1,2 3,4 5,6`.  A polygon contains at
// least 3 points.  If 2 points are given it is assumed, that the
// region is a rectangle with an upper and lower point.  In this case
// a polygon with all the 4 points is returned.
func New(coordinates string) (Polygon, error) {
	var ret Polygon
	points := strings.Split(coordinates, " ")
	for _, point := range points {
		xy := strings.Split(point, ",")
		if len(xy) != 2 {
			return nil, fmt.Errorf("invalid coordinates for polygon: %s", coordinates)
		}
		x, e1 := strconv.Atoi(xy[0])
		y, e2 := strconv.Atoi(xy[1])
		if e1 != nil || e2 != nil {
			return nil, fmt.Errorf("invalid coordinates for polygon: %s", coordinates)
		}
		ret = append(ret, image.Point{X: x, Y: y})
	}
	if len(ret) < 2 {
		return nil, fmt.Errorf("invalid coordinates for polygon: %s", coordinates)
	}
	if len(ret) == 2 {
		return Polygon{
			ret[0],
			image.Pt(ret[1].X, ret[0].Y),
			ret[1],
			image.Pt(ret[0].X, ret[1].Y),
		}, nil
	}
	return ret, nil
}

const (
	maxuint = ^uint(0)
	maxint  = int(maxuint >> 1)
	minint  = -maxint - 1
)

// BoundingRectangle returns the minimal rectangle containing all of
// the polygon's points.
func (p Polygon) BoundingRectangle() image.Rectangle {
	var (
		min = image.Point{X: maxint, Y: maxint}
		max = image.Point{X: minint, Y: minint}
	)
	for _, point := range p {
		if point.X < min.X {
			min.X = point.X
		}
		if point.Y < min.Y {
			min.Y = point.Y
		}
		if point.X > max.X {
			max.X = point.X
		}
		if point.Y > max.Y {
			max.Y = point.Y
		}
	}
	return image.Rectangle{Min: min, Max: max}
}

// Inside returns true if the given point lies within the polygon.
// Implementation: https://stackoverflow.com/questions/217578/how-can-i-determine-whether-a-2d-point-is-within-a-polygon
func (p Polygon) Inside(point image.Point) bool {
	if len(p) == 0 {
		return false
	}
	rect := p.BoundingRectangle()
	min, max := rect.Min, rect.Max
	if point.X < min.X || point.X > max.X || point.Y < min.Y || point.Y > max.Y {
		return false
	}
	inside := false
	j := len(p) - 1
	for i := 0; i < len(p); i++ {
		if (p[i].Y > point.Y) != (p[j].Y > point.Y) && point.X < (p[j].X-p[i].X)*(point.Y-p[i].Y)/(p[j].Y-p[i].Y)+p[i].X {
			inside = !inside
		}
		j = i
	}
	// log.Printf("CASE 3: %t", inside)
	return inside
}

func (p Polygon) String() string {
	points := make([]string, len(p))
	for i, point := range p {
		points[i] = fmt.Sprintf("%s", point)
	}
	return strings.Join(points, "-")
}
