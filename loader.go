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
	"fmt"
	"github.com/disintegration/gift"
	_ "github.com/disintegration/imaging"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	//"github.com/NOX73/go-neural"
	//"github.com/NOX73/go-neural/learn"
	//"github.com/sajari/random-forest/RF"
)

func randomArray(n int) []float32 {
	ret := make([]float32, n)
	for i := 0; i < n; i++ {
		ret[i] = (rand.Float32() - 0.5) * 5
	}
	return ret
}

func main() {
	log.Print("use log so we don't have to put an underscore before the import")

	type Example struct {
		features []float32
		category string
	}

	paths := make(chan string)
	processed := make(chan Example)

	// the Seed for the network is 7
	rand.Seed(7)

	g := gift.New(
		// edge detector
		gift.Convolution(
			[]float32{
				-1, -1, -1,
				-1, 8, -1,
				-1, -1, -1,
			},
			false, false, false, 0.0),
		// is this max pool?
		gift.Maximum(2, true),
		gift.Resize(50, 0, gift.LinearResampling),

		// random 5x5 conv, hmm but like the color channels bro this is a shit neural network
		gift.Convolution(
			randomArray(25),
			false, false, false, 0.0),
		// is this max pool?
		gift.Maximum(2, true),
		gift.Resize(25, 0, gift.LinearResampling),

		// random 3x3 conv, hmm but like the color channels bro this is a shit neural network
		gift.Convolution(
			randomArray(9),
			false, false, false, 0.0),
		// is this max pool?
		gift.Maximum(2, true),
		gift.Resize(10, 0, gift.LinearResampling),

		// 300 features one for each spartan RIP
	)

	//n := neural.NewNetwork(300, []int{100,20,1})
	//n.RandomizeSynapses()

	// forest builder
	go func() {
		// is this a proper design pattern?
		// probs not it's awkward ROS node shit
		for {
			sample := <-processed
			fmt.Println(sample)

			// ugh no inline if?
			/*prob := []float64{0.01}
			  if sample.yes {
			    prob = []float64{0.4}
			  }

			  learn.Learn(n, sample.features, prob, 0.05)

			  println(prob[0], learn.Evaluation(n, sample.features, prob))*/
		}
	}()

	// image loader and network runner
	go func() {
		for {
			path := <-paths

			// load the image, this is 5 lines
			// i hate all this error handling does go have exceptions?
			f, err := os.Open(path)
			if err != nil {
				log.Fatal(err, path)
			}
			img, err := png.Decode(f)
			if err != nil {
				log.Fatal(err, path)
			}
			f.Close()

			dst := image.NewRGBA(g.Bounds(img.Bounds()))
			g.Draw(dst, img)

			// extract features
			// i can write much better than this wow shit
			ret := make([]float32, 300)
			cnt := 0
			for i := 0; i < 400; i++ {
				if i%4 == 3 {
					continue
				}
				ret[cnt] = float32(dst.Pix[i]) / 256.0
				cnt += 1
			}

			processed <- Example{features: ret, category: strings.Split(path, "/")[1]}

			//imaging.Save(dst, "dst.png")
			//println(dst)
		}
	}()

	files := []string{}

	filepath.Walk("imgs/", func(path string, finfo os.FileInfo, err error) error {
		if finfo.IsDir() {
			return nil
		}
		//paths <- path
		files = append(files, path)
		return nil
	})
	fmt.Println("files list built")
	println(len(files))

	/*perm := rand.Perm(len(files))
	  for _, v := range perm {
	    paths <- files[v]
	  }*/
}
