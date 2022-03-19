package main

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	inputFile := "test.slideshow"
	var expectedOutput string
	data := readData(inputFile)
	for i, slide := range data.Slide {
		if i == 0 {
			// Test background filename
			expectedOutput = "background.mp3"
			backFilename := slide.Audio.Background_Filename.Path
			if backFilename != expectedOutput {
				t.Error(fmt.Sprintf("expected background filename to be %s, but got %s", expectedOutput, backFilename))
			}
		} else {
			expectedOutput = "narration.mp3"
			audio := slide.Audio.Filename.Name
			if audio != expectedOutput {
				t.Error(fmt.Sprintf("expected audio filename to be %s, but got %s", expectedOutput, audio))
			}
		}
		expectedOutput = fmt.Sprintf("test-%d.jpg", i)
		image := slide.Image.Name
		if image != expectedOutput {
			t.Error(fmt.Sprintf("expected image filename to be %s, but got %s", expectedOutput, image))
		}
		if slide.Motion.Start != "" {
			expectedOutput = "0.0 0.1 0.2 0.3"
			start := slide.Motion.Start
			if start != expectedOutput {
				t.Error(fmt.Sprintf("expected motion start to be %s, but got %s", expectedOutput, start))
			}
			expectedOutput = "1 2 3 4"
			end := slide.Motion.End
			if end != expectedOutput {
				t.Error(fmt.Sprintf("expected motion end to be %s, but got %s", expectedOutput, end))
			}
		}
		if slide.Transition.Type != "" {
			expectedOutput = "transitionTest"
			transitionType := slide.Transition.Type
			if transitionType != expectedOutput {
				t.Error(fmt.Sprintf("expected transtion type to be %s, but got %s", expectedOutput, transitionType))
			}
			expectedOutput = "1000"
			transitionDuration := slide.Transition.Duration
			if transitionDuration != expectedOutput {
				t.Error(fmt.Sprintf("expected transtion duration to be %s, but got %s", expectedOutput, transitionDuration))
			}
		}
		if slide.Timing.Start != "" {
			expectedOutput = "1234"
			timingStart := slide.Timing.Start
			if timingStart != expectedOutput {
				t.Error(fmt.Sprintf("expected timing start to be %s, but got %s", expectedOutput, timingStart))
			}
			expectedOutput = "5678"
			timingDuration := slide.Timing.Duration
			if timingDuration != expectedOutput {
				t.Error(fmt.Sprintf("expected timing duration to be %s, but got %s", expectedOutput, timingDuration))
			}
		}

	}
}

// expected output should be a png
func TestScaleImage(t *testing.T) {
	inputFile := input_images(Images[i])
	input := height
	input2 := width
	expectedOutput := fmt.Sprintf("test_%d.jpg",i)
	if inputFile != expectedOutput {
		t.Errorf("expected image here")
	}
}
}

// func TestReadFile(t *testing.T) {
// 	data, err := ioutil.ReadFile("data.slideshow")
// 	if err != nil {

// 	}
// 	if string(readData) != nil {

// 	}
//// }//