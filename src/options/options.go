package options

import (
	"flag"
)

type options struct {
	SlideshowDirectory    string
	OutputDirectory       string
	TemporaryDirectory    string
	OverlayVideoDirectory string
	FFmpegDirectory       string
	LowQuality            bool
	SaveTemps             bool
	UseOldFade            bool
	Verbose               bool
}

/* Function to parse the command line options flags
 *  Returns:
 *			initalized options struct
 */
func ParseFlags() options {
	var slideshowDirectory string
	var outputDirectory string
	var temporaryDirectory string
	var overlayVideoDirectory string
	var ffmpegDirectory string
	var lowQuality bool
	var saveTemps bool
	var useOldFade bool
	var verbose bool

	flag.BoolVar(&lowQuality, "l", false, "Low Quality, include to generate a lower quality video (480p instead of 720p)")
	flag.BoolVar(&saveTemps, "s", false, "Save Temporaries, include to save temporary files generated during video process)")
	flag.BoolVar(&useOldFade, "f", false, "Fadetype, include to use the non-xfade default transitions for video")
	flag.BoolVar(&verbose, "v", false, "Verbose, include to increase the verbosity of the feedback provided")

	flag.StringVar(&slideshowDirectory, "t", "", "Template Name, specify a template to use (if not included searches current folder for template)")
	flag.StringVar(&outputDirectory, "o", "", "Output Location, specify where to store final result (default is current directory)")
	flag.StringVar(&temporaryDirectory, "td", "", "Temporary Directory, used to specify a location to store the temporary files used in video production (default is OS' temp folder/storybuilder-*)")
	flag.StringVar(&overlayVideoDirectory, "ov", "", "Overlay Video, specify test video location to create overlay video")
	flag.StringVar(&ffmpegDirectory, "fd", "", "FFmpeg directory, specify location of ffmpeg binaries on your machine")
	flag.Parse()

	options := options{slideshowDirectory, outputDirectory, temporaryDirectory, overlayVideoDirectory, ffmpegDirectory, lowQuality, saveTemps, useOldFade, verbose}

	return options

}

/* Function to set the slideshow directory of the options struct
 *  Parameters:
 *			directory (string) : directory of the slideshow file
 */
func (o *options) SetSlideshowDirectory(directory string) {
	o.SlideshowDirectory = directory
}
