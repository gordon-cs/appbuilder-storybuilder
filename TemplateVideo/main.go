package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var templateName string

func main() {
	var fadeType string
	flag.StringVar(&templateName, "t", "", "Specify template to use.")
	flag.StringVar(&fadeType, "f", "", "Specify transition type (x) for xfade, leave blank for old fade")
	flag.Parse()
	if templateName == "" {
		fmt.Println("No template provided, searching local folder...")
		filepath.WalkDir(".", findTemplate)
	}
	start := time.Now()
	// First we parse in the various pieces from the template
	Images := []string{}
	Audios := []string{}
	BackAudioPath := ""
	BackAudioVolume := ""
	Transitions := []string{}
	TransitionDurations := []string{}
	Timings := [][]string{}
	fmt.Println("Parsing .slideshow file...")
	var slideshow = readData(templateName)
	for _, slide := range slideshow.Slide {
		if slide.Audio.Background_Filename.Path != "" {
			Audios = append(Audios, slide.Audio.Background_Filename.Path)
			BackAudioPath = slide.Audio.Background_Filename.Path
			BackAudioVolume = slide.Audio.Background_Filename.Volume
		} else {
			if slide.Audio.Filename.Name == "" {
				Audios = append(Audios, "")
			} else {
				Audios = append(Audios, slide.Audio.Filename.Name)
			}
		}
		Images = append(Images, slide.Image.Name)

		if slide.Transition.Type == "" {
			Transitions = append(Transitions, "fade")
		} else {
			Transitions = append(Transitions, slide.Transition.Type)
		}
		if slide.Transition.Duration == "" {
			TransitionDurations = append(TransitionDurations, "1000")
		} else {
			TransitionDurations = append(TransitionDurations, slide.Transition.Duration)
		}

		temp := []string{slide.Timing.Start, slide.Timing.Duration}
		Timings = append(Timings, temp)
	}
	fmt.Println("Parsing completed...")
	fmt.Println("Scaling Images...")
	//	scaleImages(Images, "1500", "900")
	fmt.Println("Creating video...")

	//if using xfade
	if fadeType == "xfade" {
		//allImages := make_temp_videos_with_audio(Images, Transitions, TransitionDurations, Timings, Audios)
		//make_temp_videos_with_audio(Images, Transitions, TransitionDurations, Timings, Audios)
		//for testing
		//allImages := []int{0, 1, 2, 3, 4, 5, 6}

		//mergeVideos(allImages, Images, Transitions, TransitionDurations, Timings, 0)
		merge_videos_once(Images, Transitions, TransitionDurations, Timings)
	} else {
		combineVideos(Images, Transitions, TransitionDurations, Timings, Audios)
		fmt.Println("Adding intro music...")
		addBackgroundMusic(BackAudioPath, BackAudioVolume)
	}

	fmt.Println("Finished making video...")
	duration := time.Since(start)
	fmt.Println("Video completed!")
	fmt.Println(fmt.Sprintf("Time Taken: %f seconds", duration.Seconds()))
}

// Function to check errors from non-CMD output
func check(err error) {
	if err != nil {
		fmt.Println("Error", err)
		log.Fatalln(err)
	}
}

// Function to check CMD error output when running commands
func checkCMDError(output []byte, err error) {
	if err != nil {
		log.Fatalln(fmt.Sprint(err) + ": " + string(output))
	}
}

/* Function to scale all the input images to a uniform height/width
*  to prevent issues in the video creation process
 */
func scaleImages(Images []string, height string, width string) {
	var wg sync.WaitGroup
	// Tell the 'wg' WaitGroup how many threads/goroutines
	//   that are about to run concurrently.
	wg.Add(len(Images))

	for i := 0; i < len(Images); i++ {
		go func(i int) {
			defer wg.Done()
			cmd := exec.Command("ffmpeg", "-i", "./"+Images[i],
				"-vf", fmt.Sprintf("scale=%s:%s", height, width)+",setsar=1:1",
				"-y", "./"+Images[i])
			output, err := cmd.CombinedOutput()
			checkCMDError(output, err)
		}(i)
	}

	wg.Wait()
}

/* Function to find the .slideshow template if none provided
 */
func findTemplate(s string, d fs.DirEntry, err error) error {
	slideRegEx := regexp.MustCompile(`.+(.slideshow)$`)
	if err != nil {
		return err
	}
	if slideRegEx.MatchString(d.Name()) {
		fmt.Println("Found template: " + s + "\nUsing found template...")
		templateName = s
	}
	return nil
}

/** Function to create the video with all images + transitions
*	Parameters:
*		Images: ([]string) - Array of filenames for the images
*		Transitions: ([]string) - Array of XFade transition names to use
*		TransitionDurations: ([]string) - Array of durations for each transition
*		Timings: ([][]string) - 2-D array of timing data for the audio for each image
*		Audios: ([]string) - Array of filenames for the audios to be used
 */
func combineVideos(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, Audios []string) {
	input_images := []string{}
	input_filters := ""
	totalNumImages := len(Images)
	concatTransitions := ""

	fmt.Println("Getting list of images and filters...")
	for i := 0; i < totalNumImages; i++ {
		// Everything needs to be concatenated so always add the image to concatTransitions
		concatTransitions += fmt.Sprintf("[v%d]", i)
		// Everything needs to be cropped so add the crop filter to every image
		input_filters += fmt.Sprintf("[%d:v]crop=trunc(iw/2)*2:trunc(ih/2)*2", i)
		if i == totalNumImages-1 { // Credits image has no timings/transitions
			input_images = append(input_images, "-i", "./"+Images[i])
		} else {
			input_images = append(input_images, "-loop", "1", "-ss", Timings[i][0]+"ms", "-t", Timings[i][1]+"ms", "-i", "./"+Images[i])

			if i == 0 {
				input_filters += fmt.Sprintf(",fade=t=out:st=%sms:d=%sms", Timings[i][1], TransitionDurations[i])
			} else {
				half_duration, err := strconv.Atoi(TransitionDurations[i])
				check(err)
				input_filters += fmt.Sprintf(",fade=t=in:st=0:d=%dms,fade=t=out:st=%sms:d=%dms", half_duration/2, Timings[i][1], half_duration/2)
			}
		}
		input_filters += fmt.Sprintf("[v%d];", i)

	}

	concatTransitions += fmt.Sprintf("concat=n=%d:v=1:a=0,format=yuv420p[v]", totalNumImages)
	input_filters += concatTransitions

	input_images = append(input_images, "-i", "./narration-001.mp3",
		"-max_muxing_queue_size", "9999",
		"-filter_complex", input_filters, "-map", "[v]",
		"-map", fmt.Sprintf("%d:a", totalNumImages),
		"-shortest", "-y", "../output/mergedVideo.mp4")

	fmt.Println("Creating video...")
	cmd := exec.Command("ffmpeg", input_images...)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}

func addBackgroundMusic(backgroundAudio string, backgroundVolume string) {
	tempVol := 0.0
	// Convert the background volume to a number between 0 and 1, if it exists
	if backgroundVolume != "" {
		if s, err := strconv.ParseFloat(backgroundVolume, 64); err == nil {
			tempVol = s
		} else {
			fmt.Println("Error converting volume to float")
		}
		tempVol = tempVol / 100
	} else {
		tempVol = .5
	}
	cmd := exec.Command("ffmpeg",
		"-i", "../output/mergedVideo.mp4",
		"-i", backgroundAudio,
		"-filter_complex", "[1:0]volume="+fmt.Sprintf("%f", tempVol)+"[a1];[0:a][a1]amix=inputs=2:duration=first",
		"-map", "0:v:0",
		"-y", "../output/finalvideo.mp4",
	)
	output, e := cmd.CombinedOutput()
	checkCMDError(output, e)
}

func make_temp_videos_with_audio(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, Audios []string) []int {
	totalNumImages := len(Images)

	cmd := exec.Command("")

	allImages := []int{}

	var wg sync.WaitGroup
	// Tell the 'wg' WaitGroup how many threads/goroutines
	//   that are about to run concurrently.
	wg.Add(totalNumImages)

	for i := 0; i < totalNumImages; i++ {
		// Spawn a thread for each iteration in the loop.
		// Pass 'i' into the goroutine's function
		//   in order to make sure each goroutine
		//   uses a different value for 'i'.
		go func(i int) {
			// At the end of the goroutine, tell the WaitGroup
			//   that another thread has completed.
			defer wg.Done()

			totalDuration := 0
			for j := 1; j < i; j++ {
				if Audios[i] == Audios[j] {
					duration, err := strconv.Atoi(Timings[j][1])
					check(err)
					totalDuration += duration
				}
			}

			if Audios[i] == "" {
				fmt.Printf("Making temp%d-%d.mp4 video with empty audio\n", i, totalNumImages)
				cmd = exec.Command("ffmpeg", "-loop", "1", "-ss", "0ms", "-t", "3000ms", "-i", Images[i],
					"-f", "lavfi", "-i", "aevalsrc=0", "-t", "3000ms",
					"-shortest", "-pix_fmt", "yuv420p",
					"-y", fmt.Sprintf("../output/temp%d-%d.mp4", i, totalNumImages))
			} else {
				fmt.Printf("Making temp%d-%d.mp4 video\n", i, totalNumImages)
				totalDuration_int := strconv.Itoa(totalDuration)

				cmd = exec.Command("ffmpeg", "-loop", "1", "-ss", "0ms", "-t", Timings[i][1]+"ms", "-i", "./"+Images[i],
					"-ss", totalDuration_int+"ms", "-t", Timings[i][1]+"ms", "-i", Audios[i],
					"-shortest", "-pix_fmt", "yuv420p", "-y", fmt.Sprintf("../output/temp%d-%d.mp4", i, totalNumImages))

			}

			output, err := cmd.CombinedOutput()
			checkCMDError(output, err)
		}(i)

		allImages = append(allImages, i)
	}

	// Wait for `wg.Done()` to be exectued the number of times
	//   specified in the `wg.Add()` call.
	// `wg.Done()` should be called the exact number of times
	//   that was specified in `wg.Add()`.
	wg.Wait()
	return allImages
}

func mergeVideos(items []int, Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, depth int) []int {
	if len(items) < 2 {
		return items
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	first := []int{}

	go func() {
		defer wg.Done()
		first = mergeVideos(items[:len(items)/2], Images, Transitions, TransitionDurations, Timings, depth+1)
	}()

	second := mergeVideos(items[len(items)/2:], Images, Transitions, TransitionDurations, Timings, depth+1)

	wg.Wait()

	return merge(first, second, Images, Transitions, TransitionDurations, Timings, depth)
}

func merge(a []int, b []int, Images []string, Transitions []string, TransitionDurations []string, Timings [][]string, depth int) []int {
	final := []int{}
	i := 0
	j := 0

	if len(a) == 1 && len(b) == 1 {
		//combine the individual temporary videos into merged files
		totalNumImages := len(Images)
		transition := Transitions[a[0]]

		transition_duration, err := strconv.Atoi(TransitionDurations[a[0]])
		check(err)
		transition_duration_float := float64(transition_duration) / 1000

		duration, err := strconv.Atoi(Timings[a[0]][1])
		offset := (float64(duration) - float64(transition_duration)) / 1000

		//check if video has full or partial audio
		cmd := exec.Command("ffprobe",
			"-i", fmt.Sprintf("../output/temp%d-%d.mp4", a[0], totalNumImages),
			"-v", "error", "-of", "flat=s_",
			"-select_streams", "1", "-show_entries", "stream=duration", "-of", "default=noprint_wrappers=1:nokey=1")

		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)

		audio_duration_offset, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)

		if offset-audio_duration_offset > 1 {
			//the calculated offset is more than 1 seconds longer than the true duration of the video
			offset = offset - float64(transition_duration/1000)
		}

		fmt.Printf("Combining videos temp%d-%d.mp4 and temp%d-%d.mp4 with %s transition to merged%d-%d. \n", a[0], totalNumImages, b[0], totalNumImages, transition, a[0], depth)

		cmd = exec.Command("ffmpeg",
			"-i", fmt.Sprintf("../output/temp%d-%d.mp4", a[0], totalNumImages),
			"-i", fmt.Sprintf("../output/temp%d-%d.mp4", b[0], totalNumImages),
			"-filter_complex",
			fmt.Sprintf("[0:v]settb=AVTB,fps=30/1[v0];[1:v]settb=AVTB,fps=30/1[v1];[v0][v1]xfade=transition=%s:duration=%f:offset=%f,format=yuv420p[outv];[0:a][1:a]acrossfade=duration=%dms:o=0:curve1=nofade:curve2=nofade[outa]", transition, transition_duration_float, offset, transition_duration),
			"-map", "[outv]",
			"-map", "[outa]",
			"-y", fmt.Sprintf("../output/merged%d-%d.mp4", a[0], depth),
		)

		output, err = cmd.CombinedOutput()
		checkCMDError(output, err)
	} else if len(a) == 1 && len(b) == 2 {
		//if odd number of things to merge, then it merges a temporary video with merged file
		totalNumImages := len(Images)
		index := len(a) - 1
		newDepth := depth

		newDepth++

		if depth == newDepth {
			newDepth = a[index]
		}

		transition := Transitions[a[index]]

		transition_duration, err := strconv.Atoi(TransitionDurations[a[0]])
		check(err)
		transition_duration_float := float64(transition_duration) / 1000

		duration, err := strconv.Atoi(Timings[a[0]][1])
		offset := (float64(duration) - float64(transition_duration)) / 1000

		fmt.Printf("Combining videos temp%d-%d.mp4 and merged%d-%d.mp4 with %s transition to merged%d-%d. \n", a[0], totalNumImages, b[0], newDepth, transition, a[0], depth)

		cmd := exec.Command("ffmpeg",
			"-i", fmt.Sprintf("../output/temp%d-%d.mp4", a[0], totalNumImages),
			"-i", fmt.Sprintf("../output/merged%d-%d.mp4", b[0], newDepth),
			"-filter_complex",
			fmt.Sprintf("[0:v]settb=AVTB,fps=30/1[v0];[1:v]settb=AVTB,fps=30/1[v1];[v0][v1]xfade=transition=%s:duration=%f:offset=%f,format=yuv420p[outv];[0:a][1:a]acrossfade=duration=%dms:o=0:curve1=nofade:curve2=nofade[outa]", transition, transition_duration_float, offset, transition_duration),
			"-map", "[outv]",
			"-map", "[outa]",
			"-y", fmt.Sprintf("../output/merged%d-%d.mp4", a[0], depth),
		)
		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	} else {
		//merging two merged files
		index := len(a) - 1
		newDepth := depth / 2

		if newDepth == 1 {
			newDepth = 2
		}
		if newDepth == 0 {
			newDepth = 1
		}

		depth++

		if depth == newDepth {
			newDepth = depth - 1
		}

		transition := Transitions[a[index]]

		transition_duration, err := strconv.Atoi(TransitionDurations[a[index]])
		check(err)

		transition_duration_float := float64(transition_duration) / 1000

		duration := 0

		for i := 0; i < len(a); i++ {
			duration_temp, err := strconv.Atoi(Timings[a[i]][1])
			check(err)

			duration += duration_temp
		}

		offset := (float64(duration) - float64(transition_duration*len(a))) / 1000

		fmt.Printf("Combining videos merged%d-%d.mp4 and merged%d-%d.mp4 with %s transition to merged%d-%d. \n", a[0], depth, b[0], depth, transition, a[0], newDepth)

		cmd := exec.Command("ffmpeg",
			"-i", fmt.Sprintf("../output/merged%d-%d.mp4", a[0], depth),
			"-i", fmt.Sprintf("../output/merged%d-%d.mp4", b[0], depth),
			"-filter_complex",
			fmt.Sprintf("[0:v]settb=AVTB,fps=30/1[v0];[1:v]settb=AVTB,fps=30/1[v1];[v0][v1]xfade=transition=%s:duration=%f:offset=%f,format=yuv420p[outv];[0:a][1:a]acrossfade=duration=%dms:o=0:curve1=nofade:curve2=nofade[outa]", transition, transition_duration_float, offset, transition_duration),
			"-map", "[outv]",
			"-map", "[outa]",
			"-y", fmt.Sprintf("../output/merged%d-%d.mp4", a[0], newDepth),
		)
		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)
	}

	for i < len(a) && j < len(b) {
		final = append(final, a[i])
		i++
	}

	for ; j < len(b); j++ {
		final = append(final, b[j])
	}
	return final
}

func merge_videos_once(Images []string, Transitions []string, TransitionDurations []string, Timings [][]string) {
	video_fade_filter := ""
	audio_fade_filter := ""
	settb := ""

	last_fade_output := "0:v"
	last_audio_output := "0:a"

	totalNumImages := len(Images)

	video_total_length := 0.0

	video_each_length := make([]float64, totalNumImages)

	input_files := []string{}

	for i := 0; i < totalNumImages; i++ {
		settb += fmt.Sprintf("[%d]settb=AVTB[%d:v];", i, i)
		input_files = append(input_files, "-i", fmt.Sprintf("../output/temp%d-%d.mp4", i, totalNumImages))
	}

	for i := 0; i < len(Images)-1; i++ {
		transition := Transitions[i]
		transition_duration, err := strconv.ParseFloat(strings.TrimSpace(string(TransitionDurations[i])), 8)
		transition_duration = transition_duration / 1000

		cmd := exec.Command("ffprobe",
			"-v", "error",
			"-show_entries", "format=duration",
			"-of", "default=noprint_wrappers=1:nokey=1",
			fmt.Sprintf("../output/temp%d-%d.mp4", i, totalNumImages),
		)
		output, err := cmd.CombinedOutput()
		checkCMDError(output, err)

		video_each_length[i], err = strconv.ParseFloat(strings.TrimSpace(string(output)), 8)

		video_total_length += video_each_length[i]
		next_fade_output := fmt.Sprintf("v%d%d", i, i+1)

		offset_from_previous_transitions := transition_duration * 1000 / 2
		offset_from_previous_transitions = (offset_from_previous_transitions * float64(i+1)) / 1000

		offset := video_total_length - (transition_duration)*(float64(i+1)) + offset_from_previous_transitions

		video_fade_filter += fmt.Sprintf("[%s][%d:v]xfade=transition=%s:duration=%.2f:offset=%.2f", last_fade_output, i+1,
			transition, transition_duration, offset)

		last_fade_output = next_fade_output

		if i < totalNumImages-2 {
			video_fade_filter += fmt.Sprintf("[%s];", next_fade_output)
		} else {
			video_fade_filter += ",format=yuv420p;"
		}

		next_audio_output := fmt.Sprintf("a%d%d", i, i+1)
		audio_fade_filter += fmt.Sprintf("[%s][%d:a]acrossfade=d=%.2f:o=0:curve2=nofade", last_audio_output, i+1, transition_duration)

		if i < totalNumImages-2 {
			audio_fade_filter += fmt.Sprintf("[%s];", next_audio_output)
		}

		last_audio_output = next_audio_output
	}

	input_files = append(input_files, "-filter_complex", settb+video_fade_filter+audio_fade_filter, "-y", "../output/final.mp4")

	cmd := exec.Command("ffmpeg", input_files...)

	output, err := cmd.CombinedOutput()
	checkCMDError(output, err)
}
