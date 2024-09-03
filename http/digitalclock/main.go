//go:build !solution

package main

import (
	"bytes"
	"flag"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	h      = 12
	w      = 8
	wColon = 4
)

var (
	colorMap = map[rune]color.RGBA{
		'.': {255, 255, 255, 255},
		'1': Cyan,
	}
	symbolMap = map[rune]string{
		'0': Zero,
		'1': One,
		'2': Two,
		'3': Three,
		'4': Four,
		'5': Five,
		'6': Six,
		'7': Seven,
		'8': Eight,
		'9': Nine,
		':': Colon,
	}
)

func setSymbol(k int, pattern string, img *image.RGBA, shift int) {
	for y, line := range strings.Split(pattern, "\n") {
		for x, char := range line {
			c := colorMap[char]
			for i := 0; i < (x+1)*k; i++ {
				for j := 0; j < (y+1)*k; j++ {
					img.SetRGBA(shift+i+x*k, j+y*k, c)
				}
			}
		}
	}
}

func generateImage(k int, time string) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, (6*w+2*wColon)*k, h*k))
	shift := 0
	for _, char := range time {
		setSymbol(k, symbolMap[char], img, shift)
		if char == ':' {
			shift += wColon * k
		} else {
			shift += w * k
		}
	}
	return img
}

func clockHandler(w http.ResponseWriter, r *http.Request) {
	kParam := r.URL.Query().Get("k")
	if kParam == "" {
		kParam = "1"
	}
	k, err := strconv.Atoi(kParam)
	if err != nil || k < 1 || k > 30 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("invalid k"))
		return
	}

	timeParam := r.URL.Query().Get("time")
	if timeParam == "" {
		timeParam = time.Now().Format(time.TimeOnly)
	} else {
		_, err := time.Parse(time.TimeOnly, timeParam)
		if err != nil || len(timeParam) != 8 {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte("invalid time"))
			return
		}
	}

	buf := new(bytes.Buffer)
	_ = png.Encode(buf, generateImage(k, timeParam))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(buf.Bytes())
}

func main() {
	port := flag.String("port", "6029", "Specify server port. Default is 6029")
	flag.Parse()

	http.HandleFunc("/", clockHandler)
	_ = http.ListenAndServe(":"+*port, nil)
}
