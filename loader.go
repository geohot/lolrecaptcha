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
  "path/filepath"
  "image/png"
  "os"
  "log"
)


func main() {
  log.Print("use log so we don't have to put an underscore before the import")

  paths := make(chan string)

  // image loader and network runner 
  go func() {
    for {
      path := <-paths

      // load the image, this is 5 lines
      // i hate all this error handling does go have exceptions?
      f, err := os.Open(path)
      if err != nil { log.Fatal(err) }
      img, err := png.Decode(f)
      if err != nil { log.Fatal(err) }
      f.Close()

      println(img)
    }
  }()

  filepath.Walk("imgs/", func(path string, finfo os.FileInfo, err error) error {
    if finfo.IsDir() { return nil }
    paths <- path
    return nil
  });


}

