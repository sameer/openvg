// Package openvg is a high-level 2D vector graphics library built on OpenVG
package openvg

/*
#cgo CFLAGS:   -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads
#cgo LDFLAGS:  -L/opt/vc/lib -lGLESv2 -lEGL -lbcm_host -ljpeg
#include "VG/openvg.h"
#include "VG/vgu.h"
#include "EGL/egl.h"
#include "GLES/gl.h"
#include "fontinfo.h" // font information
#include "shapes.h"   // C API
*/
import "C"
import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"runtime"
	"strings"
	"unsafe"
	"image/color"
)

// VGfloat defines the basic type for coordinates, dimensions and other values
type VGfloat C.VGfloat
type VGint C.VGint

// Offcolor defines the offset, color and alpha values used in gradients
// the Offset ranges from 0..1, colors as RGB triples, alpha ranges from 0..1
type Offcolor struct {
	Offset VGfloat
	color.RGBA
}

func UnwrapRGBA(rgba color.RGBA) (r, g, b uint8, a VGfloat) {
	return rgba.R, rgba.G, rgba.B, VGfloat(rgba.A) / 255.0
}

func UnwrapRGB(rgba color.RGBA) (r, g, b uint8) {
	return rgba.R, rgba.G, rgba.B
}


// colornames maps SVG color names to RGB triples.
var colornames = map[string]color.RGBA{
	"aliceblue":            {240, 248, 255, 255},
	"antiquewhite":         {250, 235, 215, 255},
	"aqua":                 {0, 255, 255, 255},
	"aquamarine":           {127, 255, 212, 255},
	"azure":                {240, 255, 255, 255},
	"beige":                {245, 245, 220, 255},
	"bisque":               {255, 228, 196, 255},
	"black":                {0, 0, 0, 255},
	"blanchedalmond":       {255, 235, 205, 255},
	"blue":                 {0, 0, 255, 255},
	"blueviolet":           {138, 43, 226, 255},
	"brown":                {165, 42, 42, 255},
	"burlywood":            {222, 184, 135, 255},
	"cadetblue":            {95, 158, 160, 255},
	"chartreuse":           {127, 255, 0, 255},
	"chocolate":            {210, 105, 30, 255},
	"coral":                {255, 127, 80, 255},
	"cornflowerblue":       {100, 149, 237, 255},
	"cornsilk":             {255, 248, 220, 255},
	"crimson":              {220, 20, 60, 255},
	"cyan":                 {0, 255, 255, 255},
	"darkblue":             {0, 0, 139, 255},
	"darkcyan":             {0, 139, 139, 255},
	"darkgoldenrod":        {184, 134, 11, 255},
	"darkgray":             {169, 169, 169, 255},
	"darkgreen":            {0, 100, 0, 255},
	"darkgrey":             {169, 169, 169, 255},
	"darkkhaki":            {189, 183, 107, 255},
	"darkmagenta":          {139, 0, 139, 255},
	"darkolivegreen":       {85, 107, 47, 255},
	"darkorange":           {255, 140, 0, 255},
	"darkorchid":           {153, 50, 204, 255},
	"darkred":              {139, 0, 0, 255},
	"darksalmon":           {233, 150, 122, 255},
	"darkseagreen":         {143, 188, 143, 255},
	"darkslateblue":        {72, 61, 139, 255},
	"darkslategray":        {47, 79, 79, 255},
	"darkslategrey":        {47, 79, 79, 255},
	"darkturquoise":        {0, 206, 209, 255},
	"darkviolet":           {148, 0, 211, 255},
	"deeppink":             {255, 20, 147, 255},
	"deepskyblue":          {0, 191, 255, 255},
	"dimgray":              {105, 105, 105, 255},
	"dimgrey":              {105, 105, 105, 255},
	"dodgerblue":           {30, 144, 255, 255},
	"firebrick":            {178, 34, 34, 255},
	"floralwhite":          {255, 250, 240, 255},
	"forestgreen":          {34, 139, 34, 255},
	"fuchsia":              {255, 0, 255, 255},
	"gainsboro":            {220, 220, 220, 255},
	"ghostwhite":           {248, 248, 255, 255},
	"gold":                 {255, 215, 0, 255},
	"goldenrod":            {218, 165, 32, 255},
	"gray":                 {128, 128, 128, 255},
	"green":                {0, 128, 0, 255},
	"greenyellow":          {173, 255, 47, 255},
	"grey":                 {128, 128, 128, 255},
	"honeydew":             {240, 255, 240, 255},
	"hotpink":              {255, 105, 180, 255},
	"indianred":            {205, 92, 92, 255},
	"indigo":               {75, 0, 130, 255},
	"ivory":                {255, 255, 240, 255},
	"khaki":                {240, 230, 140, 255},
	"lavender":             {230, 230, 250, 255},
	"lavenderblush":        {255, 240, 245, 255},
	"lawngreen":            {124, 252, 0, 255},
	"lemonchiffon":         {255, 250, 205, 255},
	"lightblue":            {173, 216, 230, 255},
	"lightcoral":           {240, 128, 128, 255},
	"lightcyan":            {224, 255, 255, 255},
	"lightgoldenrodyellow": {250, 250, 210, 255},
	"lightgray":            {211, 211, 211, 255},
	"lightgreen":           {144, 238, 144, 255},
	"lightgrey":            {211, 211, 211, 255},
	"lightpink":            {255, 182, 193, 255},
	"lightsalmon":          {255, 160, 122, 255},
	"lightseagreen":        {32, 178, 170, 255},
	"lightskyblue":         {135, 206, 250, 255},
	"lightslategray":       {119, 136, 153, 255},
	"lightslategrey":       {119, 136, 153, 255},
	"lightsteelblue":       {176, 196, 222, 255},
	"lightyellow":          {255, 255, 224, 255},
	"lime":                 {0, 255, 0, 255},
	"limegreen":            {50, 205, 50, 255},
	"linen":                {250, 240, 230, 255},
	"magenta":              {255, 0, 255, 255},
	"maroon":               {128, 0, 0, 255},
	"mediumaquamarine":     {102, 205, 170, 255},
	"mediumblue":           {0, 0, 205, 255},
	"mediumorchid":         {186, 85, 211, 255},
	"mediumpurple":         {147, 112, 219, 255},
	"mediumseagreen":       {60, 179, 113, 255},
	"mediumslateblue":      {123, 104, 238, 255},
	"mediumspringgreen":    {0, 250, 154, 255},
	"mediumturquoise":      {72, 209, 204, 255},
	"mediumvioletred":      {199, 21, 133, 255},
	"midnightblue":         {25, 25, 112, 255},
	"mintcream":            {245, 255, 250, 255},
	"mistyrose":            {255, 228, 225, 255},
	"moccasin":             {255, 228, 181, 255},
	"navajowhite":          {255, 222, 173, 255},
	"navy":                 {0, 0, 128, 255},
	"oldlace":              {253, 245, 230, 255},
	"olive":                {128, 128, 0, 255},
	"olivedrab":            {107, 142, 35, 255},
	"orange":               {255, 165, 0, 255},
	"orangered":            {255, 69, 0, 255},
	"orchid":               {218, 112, 214, 255},
	"palegoldenrod":        {238, 232, 170, 255},
	"palegreen":            {152, 251, 152, 255},
	"paleturquoise":        {175, 238, 238, 255},
	"palevioletred":        {219, 112, 147, 255},
	"papayawhip":           {255, 239, 213, 255},
	"peachpuff":            {255, 218, 185, 255},
	"peru":                 {205, 133, 63, 255},
	"pink":                 {255, 192, 203, 255},
	"plum":                 {221, 160, 221, 255},
	"powderblue":           {176, 224, 230, 255},
	"purple":               {128, 0, 128, 255},
	"red":                  {255, 0, 0, 255},
	"rosybrown":            {188, 143, 143, 255},
	"royalblue":            {65, 105, 225, 255},
	"saddlebrown":          {139, 69, 19, 255},
	"salmon":               {250, 128, 114, 255},
	"sandybrown":           {244, 164, 96, 255},
	"seagreen":             {46, 139, 87, 255},
	"seashell":             {255, 245, 238, 255},
	"sienna":               {160, 82, 45, 255},
	"silver":               {192, 192, 192, 255},
	"skyblue":              {135, 206, 235, 255},
	"slateblue":            {106, 90, 205, 255},
	"slategray":            {112, 128, 144, 255},
	"slategrey":            {112, 128, 144, 255},
	"snow":                 {255, 250, 250, 255},
	"springgreen":          {0, 255, 127, 255},
	"steelblue":            {70, 130, 180, 255},
	"tan":                  {210, 180, 140, 255},
	"teal":                 {0, 128, 128, 255},
	"thistle":              {216, 191, 216, 255},
	"tomato":               {255, 99, 71, 255},
	"turquoise":            {64, 224, 208, 255},
	"violet":               {238, 130, 238, 255},
	"wheat":                {245, 222, 179, 255},
	"white":                {255, 255, 255, 255},
	"whitesmoke":           {245, 245, 245, 255},
	"yellow":               {255, 255, 0, 255},
	"yellowgreen":          {154, 205, 50, 255},
}

// Init initializes the graphics subsystem
func Init() (int, int) {
	runtime.LockOSThread()
	var rh, rw C.int
	C.init(&rw, &rh)
	return int(rw), int(rh)
}

// InitWidowSize initialized the graphics subsystem with specified dimensions
func InitWindowSize(x, y, w, h int) {
	C.initWindowSize(C.int(x), C.int(y), C.uint(w), C.uint(h))
}

// WindowClear clears the window to previously set background color
func WindowClear() {
	C.WindowClear()
}

// WindowPostion places a window
func WindowPosition(x, y int) {
	C.WindowPosition(C.int(x), C.int(y))
}

// WindowOpacity sets the window's opacity
func WindowOpacity(a uint) {
	C.WindowOpacity(C.uint(a))
}

// AreaClear clears a given rectangle in window coordinates
func AreaClear(x, y, w, h int) {
	C.AreaClear(C.uint(x), C.uint(y), C.uint(w), C.uint(h))
}

// Finish shuts down the graphics subsystem
func Finish() {
	C.finish()
	runtime.UnlockOSThread()
}

// Background clears the screen with the specified solid background color using RGB triples
func Background(r, g, b uint8) {
	C.Background(C.uint(r), C.uint(g), C.uint(b))
}

// BackgroundRGB clears the screen with the specified background color using a RGBA quad
func BackgroundRGB(r, g, b uint8, alpha VGfloat) {
	C.BackgroundRGB(C.uint(r), C.uint(g), C.uint(b), C.VGfloat(alpha))
}

// BackgroundColor sets the background color
func BackgroundColor(s string, alpha ...VGfloat) {
	c := Colorlookup(s)
	if len(alpha) == 0 {
		BackgroundRGB(c.R, c.G, c.B, 1)
	} else {
		BackgroundRGB(c.R, c.G, c.B, alpha[0])
	}
}

// makestops prepares the color/stop vector
func makeramp(r []Offcolor) (*C.VGfloat, C.int) {
	lr := len(r)
	nr := lr * 5
	cs := make([]C.VGfloat, nr)
	j := 0
	for i := 0; i < lr; i++ {
		cs[j] = C.VGfloat(r[i].Offset)
		j++
		cs[j] = C.VGfloat(VGfloat(r[i].R) / 255.0)
		j++
		cs[j] = C.VGfloat(VGfloat(r[i].G) / 255.0)
		j++
		cs[j] = C.VGfloat(VGfloat(r[i].B) / 255.0)
		j++
		cs[j] = C.VGfloat(VGfloat(r[i].A) / 255.0)
		j++
	}
	return &cs[0], C.int(lr)
}

// FillLinearGradient sets up a linear gradient between (x1,y2) and (x2, y2)
// using the specified offsets and colors in ramp
func FillLinearGradient(x1, y1, x2, y2 VGfloat, ramp []Offcolor) {
	cr, nr := makeramp(ramp)
	C.FillLinearGradient(C.VGfloat(x1), C.VGfloat(y1), C.VGfloat(x2), C.VGfloat(y2), cr, nr)
}

// FillRadialGradient sets up a radial gradient centered at (cx, cy), radius r,
// with a focal point at (fx, fy) using the specified offsets and colors in ramp
func FillRadialGradient(cx, cy, fx, fy, radius VGfloat, ramp []Offcolor) {
	cr, nr := makeramp(ramp)
	C.FillRadialGradient(C.VGfloat(cx), C.VGfloat(cy), C.VGfloat(fx), C.VGfloat(fy), C.VGfloat(radius), cr, nr)
}

// FillRGB sets the fill color, using RGB triples and alpha values
func FillRGB(r, g, b uint8, alpha VGfloat) {
	C.Fill(C.uint(r), C.uint(g), C.uint(b), C.VGfloat(alpha))
}

// StrokeRGB sets the stroke color, using RGB triples
func StrokeRGB(r, g, b uint8, alpha VGfloat) {
	C.Stroke(C.uint(r), C.uint(g), C.uint(b), C.VGfloat(alpha))
}

// StrokeWidth sets the stroke width
func StrokeWidth(w VGfloat) {
	C.StrokeWidth(C.VGfloat(w))
}

// Colorlookup returns a RGB triple corresponding to the named color,
// or "rgb(r,g,b)" string. On error, return black.
func Colorlookup(s string) color.RGBA {
	var rcolor = color.RGBA{0, 0, 0, 255}
	col, ok := colornames[s]
	if ok {
		return col
	}
	if strings.HasPrefix(s, "rgb(") {
		n, err := fmt.Sscanf(s[3:], "(%d,%d,%d,%d)", &rcolor.R, &rcolor.G, &rcolor.B, &rcolor.A)
		if n != 3 || err != nil {
			return color.RGBA{0, 0, 0, 255}
		}
		return rcolor
	}
	return rcolor
}

// FillColor sets the fill color using names to specify the color, optionally applying alpha.
func FillColor(s string, alpha ...VGfloat) {
	fc := Colorlookup(s)
	if len(alpha) == 0 {
		FillRGB(fc.R, fc.G, fc.B, 1)
	} else {
		FillRGB(fc.R, fc.G, fc.B, alpha[0])
	}
}

// StrokeColor sets the fill color using names to specify the color, optionally applying alpha.
func StrokeColor(s string, alpha ...VGfloat) {
	fc := Colorlookup(s)
	if len(alpha) == 0 {
		StrokeRGB(fc.R, fc.G, fc.B, 1)
	} else {
		StrokeRGB(fc.R, fc.G, fc.B, alpha[0])
	}
}

// Start begins a picture
func Start(w, h int, color ...uint8) {
	C.Start(C.int(w), C.int(h))
	if len(color) == 3 {
		Background(color[0], color[1], color[2])
	}
}

// StartColor begins the picture with the specified color background
func StartColor(w, h int, color string, alpha ...VGfloat) {
	C.Start(C.int(w), C.int(h))
	BackgroundColor(color, alpha...)
}

// End ends the picture
func End() {
	C.End()
}

// SaveEnd ends the picture, saving the raw raster
func SaveEnd(filename string) {
	s := C.CString(filename)
	defer C.free(unsafe.Pointer(s))
	C.SaveEnd(s)
}

// fakeimage makes a placeholder for a missing image
func fakeimage(x, y VGfloat, w, h int, s string) {
	fw := VGfloat(w)
	fh := VGfloat(h)
	FillColor("lightgray")
	Rect(x, y, fw, fh)
	StrokeWidth(1)
	StrokeColor("gray")
	Line(x, y, x+fw, y+fh)
	Line(x, y+fh, x+fw, y)
	StrokeWidth(0)
	FillColor("black")
	TextMid(x+(fw/2), y+(fh/2), s, "sans", w/20)
}

// Img places an image object at (x,y)
func Img(x, y VGfloat, im image.Image) {
	bounds := im.Bounds()
	minx := bounds.Min.X
	maxx := bounds.Max.X
	miny := bounds.Min.Y
	maxy := bounds.Max.Y
	data := make([]C.VGubyte, bounds.Dx()*bounds.Dy()*4)
	n := 0
	var r, g, b, a uint32
	for yp := miny; yp < maxy; yp++ {
		for xp := minx; xp < maxx; xp++ {
			r, g, b, a = im.At(xp, (maxy-1)-yp).RGBA() // OpenVG has origin at lower left, y increasing up
			data[n] = C.VGubyte(r >> 8)
			n++
			data[n] = C.VGubyte(g >> 8)
			n++
			data[n] = C.VGubyte(b >> 8)
			n++
			data[n] = C.VGubyte(a >> 8)
			n++
		}
	}
	C.makeimage(C.VGfloat(x), C.VGfloat(y), C.int(bounds.Dx()), C.int(bounds.Dy()), &data[0])
}

// Image places the named image at (x,y) with dimensions (w,h)
// the specified derived image dimensions override the native ones.
func Image(x, y VGfloat, w, h int, s string) {

	var img image.Image
	var derr error
	f, err := os.Open(s)
	if err != nil {
		fakeimage(x, y, w, h, s)
		return
	}
	img, _, derr = image.Decode(f)
	defer f.Close()
	if derr != nil {
		fakeimage(x, y, w, h, s)
		return
	}
	Img(x, y, img)
}

// Line draws a line between two points
func Line(x1, y1, x2, y2 VGfloat) {
	C.Line(C.VGfloat(x1), C.VGfloat(y1), C.VGfloat(x2), C.VGfloat(y2))
}

// Rect draws a rectangle at (x,y) with dimesions (w,h)
func Rect(x, y, w, h VGfloat) {
	C.Rect(C.VGfloat(x), C.VGfloat(y), C.VGfloat(w), C.VGfloat(h))
}

// Roundrect draws a rounded rectangle at (x,y) with dimesions (w,h).
// the corner radii are at (rw, rh)
func Roundrect(x, y, w, h, rw, rh VGfloat) {
	C.Roundrect(C.VGfloat(x), C.VGfloat(y), C.VGfloat(w), C.VGfloat(h), C.VGfloat(rw), C.VGfloat(rh))
}

// Ellipse draws an ellipse at (x,y) with dimensions (w,h)
func Ellipse(x, y, w, h VGfloat) {
	C.Ellipse(C.VGfloat(x), C.VGfloat(y), C.VGfloat(w), C.VGfloat(h))
}

// Circle draws a circle centered at (x,y), with radius r
func Circle(x, y, r VGfloat) {
	C.Circle(C.VGfloat(x), C.VGfloat(y), C.VGfloat(r))
}

// Qbezier draws a quadratic bezier curve with extrema (sx, sy) and (ex, ey)
// Control points are at (cx, cy)
func Qbezier(sx, sy, cx, cy, ex, ey VGfloat) {
	C.Qbezier(C.VGfloat(sx), C.VGfloat(sy), C.VGfloat(cx), C.VGfloat(cy), C.VGfloat(ex), C.VGfloat(ey))
}

// Cbezier draws a cubic bezier curve with extrema (sx, sy) and (ex, ey).
// Control points at (cx, cy) and (px, py)
func Cbezier(sx, sy, cx, cy, px, py, ex, ey VGfloat) {
	C.Cbezier(C.VGfloat(sx), C.VGfloat(sy), C.VGfloat(cx), C.VGfloat(cy), C.VGfloat(px), C.VGfloat(py), C.VGfloat(ex), C.VGfloat(ey))
}

// Arc draws an arc at (x,y) with dimensions (w,h).
// the arc starts at the angle sa, extended to aext
func Arc(x, y, w, h, sa, aext VGfloat) {
	C.Arc(C.VGfloat(x), C.VGfloat(y), C.VGfloat(w), C.VGfloat(h), C.VGfloat(sa), C.VGfloat(aext))
}

// poly converts coordinate slices
func poly(x, y []VGfloat) (*C.VGfloat, *C.VGfloat, C.VGint) {
	size := len(x)
	if size != len(y) {
		return nil, nil, 0
	}
	px := make([]C.VGfloat, size)
	py := make([]C.VGfloat, size)
	for i := 0; i < size; i++ {
		px[i] = C.VGfloat(x[i])
		py[i] = C.VGfloat(y[i])
	}
	return &px[0], &py[0], C.VGint(size)
}

// Polygon draws a polygon with coordinate in x,y
func Polygon(x, y []VGfloat) {
	px, py, np := poly(x, y)
	if np > 0 {
		C.Polygon(px, py, np)
	}
}

// Polyline draws a polyline with coordinates in x, y
func Polyline(x, y []VGfloat) {
	px, py, np := poly(x, y)
	if np > 0 {
		C.Polyline(px, py, np)
	}
}

// selectfont specifies the font by generic name
func selectfont(s string) C.Fontinfo {
	switch s {
	case "sans":
		return C.SansTypeface
	case "serif":
		return C.SerifTypeface
	case "mono":
		return C.MonoTypeface
	case "helvetica":
		return C.HelveticaTypeface
	}
	return C.SerifTypeface
}

// ClipRect limits the drawing area to specified rectangle
func ClipRect(x, y, w, h int) {
	C.ClipRect(C.VGint(x), C.VGint(y), C.VGint(w), C.VGint(h))
}

// ClipEnd stops limiting drawing area to specified rectangle
func ClipEnd() {
	C.ClipEnd()
}

// Text draws text whose aligment begins (x,y)
func Text(x, y VGfloat, s string, font string, size int) {
	t := C.CString(s)
	C.Text(C.VGfloat(x), C.VGfloat(y), t, selectfont(font), C.int(size))
	C.free(unsafe.Pointer(t))
}

// TextMid draws text centered at (x,y)
func TextMid(x, y VGfloat, s string, font string, size int) {
	t := C.CString(s)
	C.TextMid(C.VGfloat(x), C.VGfloat(y), t, selectfont(font), C.int(size))
	C.free(unsafe.Pointer(t))
}

// TextEnd draws text end-aligned at (x,y)
func TextEnd(x, y VGfloat, s string, font string, size int) {
	t := C.CString(s)
	C.TextEnd(C.VGfloat(x), C.VGfloat(y), t, selectfont(font), C.int(size))
	C.free(unsafe.Pointer(t))
}

// TextWidth returns the length of text at a specified font and size
func TextWidth(s string, font string, size int) VGfloat {
	t := C.CString(s)
	defer C.free(unsafe.Pointer(t))
	return VGfloat(C.TextWidth(t, selectfont(font), C.int(size)))
}

// TextHeight returns a font's height (ascent)
func TextHeight(font string, size int) VGfloat {
	return VGfloat(C.TextHeight(selectfont(font), C.int(size)))
}

// TextDepth returns the distance below the baseline for a specified font
func TextDepth(font string, size int) VGfloat {
	return VGfloat(C.TextDepth(selectfont(font), C.int(size)))
}

// Translate translates the coordinate system to (x,y)
func Translate(x, y VGfloat) {
	C.Translate(C.VGfloat(x), C.VGfloat(y))
}

// Rotate rotates the coordinate system around the specifed angle
func Rotate(r VGfloat) {
	C.Rotate(C.VGfloat(r))
}

// Shear warps the coordinate system by (x,y)
func Shear(x, y VGfloat) {
	C.Shear(C.VGfloat(x), C.VGfloat(y))
}

// Scale scales the coordinate system by (x,y)
func Scale(x, y VGfloat) {
	C.Scale(C.VGfloat(x), C.VGfloat(y))
}

// SaveTerm saves terminal settings
func SaveTerm() {
	C.saveterm()
}

// RestoreTerm retores terminal settings
func RestoreTerm() {
	C.restoreterm()
}

// RawTerm sets the terminal to raw mode
func RawTerm() {
	C.rawterm()
}
