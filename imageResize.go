package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

func sl(oriSli []os.FileInfo, cpuNum int, toSlice [][]os.FileInfo) {
	perlen := len(oriSli)/cpuNum + 1
	for i := 0; i < cpuNum-1; i++ {

		toSlice[i] = oriSli[i*perlen : (i+1)*perlen]
	}
	toSlice[cpuNum-1] = oriSli[(cpuNum-1)*perlen:]

}

func resize(s []os.FileInfo, c chan int) {
	for _, v := range s {
		if strings.HasSuffix(v.Name(), ".jpg") || strings.HasSuffix(v.Name(), ".JPG") {
			f, err := imaging.Open(v.Name())
			if err != nil {
				log.Fatalln("Open file fail!", v.Name())
			}
			outf := imaging.Resize(f, f.Bounds().Dx(), 0, imaging.Lanczos)
			imaging.Save(outf, v.Name(), imaging.JPEGQuality(80))
			runtime.Gosched()

		}

	}

	c <- 1
}

func main() {

	stime := time.Now()
	if len(os.Args) == 1 {
		log.Fatal("必须输入一个目录")
	}
	wd := os.Args[1]

	err := os.Chdir(wd)
	if err != nil {
		log.Fatal(err)
	}
	dirs, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatalln("Dirs get fail!", wd)
	}
	fmt.Println("Total file num is ", len(dirs))
	fmt.Println("Now starts resizing!")
	cpuNum := runtime.NumCPU()
	oriSli := dirs
	c := make(chan int, 2)
	toSlice := make([][]os.FileInfo, cpuNum)
	sl(oriSli, cpuNum, toSlice)
	for i := 0; i < cpuNum; i++ {
		go resize(toSlice[i], c)
	}
	for k := 0; k < cpuNum; k++ {
		<-c
	}

	etime := time.Now()
	fmt.Println(etime.Sub(stime))
	time.Sleep(10 * time.Second)
}
