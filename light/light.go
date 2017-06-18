// package light provides functions for converting between properties of
// physical light and perceived colors
// Some interesting links
// WolframAlpha color convert: http://www.wolframalpha.com/widgets/view.jsp?id=628bfed9ce559d754c3eabcfca44366b
// http://www.midnightkite.com/color.html
// http://lsrtools.1apps.com/wavetorgb/index.asp?wavelength=0
// http://stackoverflow.com/questions/3407942/rgb-values-of-visible-spectrum
// http://stackoverflow.com/questions/1472514/convert-light-frequency-to-rgb
// http://pubs.acs.org/doi/pdf/10.1021/acs.jchemed.5b00844
// http://codepen.io/pen/?&editors=011
// https://academo.org/demos/wavelength-to-colour-relationship/
// https://en.wikipedia.org/wiki/CIE_1931_color_space#Meaning_of_X.2C_Y.2C_and_Z
package light

import (
	"image/color"
)

var (
	MinWavelength = 380.0
	MaxWavelength = 780.0
)

// WavelightToColor converts a wavelength in nm from the range
// 380nm to 780nm inclusive to an RGB color.
// The algorthm is taken from this FORTRAN code
// http://www.physics.sfasu.edu/astro/color/spectra.html
func WavelengthToColor(wavelength float64) color.Color {
	// Return a 100% transparent color for wavelengths that are out of bounds
	// if wavelenth < 380.0 || wavelenth > 780.0 {
	//     return color.Alpha{255}
	// }
	c := Color{}
	switch {
	// Return a 100% transparent color for wavelengths that are out of bounds
	case wavelength < MinWavelength || wavelength > MaxWavelength:
		return color.Transparent
	case wavelength >= MinWavelength && wavelength < 490.0:
		c.G = (wavelength - 440.0) / (490.0 - 440.0)
		return c
	case wavelength >= 490.0 && wavelength < 510.0:
		c.G = 1
		c.B = -1 * ((wavelength - 510.0) / (510.0 - 490.0))
		return c
	case wavelength >= 510.0 && wavelength < 580.0:
		c.R = (wavelength - 510.0) / (580.0 - 510.0)
		c.G = 1
		return c
	case wavelength >= 580.0 && wavelength < 645.0:
		c.R = 1
		c.G = -1 * (wavelength - 645.0) / (645.0 - 580.0)
		return c
	case wavelength >= 645.0 && wavelength < MaxWavelength:
		c.R = 1
		return c
	}

	var sss float64

	switch {
	case wavelength > 700:
		sss = .3 + .7*(MaxWavelength-wavelength)/(MaxWavelength-700.0)
	case wavelength < 420.0:
		sss = .3 + .7*(wavelength-380.0)/(420.0-380.0)
	default:
		sss = 1
	}

	_ = sss

	return color.Transparent
}

//r = uint32(col.R*65535.0)
//g = uint32(col.G*65535.0)
//b = uint32(col.B*65535.0)

// Color is stored internally using sRGB (standard RGB) values in the range 0-1
type Color struct {
	R, G, B float64
}

// Implement the Go color.Color interface.
func (col Color) RGBA() (r, g, b, a uint32) {
	r = uint32(col.R * 65535.0)
	g = uint32(col.G * 65535.0)
	b = uint32(col.B * 65535.0)
	a = 0xFFFF
	return
}
