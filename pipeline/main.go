package main

import (
	"fmt"
	"image"
	"strings"

	imageprocessing "github.com/kcalixto/poc-go-concurrency/pipeline/image_processing"
)

type Job struct {
	InputPath  string
	Image      image.Image
	OutputPath string
}

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		for _, p := range paths {
			job := Job{
				InputPath:  p,
				OutputPath: strings.Replace(p, "images/", "images/output/", 1),
				Image:      imageprocessing.ReadImage(p),
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			job.Image = imageprocessing.Resize(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			job.Image = imageprocessing.Grayscale(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func saveImage(input <-chan Job) <-chan bool {
	out := make(chan bool)
	go func() {
		for job := range input {
			imageprocessing.WriteImage(job.OutputPath, job.Image)
			out <- true
		}
		close(out)
	}()
	return out
}

func main() {
	imagePaths := []string{
		"images/image1.jpg",
		"images/image2.jpg",
		"images/image3.jpg",
		"images/image4.jpg",
		"images/image5.jpg",
	}

	channel1 := loadImage(imagePaths)
	channel2 := resize(channel1)
	channel3 := convertToGrayscale(channel2)
	writeResults := saveImage(channel3)

	// continuously read from the writeResults channel
	for success := range writeResults {
		if success {
			fmt.Println("Success!")
		} else {
			fmt.Println("Failed!")
		}
	}
}
