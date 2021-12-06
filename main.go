package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	//"os/exec"
)

// File Location of Repository **CHANGE THIS FILEPATH TO YOUR REPOSITORY FILEPATH**
// var basePath = "/Users/gordon.loaner/OneDrive - Gordon College/Desktop/Gordon/Senior/Senior Project/SIL-Video" //sehee
var basePath = "/Users/hyungyu/Documents/SIL-Video" //hyungyu
// var basePath = "C:/Users/damar/Documents/GitHub/SIL-Video" // david
// var basePath = "/Users/roddy/Desktop/SeniorProject/SIL-Video/"

//location of where you downloaded FFmpeg
// var baseFFmpegPath = "C:/FFmpeg" //windows
var baseFFmpegPath = "/usr/local/" //mac

var FfmpegBinPath = baseFFmpegPath + "/bin/ffmpeg"
var FfprobeBinPath = baseFFmpegPath + "/bin/ffprobe"

func main() {
	// First we read in the input file and parse the json
	//convertToVideo()

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
	var img1_motion_start_left = img1_motion_start_slice[0]
	var img1_motion_start_top = img1_motion_start_slice[1]
	var img1_motion_start_width = img1_motion_start_slice[2]
	var img1_motion_start_height = img1_motion_start_slice[3]
	var img1_motion_end_left = img1_motion_end_slice[0]
	var img1_motion_end_top = img1_motion_end_slice[1]
	var img1_motion_end_width = img1_motion_end_slice[2]
	var img1_motion_end_height = img1_motion_end_slice[3]
	// fmt.Println(img1_motion_start)
	// fmt.Println(img1_motion_end)
	// fmt.Println(img1_motion_start_left)
	// fmt.Println(img1_motion_start_top)
	// fmt.Println(img1_motion_start_width)
	// fmt.Println(img1_motion_start_height)
	// fmt.Println(img1_motion_end_left)
	// fmt.Println(img1_motion_end_top)
	// fmt.Println(img1_motion_end_width)
	// fmt.Println(img1_motion_end_height)

	var img2_motion_start = slideshow.Slide[2].Motion.Start
	var img2_motion_end = slideshow.Slide[2].Motion.End
	img2_motion_start_slice := strings.Split(img2_motion_start, " ")
	img2_motion_end_slice := strings.Split(img2_motion_end, " ")
	var img2_motion_start_left = img2_motion_start_slice[0]
	var img2_motion_start_top = img2_motion_start_slice[1]
	var img2_motion_start_width = img2_motion_start_slice[2]
	var img2_motion_start_height = img2_motion_start_slice[3]
	var img2_motion_end_left = img2_motion_end_slice[0]
	var img2_motion_end_top = img2_motion_end_slice[1]
	var img2_motion_end_width = img2_motion_end_slice[2]
	var img2_motion_end_height = img2_motion_end_slice[3]

	var img3_motion_start = slideshow.Slide[3].Motion.Start
	var img3_motion_end = slideshow.Slide[3].Motion.End
	img3_motion_start_slice := strings.Split(img3_motion_start, " ")
	img3_motion_end_slice := strings.Split(img3_motion_end, " ")
	var img3_motion_start_left = img3_motion_start_slice[0]
	var img3_motion_start_top = img3_motion_start_slice[1]
	var img3_motion_start_width = img3_motion_start_slice[2]
	var img3_motion_start_height = img3_motion_start_slice[3]
	var img3_motion_end_left = img3_motion_end_slice[0]
	var img3_motion_end_top = img3_motion_end_slice[1]
	var img3_motion_end_width = img3_motion_end_slice[2]
	var img3_motion_end_height = img3_motion_end_slice[3]

	// Place them all inside a string slice
	// paths := []string{outputPath, titleimg, img1, img2, img3, introAudio, audio1, title_start, title_duration, img1_start, img1_duration, img2_start, img2_duration, img3_start, img3_duration}
	paths := []string{outputPath, titleimg, img1, img2, img3, introAudio, audio1, title_start, title_duration, img1_start, img1_duration, img2_start, img2_duration, img3_start, img3_duration, img1_motion_start_left, img1_motion_start_top, img1_motion_start_width, img1_motion_start_height, img1_motion_end_left, img1_motion_end_top, img1_motion_end_width, img1_motion_end_height, img2_motion_start_left, img2_motion_start_top, img2_motion_start_width, img2_motion_start_height, img2_motion_end_left, img2_motion_end_top, img2_motion_end_width, img2_motion_end_height, img3_motion_start_left, img3_motion_start_top, img3_motion_start_width, img3_motion_start_height, img3_motion_end_left, img3_motion_end_top, img3_motion_end_width, img3_motion_end_height}

	createTempVideos(paths...)
	combineVideos()

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
			"-framerate", "1", // frame  to define how fast the pictures are read in, in this case, 1 picture per second
			"-i", fmt.Sprintf("%s/image-%d.jpg", basePath, i), // input image
			"-r", "30", // the framerate of the output video
			"-ss", paths[9+2*i-2]+"ms",
			"-t", paths[10+2*i-2]+"ms",
			"-i", basePath+"/narration-001.mp3", // input audio
			"-vf", "zoompan=z='zoom + 0.001':x='"+paths[15+8*i-8]+"':y='"+paths[16+8*i-8]+"':iw='"+paths[17+8*i-8]+"':ih='"+paths[18+8*i-8]+"':ow='"+paths[21+8*i-8]+"':oh='"+paths[22+8*i-8]+"':x='"+paths[19+8*i-8]+"':y='"+paths[20+8*i-8]+"'",
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
	findVideos()

	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", basePath+"/output/text.txt",
		basePath+"/output/mergedVideo.mp4",
	)

	err := cmd.Run() // Start a process on another goroutine
	check(err)
}
