package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/muesli/termenv"
	"github.com/nfnt/resize"
)

const asciiChars = "@#Ss:. "

func main() {
	// Get resolution from user (default: width = 100)
	width := 50
	if len(os.Args) > 1 {
		w, err := strconv.Atoi(os.Args[1])
		if err == nil && w > 0 {
			width = w
		}
	}

	// Read file
	data, err := ioutil.ReadFile("input.png")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Decode image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// Set height based on aspect ratio
	newHeight := uint(float64(img.Bounds().Dy()) / float64(img.Bounds().Dx()) * float64(width) * 0.5)
	img = resize.Resize(uint(width), newHeight, img, resize.Lanczos3)

	asciiArt := imageToColorASCII(img)
	err = saveToFile("ascii_output.txt", asciiArt)
	if err != nil {
		fmt.Println("Error saving ASCII art:", err)
	}

	// Convert and print ASCII
	fmt.Println(imageToColorASCII(img))
}

func imageToColorASCII(img image.Image) string {
	bounds := img.Bounds()
	output := ""

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := img.At(x, y)
			r, g, b, _ := col.RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			gray := grayscale(col)
			char := brightnessToChar(gray)

			colorized := termenv.String(string(char)).Foreground(termenv.RGBColor(
				fmt.Sprintf("#%02x%02x%02x", r8, g8, b8),
			))

			output += colorized.String()
		}
		output += "\n"
	}

	return output
}

func grayscale(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	return uint8(0.299*float64(r/256) + 0.587*float64(g/256) + 0.114*float64(b/256))
}

func brightnessToChar(brightness uint8) byte {
	scale := float64(brightness) / 255.0
	index := int(scale * float64(len(asciiChars)-1))
	return asciiChars[index]
}

func saveToFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}
