package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

// Stores our image path, image object and ASCII slice after it has been
// converted
type ImageToASCII struct {
	ImagePath string
	Image     image.Image
	ASCII     [][]rune
}

// Runs the conversion process which will loop through the image and convert
// each pixel into a character. The passed width will be used as the width of
// the outputted image, the height will be automatically calculated.
func (img *ImageToASCII) Convert(width int) {

	// Reference: http://members.optusnet.com.au/astroblue/grey_scale.txt
	valueMap := []rune{
		' ',
		'.',
		'\'',
		',',
		';',
		'"',
		'o',
		'O',
		'%',
		'8',
		'@',
		'#'}

	// Grab the image dimensions
	bounds := img.Image.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	// @TODO(mark): Resize and stretch the image to account for the fact that
	//              the font is twice as tall as wide (roughly). Also take into
	//              account the width argument to calculate the required
	//              height.

	// Some precalculations
	img.ASCII = make([][]rune, h)
	valueMapWeight := float64(len(valueMap)-1) / 255

	for y := 0; y < h; y++ {

		// Assign memory for Y
		img.ASCII[y] = make([]rune, w)

		for x := 0; x < w; x++ {

			// @TODO(mark): Move this out into a goroutine with a channel which
			//              reads into our ASCII slice

			// Grab the gray value
			rgb := img.Image.At(x, y)
			gray := color.GrayModel.Convert(rgb).(color.Gray).Y

			// Pick the matching value
			bestMatch := math.Floor(valueMapWeight * float64(gray))
			img.ASCII[y][x] = valueMap[int(bestMatch)]
		}
	}
}

// Loops through the ASCII slice and prints out each rune with a new line after
// each row.
func (img *ImageToASCII) Print() {

	for y := 0; y < len(img.ASCII); y++ {

		// @TODO(mark): Remove this work around once the conversion algorithm
		//              has been fixed. This just halves x so it is twice as
		//              wide
		for x := 0.0; int(x) < len(img.ASCII[y]); x += 0.5 {

			fmt.Printf("%c", img.ASCII[y][int(x)])
		}
		fmt.Printf("\n")
	}
}

// Creates a new ImageToASCII struct, will open the image that is passed in and
// report any errors in the opening or decoding of it
func NewImageToASCII(imagePath string) (ImageToASCII, error) {

	reader, err := os.Open(imagePath)

	if err != nil {
		return ImageToASCII{}, err
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		return ImageToASCII{}, err
	}

	imgAscii := ImageToASCII{
		ImagePath: imagePath,
		Image:     img,
		ASCII:     make([][]rune, 0)}
	return imgAscii, nil
}

// Uses the code above
func main() {

	// Stores for flags
	var imagePath string
	var imageWidth int

	// Read in the path and image width flags
	flag.StringVar(&imagePath, "path", "", "Path of the image you wish to convert")
	flag.IntVar(&imageWidth, "width", 0, "Width to convert")
	flag.Parse()

	image, err := NewImageToASCII(imagePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Run conversion and immediatly print
	image.Convert(imageWidth)
	image.Print()
}
