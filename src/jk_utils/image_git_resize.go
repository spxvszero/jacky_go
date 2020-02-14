package jk_utils

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"log"
	"os"
	"reflect"
)


func Thumbnails() {
	// open "test.jpg"
	//file, err := os.Open("/Users/jacky/Desktop/test.jpg")
	file, err := os.Open("/Users/jason/Desktop/goTest/test.png")
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	//img, err := jpeg.Decode(file)
	img, err:= png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	for i:=resize.NearestNeighbor;i <= resize.Lanczos3;i++  {
		creatImage(img,i)
	}

}

func creatImage(img image.Image, resizeType resize.InterpolationFunction){
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(0, 0, img, resizeType)

	//out, err := os.Create("test_resized.jpg")
	fmt.Println("?? ",reflect.TypeOf(resizeType))


	name := fmt.Sprintf("/Users/jason/Desktop/goTest/test_%d.png",resizeType)
	out, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	//jpeg.Encode(out, m, nil)
	png.Encode(out, m)

}