package ffmpeg

import (
	"fmt"
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

	// Check if intro, outro, or background music are set
	hasIntro := v.introPath != nil
	hasOutro := v.outroPath != nil
	hasMusic := v.backgroundMusic != nil

	// If we have intro/outro or background music, we need filter_complex
	if hasIntro || hasOutro || hasMusic {
		cmdline := []string{
			"ffmpeg",
			"-y",
		}

		var filterComplex string
		var videoInputCount int
		var musicInputIndex int

		// Add video inputs (intro, main, outro)
		var videoPaths []string
		if hasIntro {
			videoPaths = append(videoPaths, *v.introPath)
		}
		videoPaths = append(videoPaths, v.filepath)
		if hasOutro {
			videoPaths = append(videoPaths, *v.outroPath)
		}

		for _, path := range videoPaths {
			cmdline = append(cmdline, "-i", path)
		}
		videoInputCount = len(videoPaths)

		// Add music input if present
		if hasMusic {
			cmdline = append(cmdline, "-i", v.backgroundMusic.FilePath)
			musicInputIndex = videoInputCount
		}

		// Build video concatenation filter if needed
		var videoOutputLabel, audioOutputLabel string

		if len(videoPaths) > 1 {
			// Scale all videos to the main video's resolution
			for i := 0; i < len(videoPaths); i++ {
				filterComplex += fmt.Sprintf("[%d:v]scale=%d:%d:force_original_aspect_ratio=decrease,setsar=1,pad=%d:%d:(ow-iw)/2:(oh-ih)/2[v%d];",
					i, v.width, v.height, v.width, v.height, i)
			}

			// Build concat part with video and audio streams
			for i := 0; i < len(videoPaths); i++ {
				filterComplex += fmt.Sprintf("[v%d][%d:a:0]", i, i)
			}

			// Add concat filter
			filterComplex += fmt.Sprintf("concat=n=%d:v=1:a=1[outv][videoa];", len(videoPaths))
			videoOutputLabel = "[outv]"
			audioOutputLabel = "[videoa]"
		} else {
			// Single video, just use input directly
			videoOutputLabel = "[0:v]"
			audioOutputLabel = "[0:a:0]"
		}

		// Build background music filter chain if needed
		if hasMusic {
			music := v.backgroundMusic

			// Calculate total video duration including intro and outro
			totalDuration := v.duration

			// Add intro duration if present
			if hasIntro {
				introVideo, err := LoadVideo(*v.introPath)
				if err == nil {
					totalDuration += introVideo.duration
				}
			}

			// Add outro duration if present
			if hasOutro {
				outroVideo, err := LoadVideo(*v.outroPath)
				if err == nil {
					totalDuration += outroVideo.duration
				}
			}

			// Build audio filter for music
			musicFilter := fmt.Sprintf("[%d:a]", musicInputIndex)

			// Apply start time offset if specified
			if music.StartTime > 0 {
				musicFilter += fmt.Sprintf("atrim=start=%.3f,asetpts=PTS-STARTPTS,", music.StartTime.Seconds())
			}

			// Apply looping if needed
			if music.Loop {
				// Calculate number of loops needed (add extra to be safe)
				musicFilter += "aloop=loop=-1:size=2e+09,"
			}

			// Apply volume adjustment
			musicFilter += fmt.Sprintf("volume=%.3f,", music.Volume)

			// Apply fade in if specified
			if music.FadeInDuration > 0 {
				musicFilter += fmt.Sprintf("afade=t=in:st=0:d=%.3f,", music.FadeInDuration.Seconds())
			}

			// Apply fade out if specified
			if music.FadeOutDuration > 0 {
				fadeOutStart := totalDuration.Seconds() - music.FadeOutDuration.Seconds()
				if fadeOutStart < 0 {
					fadeOutStart = 0
				}
				musicFilter += fmt.Sprintf("afade=t=out:st=%.3f:d=%.3f,", fadeOutStart, music.FadeOutDuration.Seconds())
			}

			// Trim to video duration
			musicFilter += fmt.Sprintf("atrim=0:%.3f[music];", totalDuration.Seconds())

			filterComplex += musicFilter

			// Mix music with video audio
			filterComplex += fmt.Sprintf("%s[music]amix=inputs=2:duration=first:dropout_transition=2[outa]", audioOutputLabel)
			audioOutputLabel = "[outa]"
		}

		// Add filter_complex if we built one
		if filterComplex != "" {
			cmdline = append(cmdline, "-filter_complex", filterComplex)
		}

		// Map outputs
		cmdline = append(cmdline, "-map", videoOutputLabel, "-map", audioOutputLabel)

		// Add codec
		cmdline = append(cmdline, "-vcodec", "libx264")

		// Add additional args (preset, crf, etc.)
		cmdline = append(cmdline, additionalArgs...)

		// Add output
		cmdline = append(cmdline, output)

		return cmdline
	}

	// Original simple command when no intro/outro/music
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
