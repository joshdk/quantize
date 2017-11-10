// Copyright 2017 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package quantize

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
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
