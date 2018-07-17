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

func resize(src string, status chan int) {
	if strings.HasSuffix(src, ".jpg") || strings.HasSuffix(src, ".JPG") {
		f, err := imaging.Open(src)
		if err != nil {
			log.Println("Open file fail! ", src)
			return
		}
		fSize1, _ := os.Stat(src)
		width := f.Bounds().Dx()
		if width > 4000 {
			width = 4000
		}
		outf := imaging.Resize(f, width, 0, imaging.Lanczos)
		os.Remove(src)
		imaging.Save(outf, src, imaging.JPEGQuality(80))
		fSize2, _ := os.Stat(src)
		fmt.Printf("%v\t Before: %vKb\t After: %vKb.\n", src, fSize1.Size()/(1<<10), fSize2.Size()/(1<<10))
		status <- 1
	}
}

func distribut(file string, status chan int) {
	<-status
	go resize(file, status)
}

func main() {
	//Init start flag
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	fmt.Printf("Now start Processing...The number of CPU is %v\n", cpus)
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
	for _, v := range files {
		distribut(v.Name(), status)
	}
	fmt.Println("Done!")
}
