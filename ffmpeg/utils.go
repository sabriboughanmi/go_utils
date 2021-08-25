package ffmpeg

import (
	"time"
)

func (v *Video) clampToDuration(t time.Duration) time.Duration {
	if t < 0 {
		t = 0
	}
	if t > v.duration {
		t = v.duration
	}
	return t
}

func isEvenNumber(n int) bool {
	return n%2 == 0
}

func toEvenNumber(n int) int {
	if isEvenNumber(n) {
		return n
	}
	return n + 1
}

// commandLine returns the command line that will be used to convert the Video
// if you were to call Render.
func (v *EditableVideo) commandLine(output string) []string {

	additionalArgs := v.additionalArgs

	cmdline := []string{
		"ffmpeg",
		"-y",
		"-i", v.filepath,
		"-vcodec", "libx264",
		//	"-ss", strconv.FormatFloat(v.start.Seconds(), 'f', -1, 64),
		//	"-t", strconv.FormatFloat((v.end - v.start).Seconds(), 'f', -1, 64),
		//	"-vb", strconv.Itoa(v.bitrate),
	}
	cmdline = append(cmdline, additionalArgs...)
	cmdline = append(cmdline, output)
	return cmdline
}
