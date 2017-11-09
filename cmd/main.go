// Copyright 2017 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"

	"github.com/joshdk/quantize"
)

func die(err error) {
	fmt.Printf("quantize: %s\n", err.Error())
	os.Exit(1)
}

func render(clr color.RGBA) {
	fmt.Printf("#%02X%02X%02X\n", clr.R, clr.G, clr.B)
}

func main() {

	if len(os.Args) < 2 {
		die(errors.New("image file not specified"))
	}

	path := os.Args[1]

	levels := 4

	if len(os.Args) >= 3 {
		var err error
		levels, err = strconv.Atoi(os.Args[2])
		if err != nil {
			die(err)
		}
	}

	file, err := os.Open(path)
	if err != nil {
		die(err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		die(err)
	}

	colors := quantize.Image(img, levels)

	for _, clr := range colors {
		render(clr)
	}

}
