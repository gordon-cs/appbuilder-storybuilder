package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// File Location of Repository **CHANGE THIS FILEPATH TO YOUR REPOSITORY FILEPATH**
// var basePath = "/Users/gordon.loaner/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video" //sehee
var basePath = "/Users/hyungyu/Documents/SIL-Video" //hyungyu
//var basePath = "C:/Users/damar/Documents/GitHub/SIL-Video" // david
// var basePath = "/Users/roddy/Desktop/SeniorProject/SIL-Video/"

func main() {
	// First we parse in the various pieces from the template
	var outputPath = "./output"
	var slideshow = readData()

	var titleimg = slideshow.Slide[0].Image.Name
	var img1 = slideshow.Slide[1].Image.Name
	var img2 = slideshow.Slide[2].Image.Name
	var img3 = slideshow.Slide[3].Image.Name
	var introAudio = slideshow.Slide[0].Audio.Background_Filename.Path
	var audio1 = slideshow.Slide[1].Audio.Filename.Name
	var title_start = slideshow.Slide[0].Timing.Start
	var title_duration = slideshow.Slide[0].Timing.Duration
	var img1_start = slideshow.Slide[1].Timing.Start
	var img1_duration = slideshow.Slide[1].Timing.Duration
	var img2_start = slideshow.Slide[2].Timing.Start
	var img2_duration = slideshow.Slide[2].Timing.Duration
	var img3_start = slideshow.Slide[3].Timing.Start
	var img3_duration = slideshow.Slide[3].Timing.Duration

	var img1_motion_start = slideshow.Slide[1].Motion.Start
	var img1_motion_end = slideshow.Slide[1].Motion.End
	img1_motion_start_slice := strings.Split(img1_motion_start, " ")
	img1_motion_end_slice := strings.Split(img1_motion_end, " ")
	img1_motion_start_left, _ := strconv.ParseFloat(img1_motion_start_slice[0], 8) // string to float
	img1_motion_start_top, _ := strconv.ParseFloat(img1_motion_start_slice[1], 8)
	// img1_motion_start_width, _ := strconv.ParseFloat(img1_motion_start_slice[2], 8)
	img1_motion_start_height, _ := strconv.ParseFloat(img1_motion_start_slice[3], 8)
	img1_motion_end_left, _ := strconv.ParseFloat(img1_motion_end_slice[0], 8)
	img1_motion_end_top, _ := strconv.ParseFloat(img1_motion_end_slice[1], 8)
	// img1_motion_end_width, _ := strconv.ParseFloat(img1_motion_end_slice[2], 8)
	img1_motion_end_height, _ := strconv.ParseFloat(img1_motion_end_slice[3], 8)

	img1_timing_start, _ := strconv.ParseFloat(slideshow.Slide[1].Timing.Start, 8)
	img1_timing_duration, _ := strconv.ParseFloat(slideshow.Slide[1].Timing.Duration, 8)

	// var img2_motion_start = slideshow.Slide[2].Motion.Start
	// var img2_motion_end = slideshow.Slide[2].Motion.End
	// img2_motion_start_slice := strings.Split(img2_motion_start, " ")
	// img2_motion_end_slice := strings.Split(img2_motion_end, " ")
	// img2_motion_start_left, err := strconv.ParseFloat(img2_motion_start_slice[0], 8)
	// img2_motion_start_top, err := strconv.ParseFloat(img2_motion_start_slice[1], 8)
	// img2_motion_start_width, err := strconv.ParseFloat(img2_motion_start_slice[2], 8)
	// img2_motion_start_height, err := strconv.ParseFloat(img2_motion_start_slice[3], 8)
	// img2_motion_end_left, err := strconv.ParseFloat(img2_motion_end_slice[0], 8)
	// img2_motion_end_top, err := strconv.ParseFloat(img2_motion_end_slice[1], 8)
	// img2_motion_end_width, err := strconv.ParseFloat(img2_motion_end_slice[2], 8)
	// img2_motion_end_height, err := strconv.ParseFloat(img2_motion_end_slice[3], 8)

	// var img3_motion_start = slideshow.Slide[3].Motion.Start
	// var img3_motion_end = slideshow.Slide[3].Motion.End
	// img3_motion_start_slice := strings.Split(img3_motion_start, " ")
	// img3_motion_end_slice := strings.Split(img3_motion_end, " ")
	// img3_motion_start_left, err := strconv.ParseFloat(img3_motion_start_slice[0], 8)
	// img3_motion_start_top, err := strconv.ParseFloat(img3_motion_start_slice[1], 8)
	// img3_motion_start_width, err := strconv.ParseFloat(img3_motion_start_slice[2], 8)
	// img3_motion_start_height, err := strconv.ParseFloat(img3_motion_start_slice[3], 8)
	// img3_motion_end_left, err := strconv.ParseFloat(img3_motion_end_slice[0], 8)
	// img3_motion_end_top, err := strconv.ParseFloat(img3_motion_end_slice[1], 8)
	// img3_motion_end_width, err := strconv.ParseFloat(img3_motion_end_slice[2], 8)
	// img3_motion_end_height, err := strconv.ParseFloat(img3_motion_end_slice[3], 8)

	// generate params for ffmpeg
	var num_frames = ((img1_timing_duration - img1_timing_start) / (1000.0 / 30))
	num_frames_string := fmt.Sprintf("%f", num_frames)

	var size_init = img1_motion_start_height
	var size_change = img1_motion_end_height - img1_motion_start_height
	var size_incr = size_change / num_frames

	// var zoom_init = 1.0 / img1_motion_start_height
	// var zoom_change = 1.0/img1_motion_end_height - 1.0/img1_motion_start_height
	// var zoom_incr = zoom_change / num_frames

	var x_init = img1_motion_start_left
	var x_end = img1_motion_end_left
	var x_change = x_end - x_init
	var x_incr = x_change / num_frames

	var y_init = img1_motion_start_top
	var y_end = img1_motion_end_top
	var y_change = y_end - y_init
	var y_incr = y_change / num_frames

	zoom_cmd := fmt.Sprintf("1/(%f*%f*%f*on)", size_init-size_incr, checkSign(size_incr), math.Abs(size_incr))
	x_cmd := fmt.Sprintf("%f*iw*%f*%f*iw*on", x_init-x_incr, checkSign(x_incr), math.Abs(x_incr))
	y_cmd := fmt.Sprintf("%f*ih*%f*%f*ih*on", y_init-y_incr, checkSign(y_incr), math.Abs(y_incr))

	// Place them all inside a string slice
	// paths := []string{outputPath, titleimg, img1, img2, img3, introAudio, audio1, title_start, title_duration, img1_start, img1_duration, img2_start, img2_duration, img3_start, img3_duration}
	paths := []string{outputPath, titleimg, img1, img2, img3, introAudio, audio1, title_start, title_duration, img1_start, img1_duration, img2_start, img2_duration, img3_start, img3_duration, zoom_cmd, x_cmd, y_cmd, num_frames_string}

	createTempVideos(paths...)
	findVideos()
	combineVideos()
}

func checkSign(num float64) float64 {

	//return true for negative
	//return false for positive
	result := math.Signbit(num)

	if result == true {
		num = -1
	} else {
		num = 1
	}

	return num
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func createTempVideos(paths ...string) {
	fmt.Println(paths)
	for i := 1; i <= 3; i++ {
		cmd := exec.Command("ffmpeg",
			// "-i", fmt.Sprintf("%s/input/image-%d.jpg", basePath, i), // input image
			"-i", basePath+"/input/"+paths[i+1],
			"-r", "30", // the framerate of the output video
			"-ss", paths[9+2*i-2]+"ms",
			"-t", paths[10+2*i-2]+"ms",
			"-i", basePath+"/input/narration-001.mp3", // input audio
			"-vf", "zoompan=z="+paths[15+4*i-4]+":x="+paths[16+4*i-4]+":y="+paths[17+4*i-4]+":d="+paths[18+4*i-4]+":fps=30",
			"-pix_fmt", "yuv420p",
			"-vf", "crop=trunc(iw/2)*2:trunc(ih/2)*2",
			fmt.Sprintf("%s/output/output%d.mp4", basePath, i), // output
		)

		err := cmd.Start() // Start a process on another goroutine
		check(err)

		err = cmd.Wait() // wait until ffmpeg finishg
		check(err)
	}
}

func findVideos() {
	textfile, err := os.Create(basePath + "/output/text.txt")
	check(err)

	defer textfile.Close()

	files, err := ioutil.ReadDir(basePath + "/output")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".mp4") {
			textfile.WriteString("file ")
			textfile.WriteString(file.Name())
			textfile.WriteString("\n")
		}
	}

	textfile.Sync()
}

func combineVideos() {
	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", basePath+"/output/text.txt",
		basePath+"/output/mergedVideo.mp4",
	)

	err := cmd.Run() // Start a process on another goroutine
	check(err)
}
