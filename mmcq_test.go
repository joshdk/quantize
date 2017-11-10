// Copyright 2017 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package quantize

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinMax(t *testing.T) {

	tests := []struct {
		title string
		min   uint8
		max   uint8
	}{
		{
			title: "both zeros",
			min:   0,
			max:   0,
		},
		{
			title: "both 255",
			min:   255,
			max:   255,
		},
		{
			title: "off by one",
			min:   0,
			max:   1,
		},
		{
			title: "full spread",
			min:   0,
			max:   255,
		},
		{
			title: "half spread",
			min:   127,
			max:   255,
		},
		{
			title: "random values",
			min:   65,
			max:   197,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d - %s", index, test.title)

		t.Run(name, func(t *testing.T) {

			actualMin, actualMinRev := min(test.min, test.max), min(test.max, test.min)

			assert.Equal(t, test.min, actualMin)
			assert.Equal(t, test.min, actualMinRev)

			actualMax, actualMaxRev := max(test.min, test.max), max(test.max, test.min)

			assert.Equal(t, test.max, actualMax)
			assert.Equal(t, test.max, actualMaxRev)

		})
	}

}

func TestAverage(t *testing.T) {

	tests := []struct {
		title   string
		pixels  []color.RGBA
		average color.RGBA
	}{
		{
			title:   "no pixels",
			pixels:  []color.RGBA{},
			average: color.RGBA{0, 0, 0, 0xFF},
		},
		{
			title: "ignore alpha",
			pixels: []color.RGBA{
				{105, 32, 165, 0}, // random values
			},
			average: color.RGBA{105, 32, 165, 0xFF},
		},
		{
			title: "single pixel",
			pixels: []color.RGBA{
				{105, 32, 165, 0xFF}, // random values
			},
			average: color.RGBA{105, 32, 165, 0xFF},
		},
		{
			title: "double pixels",
			pixels: []color.RGBA{
				{105, 32, 165, 0xFF}, // random values
				{105, 32, 165, 0xFF},
			},
			average: color.RGBA{105, 32, 165, 0xFF},
		},
		{
			title: "orthogonal pixels",
			pixels: []color.RGBA{
				{255, 0, 0, 0xFF},
				{0, 255, 0, 0xFF},
				{0, 0, 255, 0xFF},
			},
			average: color.RGBA{85, 85, 85, 0xFF},
		},
		{
			title: "random pixels",
			pixels: []color.RGBA{
				{54, 67, 124, 0xFF}, // all random values
				{45, 186, 21, 0xFF},
				{25, 178, 79, 0xFF},
				{213, 125, 245, 0xFF},
				{251, 125, 26, 0xFF},
			},
			average: color.RGBA{117, 136, 99, 0xFF},
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d - %s", index, test.title)

		t.Run(name, func(t *testing.T) {

			actual := Average(test.pixels)

			assert.Equal(t, test.average, actual)

		})
	}

}

func TestSpread(t *testing.T) {

	tests := []struct {
		title  string
		pixels []color.RGBA
		sr     uint8
		sg     uint8
		sb     uint8
	}{
		{
			title:  "no pixels",
			pixels: []color.RGBA{},
		},
		{
			title: "one pixel",
			pixels: []color.RGBA{
				{105, 32, 165, 0xFF}, // random values
			},
		},
		{
			title: "identical pixels",
			pixels: []color.RGBA{
				{105, 32, 165, 0xFF}, // random values
				{105, 32, 165, 0xFF},
				{105, 32, 165, 0xFF},
			},
		},
		{
			title: "max spread",
			pixels: []color.RGBA{
				{255, 0, 0, 0xFF},
				{0, 255, 0, 0xFF},
				{0, 0, 255, 0xFF},
			},
			sr: 255,
			sg: 255,
			sb: 255,
		},
		{
			title: "independent spread",
			pixels: []color.RGBA{
				{105, 36, 168, 0xFF}, // random values
				{106, 32, 171, 0xFF},
				{107, 34, 165, 0xFF},
			},
			sr: 2,
			sg: 4,
			sb: 6,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d - %s", index, test.title)

		t.Run(name, func(t *testing.T) {

			sr, sg, sb := Spread(test.pixels)

			assert.Equal(t, test.sr, sr)
			assert.Equal(t, test.sg, sg)
			assert.Equal(t, test.sb, sb)

		})
	}

}

func TestPartition(t *testing.T) {

	tests := []struct {
		title  string
		pixels []color.RGBA
		left   []color.RGBA
		right  []color.RGBA
	}{
		{
			title:  "no pixels",
			pixels: []color.RGBA{},
			left:   []color.RGBA{},
			right:  []color.RGBA{},
		},
		{
			title: "one pixel",
			pixels: []color.RGBA{
				{0, 0, 0, 0xFF},
			},
			left: []color.RGBA{},
			right: []color.RGBA{
				{0, 0, 0, 0xFF},
			},
		},
		{
			title: "two pixel",
			pixels: []color.RGBA{
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
			},
			left: []color.RGBA{
				{0, 0, 0, 0xFF},
			},
			right: []color.RGBA{
				{0, 0, 0, 0xFF},
			},
		},
		{
			title: "partition by red",
			pixels: []color.RGBA{
				{21, 0, 0, 0xFF},
				{15, 5, 5, 0xFF},
				{10, 10, 10, 0xFF},
				{5, 15, 15, 0xFF},
				{0, 20, 20, 0xFF},
			},
			left: []color.RGBA{
				{0, 20, 20, 0xFF},
				{5, 15, 15, 0xFF},
			},
			right: []color.RGBA{
				{10, 10, 10, 0xFF},
				{15, 5, 5, 0xFF},
				{21, 0, 0, 0xFF},
			},
		},
		{
			title: "partition by green",
			pixels: []color.RGBA{
				{0, 21, 0, 0xFF},
				{5, 15, 5, 0xFF},
				{10, 10, 10, 0xFF},
				{15, 5, 15, 0xFF},
				{20, 0, 20, 0xFF},
			},
			left: []color.RGBA{
				{20, 0, 20, 0xFF},
				{15, 5, 15, 0xFF},
			},
			right: []color.RGBA{
				{10, 10, 10, 0xFF},
				{5, 15, 5, 0xFF},
				{0, 21, 0, 0xFF},
			},
		},
		{
			title: "partition by blue",
			pixels: []color.RGBA{
				{0, 0, 21, 0xFF},
				{5, 5, 15, 0xFF},
				{10, 10, 10, 0xFF},
				{15, 15, 5, 0xFF},
				{20, 20, 0, 0xFF},
			},
			left: []color.RGBA{
				{20, 20, 0, 0xFF},
				{15, 15, 5, 0xFF},
			},
			right: []color.RGBA{
				{10, 10, 10, 0xFF},
				{5, 5, 15, 0xFF},
				{0, 0, 21, 0xFF},
			},
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d - %s", index, test.title)

		t.Run(name, func(t *testing.T) {

			left, right := Partition(test.pixels)

			assert.Equal(t, len(test.pixels), len(left)+len(right))

			assert.Equal(t, test.left, left)
			assert.Equal(t, test.right, right)

		})
	}

}

func TestPixels(t *testing.T) {

	tests := []struct {
		title   string
		pixels  []color.RGBA
		levels  int
		palette []color.RGBA
	}{
		{
			title:  "0 pixels zero levels",
			pixels: []color.RGBA{},
			levels: 0,
			palette: []color.RGBA{
				{0, 0, 0, 0xFF},
			},
		},
		{
			title:  "0 pixels 1 level",
			pixels: []color.RGBA{},
			levels: 1,
			palette: []color.RGBA{
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
			},
		},
		{
			title:  "0 pixels 3 levels",
			pixels: []color.RGBA{},
			levels: 3,
			palette: []color.RGBA{
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
			},
		},
		{
			title: "1 pixel 0 levels",
			pixels: []color.RGBA{
				{255, 0, 0, 0xFF},
			},
			levels: 0,
			palette: []color.RGBA{
				{255, 0, 0, 0xFF},
			},
		},
		{
			title: "1 pixel 1 level",
			pixels: []color.RGBA{
				{255, 0, 0, 0xFF},
			},
			levels: 1,
			palette: []color.RGBA{
				{0, 0, 0, 0xFF},
				{255, 0, 0, 0xFF},
			},
		},
		{
			title: "1 pixel 3 levels",
			pixels: []color.RGBA{
				{255, 0, 0, 0xFF},
			},
			levels: 3,
			palette: []color.RGBA{
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
				{0, 0, 0, 0xFF},
				{255, 0, 0, 0xFF},
			},
		},
		{
			title: "order level 0",
			pixels: []color.RGBA{
				{8, 0, 0, 0xFF},
				{0, 4, 0, 0xFF},
				{8, 4, 0, 0xFF},
				{0, 0, 6, 0xFF},
			},
			levels: 0,
			palette: []color.RGBA{
				{4, 2, 1, 0xFF},
			},
		},
		{
			title: "order level 1",
			pixels: []color.RGBA{
				{8, 0, 0, 0xFF},
				{0, 4, 0, 0xFF},
				{8, 4, 0, 0xFF},
				{0, 0, 6, 0xFF},
			},
			levels: 1,
			palette: []color.RGBA{
				{0, 2, 3, 0xFF},
				{8, 2, 0, 0xFF},
			},
		},
		{
			title: "order level 2",
			pixels: []color.RGBA{
				{8, 0, 0, 0xFF},
				{0, 4, 0, 0xFF},
				{8, 4, 0, 0xFF},
				{0, 0, 6, 0xFF},
			},
			levels: 2,
			palette: []color.RGBA{
				{0, 4, 0, 0xFF},
				{0, 0, 6, 0xFF},
				{8, 0, 0, 0xFF},
				{8, 4, 0, 0xFF},
			},
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d - %s", index, test.title)

		t.Run(name, func(t *testing.T) {

			palette := Pixels(test.pixels, test.levels)

			assert.Equal(t, int(math.Pow(2, float64(test.levels))), len(palette))

			assert.Equal(t, test.palette, palette)

		})
	}

}

func TestImage(t *testing.T) {

	tests := []struct {
		title   string
		path    string
		levels  int
		palette []color.RGBA
	}{
		{
			title:  "jpg file",
			path:   "plush.jpg",
			levels: 3,
			palette: []color.RGBA{
				{R: 0x13, G: 0x25, B: 0x5c, A: 0xff},
				{R: 0x76, G: 0x5b, B: 0x4b, A: 0xff},
				{R: 0x31, G: 0x52, B: 0x99, A: 0xff},
				{R: 0x7f, G: 0x94, B: 0xb1, A: 0xff},
				{R: 0xb9, G: 0x8c, B: 0x5f, A: 0xff},
				{R: 0xd8, G: 0xcd, B: 0xbe, A: 0xff},
				{R: 0xe5, G: 0xe1, B: 0xd8, A: 0xff},
				{R: 0xf8, G: 0xf3, B: 0xe9, A: 0xff},
			},
		},
		{
			title:  "png file",
			path:   "plush.png",
			levels: 3,
			palette: []color.RGBA{
				{R: 0x14, G: 0x25, B: 0x5d, A: 0xff},
				{R: 0x76, G: 0x5b, B: 0x4b, A: 0xff},
				{R: 0x32, G: 0x52, B: 0x99, A: 0xff},
				{R: 0x7f, G: 0x94, B: 0xb1, A: 0xff},
				{R: 0xb9, G: 0x8c, B: 0x5f, A: 0xff},
				{R: 0xd8, G: 0xcc, B: 0xbe, A: 0xff},
				{R: 0xe3, G: 0xe2, B: 0xd9, A: 0xff},
				{R: 0xf8, G: 0xf2, B: 0xe8, A: 0xff},
			},
		},
		{
			title:  "gif file",
			path:   "plush.gif",
			levels: 3,
			palette: []color.RGBA{
				{R: 0x13, G: 0x26, B: 0x5d, A: 0xff},
				{R: 0x78, G: 0x5a, B: 0x49, A: 0xff},
				{R: 0x31, G: 0x53, B: 0x9b, A: 0xff},
				{R: 0x7f, G: 0x92, B: 0xae, A: 0xff},
				{R: 0xb9, G: 0x8c, B: 0x5e, A: 0xff},
				{R: 0xd9, G: 0xce, B: 0xbe, A: 0xff},
				{R: 0xe2, G: 0xe1, B: 0xd9, A: 0xff},
				{R: 0xf8, G: 0xf2, B: 0xe6, A: 0xff},
			},
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d - %s", index, test.title)

		t.Run(name, func(t *testing.T) {

			file, err := os.Open(path.Join("testdata", test.path))
			require.Nil(t, err)
			defer func() {
				if err := file.Close(); err != nil {
					panic(err.Error())
				}
			}()

			img, _, err := image.Decode(file)
			require.Nil(t, err)

			palette := Image(img, test.levels)

			assert.Equal(t, int(math.Pow(2, float64(test.levels))), len(palette))

			assert.Equal(t, test.palette, palette)

		})
	}

}
