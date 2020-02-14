package jk_utils

import (
	"fmt"
	"image/jpeg"
	"io"
)

func CompressImage(reader io.Reader,writer io.Writer) error  {


	img,err := jpeg.Decode(reader)

	if err != nil {
		fmt.Println("decode error ",err)
		return err
	}

	fmt.Println("read file size ",img.Bounds().Size())

	return nil

}

