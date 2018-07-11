package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/disintegration/imaging"
)

type sum struct {
	sum int
}

func resize(src string, status chan int) {
	if strings.HasSuffix(src, ".jpg") || strings.HasSuffix(src, ".JPG") {
		f, err := imaging.Open(src)
		if err != nil {
			log.Fatalln("Open file fail! ", src)
		}
		fmt.Println("Now processing ", src, "........")
		os.Remove(src)
		outf := imaging.Resize(f, f.Bounds().Dx(), 0, imaging.Lanczos)
		imaging.Save(outf, src, imaging.JPEGQuality(80))
		status <- 1
	}
}

func distribut(file string, id *sum, status chan int) {
	<-status
	id.sum++
	fmt.Println("Current file ID is ", id.sum)
	go resize(file, status)
}

func main() {
	//Init start flag
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	status := make(chan int, cpus)
	for i := 0; i < cpus; i++ {
		status <- 1
	}

	err := os.Chdir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	a := &sum{0}
	for _, v := range files {
		distribut(v.Name(), a, status)
	}

}
