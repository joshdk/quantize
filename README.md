[![License](https://img.shields.io/github/license/joshdk/quantize.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/joshdk/quantize?status.svg)](https://godoc.org/github.com/joshdk/quantize)
[![Go Report Card](https://goreportcard.com/badge/github.com/joshdk/quantize)](https://goreportcard.com/report/github.com/joshdk/quantize)
[![CircleCI](https://circleci.com/gh/joshdk/quantize.svg?&style=shield)](https://circleci.com/gh/joshdk/quantize/tree/master)
[![CodeCov](https://codecov.io/gh/joshdk/quantize/branch/master/graph/badge.svg)](https://codecov.io/gh/joshdk/quantize)

# Quantize

🎨 Simple color palette quantization using MMCQ

## Installing

You can fetch this library by running the following

    go get -u github.com/joshdk/quantize

## Usage

```go
import (
	"image/color"
	"image/jpeg"
	"net/http"
	"github.com/joshdk/preview"
	"github.com/joshdk/quantize"
)

resp, err := http.Get("https://i.imgur.com/X9GB4Pu.jpg")
if err != nil {
	panic(err.Error())
}

img, err := jpeg.Decode(resp.Body)
if err != nil {
	panic(err.Error())
}

colors := quantize.Image(img, 4)

palette := make([]color.Color, len(colors))
for index, clr := range colors {
	palette[index] = clr
}

preview.Show(palette)
```

## License

This library is distributed under the [MIT License](https://opensource.org/licenses/MIT), see LICENSE.txt for more information.