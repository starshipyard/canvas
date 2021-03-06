package canvas

import (
	"math"
	"strings"
	"testing"

	"github.com/tdewolff/test"
)

func TestPathBounds(t *testing.T) {
	var tts = []struct {
		orig   string
		bounds Rect
	}{
		{"Q50 100 100 0z", Rect{0, 0, 100, 50}},
		{"Q0 0 100 0z", Rect{0, 0, 100, 0}},
		{"Q100 0 100 0z", Rect{0, 0, 100, 0}},
		{"C0 100 100 100 100 0z", Rect{0, 0, 100, 75}},
		{"C0 0 100 90 100 0z", Rect{0, 0, 100, 40}},
		{"C0 90 100 0 100 0z", Rect{0, 0, 100, 40}},
		{"C0 0 100 0 100 0z", Rect{0, 0, 100, 0}},
	}
	for _, tt := range tts {
		t.Run(tt.orig, func(t *testing.T) {
			p, _ := ParseSVGPath(tt.orig)
			bounds := p.Bounds()
			test.Float(t, bounds.X, tt.bounds.X)
			test.Float(t, bounds.Y, tt.bounds.Y)
			test.Float(t, bounds.W, tt.bounds.W)
			test.Float(t, bounds.H, tt.bounds.H)
		})
	}
}

func TestPathLength(t *testing.T) {
	var tts = []struct {
		orig   string
		length float64
	}{
		{"C0 66.67 100 66.67 100 0", 159.00}, // TODO: is value correct?
		{"A30 60 0 0 0 60 0", 145.33},
	}
	for _, tt := range tts {
		t.Run(tt.orig, func(t *testing.T) {
			p, _ := ParseSVGPath(tt.orig)
			length := p.Length()
			if math.Abs(length-tt.length) > 0.1 {
				test.Fail(t, length, "!=", tt.length)
			}
		})
	}
}

//func TestCubicBezierLength(t *testing.T) {
//	p0 := Point{0.0, 0.0}
//	p1 := Point{0.0, 66.67}
//	p2 := Point{100.0, 66.67}
//	p3 := Point{100.0, 0.0}
//	S := cubicBezierLengthFunc(p0, p1, p2, p3)
//	T := cubicBezierInverseLengthFunc(p0, p1, p2, p3)
//	for _, x := range []float64{0.0, 0.5, 1.0} {
//		s := S(x)
//		xx := T(s)
//		if math.Abs(x-xx) > 0.1 {
//			test.Fail(t, x, "!=", xx)
//		}
//	}
//}
//
//func TestEllipseLength(t *testing.T) {
//}

func TestPathSplit(t *testing.T) {
	var tts = []struct {
		orig  string
		split []string
	}{
		{"M5 5L6 6z", []string{"M5 5L6 6z"}},
		{"L5 5M10 10L20 20z", []string{"M0 0L5 5", "M10 10L20 20z"}},
		{"L5 5zL10 10", []string{"M0 0L5 5z", "M0 0L10 10"}},
		{"M5 5M10 10z", []string{"M5 5", "M10 10z"}},
	}
	for _, tt := range tts {
		t.Run(tt.orig, func(t *testing.T) {
			p, _ := ParseSVGPath(tt.orig)
			ps := p.Split()
			if len(ps) != len(tt.split) {
				origs := []string{}
				for _, p := range ps {
					origs = append(origs, p.String())
				}
				test.T(t, strings.Join(origs, "\n"), strings.Join(tt.split, "\n"))
			} else {
				for i, p := range ps {
					test.T(t, p.String(), tt.split[i])
				}
			}
		})
	}
}

func TestPathSplitAt(t *testing.T) {
	var tts = []struct {
		orig  string
		d     []float64
		split []string
	}{
		{"L4 3L8 0z", []float64{0.0, 5.0, 10.0, 18.0}, []string{"M0 0L4 3", "M4 3L8 0", "M8 0H0"}},
		{"L4 3L8 0z", []float64{5.0, 20.0}, []string{"M0 0L4 3", "M4 3L8 0H0"}},
		{"L4 3L8 0z", []float64{2.5, 7.5, 14.0}, []string{"M0 0L2 1.5", "M2 1.5L4 3L6 1.5", "M6 1.5L8 0H4", "M4 0H0"}},
		{"C0 0 10 0 10 0", []float64{5.0}, []string{"M0 0C0 0 2.5 0 5 0", "M5 0C7.5 0 10 0 10 0"}},
	}
	for _, tt := range tts {
		t.Run(tt.orig, func(t *testing.T) {
			p, _ := ParseSVGPath(tt.orig)
			ps := p.SplitAt(tt.d...)
			if len(ps) != len(tt.split) {
				origs := []string{}
				for _, p := range ps {
					origs = append(origs, p.String())
				}
				test.T(t, strings.Join(origs, "\n"), strings.Join(tt.split, "\n"))
			} else {
				for i, p := range ps {
					test.T(t, p.String(), tt.split[i])
				}
			}
		})
	}
}

func TestPathTranslate(t *testing.T) {
	var tts = []struct {
		dx, dy     float64
		orig       string
		translated string
	}{
		{10.0, 10.0, "M5 5L10 0Q10 10 15 5C15 0 20 0 20 5A5 5 0 0 0 30 5z", "M15 15L20 10Q20 20 25 15C25 10 30 10 30 15A5 5 0 0 0 40 15z"},
	}
	for _, tt := range tts {
		t.Run(tt.orig, func(t *testing.T) {
			p, _ := ParseSVGPath(tt.orig)
			p = p.Translate(tt.dx, tt.dy)
			test.T(t, p.String(), tt.translated)
		})
	}
}

func TestPathReverse(t *testing.T) {
	var tts = []struct {
		orig string
		inv  string
	}{
		{"M5 5", "M5 5"},
		{"M5 5z", "M5 5z"},
		{"M5 5V10L10 5", "M10 5L5 10V5"},
		{"M5 5V10L10 5z", "M5 5H10L5 10z"},
		{"M5 5V10L10 5M10 10V20L20 10z", "M10 10H20L10 20zM10 5L5 10V5"},
		{"M5 5V10L10 5zM10 10V20L20 10z", "M10 10H20L10 20zM5 5H10L5 10z"},
		{"M5 5Q10 10 15 5", "M15 5Q10 10 5 5"},
		{"M5 5Q10 10 15 5z", "M5 5H15Q10 10 5 5z"},
		{"M5 5C5 10 10 10 10 5", "M10 5C10 10 5 10 5 5"},
		{"M5 5C5 10 10 10 10 5z", "M5 5H10C10 10 5 10 5 5z"},
		{"M5 5A2.5 5 0 0 0 10 5", "M10 5A2.5 5 0 0 1 5 5"}, // bottom-half of ellipse along y
		{"M5 5A2.5 5 0 0 1 10 5", "M10 5A2.5 5 0 0 0 5 5"},
		{"M5 5A2.5 5 0 1 0 10 5", "M10 5A2.5 5 0 1 1 5 5"},
		{"M5 5A2.5 5 0 1 1 10 5", "M10 5A2.5 5 0 1 0 5 5"},
		{"M5 5A5 2.5 90 0 0 10 5", "M10 5A5 2.5 90 0 1 5 5"}, // same shape
		{"M5 5A2.5 5 0 0 0 10 5z", "M5 5H10A2.5 5 0 0 1 5 5z"},
		{"L0 5L5 5", "M5 5H0V0"},
		{"L-1 5L5 5z", "M0 0L5 5H-1z"},
		{"Q0 5 5 5", "M5 5Q0 5 0 0"},
		{"Q0 5 5 5z", "M0 0L5 5Q0 5 0 0z"},
		{"C0 5 5 5 5 0", "M5 0C5 5 0 5 0 0"},
		{"C0 5 5 5 5 0z", "M0 0H5C5 5 0 5 0 0z"},
		{"A2.5 5 0 0 0 5 0", "M5 0A2.5 5 0 0 1 0 0"},
		{"A2.5 5 0 0 0 5 0z", "M0 0H5A2.5 5 0 0 1 0 0z"},
		{"M5 5L10 10zL15 10", "M15 10L5 5L10 10z"},
	}
	for _, tt := range tts {
		t.Run(tt.orig, func(t *testing.T) {
			p, _ := ParseSVGPath(tt.orig)
			p = p.Reverse()
			test.T(t, p.String(), tt.inv)
		})
	}
}
