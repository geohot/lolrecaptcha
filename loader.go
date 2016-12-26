package main

/*
so like wow there's no neural networks for go, CNN anyone?
idea is this and if you cheat and use python you are a big cheater
give all alleged street sign images 0.4 of street sign
and give all other images 0.01 chance of street sign
and maybe with the magic of neural networks we will learn?

TODO: don't be cheater and use python only golang pull request accepted
*/

import (
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func main() {
	imgarr := []image.Image{}
	filepath.Walk("imgs/", func(path string, finfo os.FileInfo, err error) error {
		if finfo.IsDir() {
			return nil
		}
		println(path)
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		img, err := png.Decode(f)
		imgarr = append(imgarr, img)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
		return nil
	})
}
