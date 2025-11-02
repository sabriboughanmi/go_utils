package ffmpeg

import (
	"fmt"
	"sort"
	"time"
)

// buildVolumeExpression creates an FFmpeg volume filter expression from volume segments
func buildVolumeExpression(segments []VolumeSegment) string {
	if len(segments) == 0 {
		return "1.0" // Default volume
	}

	// Sort segments by start time
	sortedSegments := make([]VolumeSegment, len(segments))
	copy(sortedSegments, segments)
	sort.Slice(sortedSegments, func(i, j int) bool {
		return sortedSegments[i].StartTime < sortedSegments[j].StartTime
	})

	// Build nested if expression
	expr := ""
	defaultVolume := float32(1.0)

	// Build the expression from last to first (nested ifs)
	for i := len(sortedSegments) - 1; i >= 0; i-- {
		seg := sortedSegments[i]
		startSec := seg.StartTime.Seconds()
		endSec := seg.EndTime.Seconds()

		if seg.TransitionType == VolumeTransitionLinear && i > 0 {
			// Linear transition from previous segment's volume
			prevVol := sortedSegments[i-1].Volume
			transStartSec := startSec
			transEndSec := startSec + seg.TransitionDuration.Seconds()

			// Clamp transition end to segment end
			if transEndSec > endSec {
				transEndSec = endSec
			}

			// Build linear ramp expression: prevVol + (t-transStart)/(transEnd-transStart) * (newVol-prevVol)
			if transEndSec > transStartSec {
				rampExpr := fmt.Sprintf("%.3f+(t-%.3f)/(%.3f-%.3f)*(%.3f-%.3f)",
					prevVol, transStartSec, transEndSec, transStartSec, seg.Volume, prevVol)

				// Three zones: transition period, constant period, and rest
				if expr == "" {
					expr = fmt.Sprintf("if(between(t,%.3f,%.3f),%s,if(between(t,%.3f,%.3f),%.3f,%.3f))",
						transStartSec, transEndSec, rampExpr,
						transEndSec, endSec, seg.Volume,
						defaultVolume)
				} else {
					expr = fmt.Sprintf("if(between(t,%.3f,%.3f),%s,if(between(t,%.3f,%.3f),%.3f,%s))",
						transStartSec, transEndSec, rampExpr,
						transEndSec, endSec, seg.Volume,
						expr)
				}
			} else {
				// Transition duration is 0 or invalid, treat as instant
				if expr == "" {
					expr = fmt.Sprintf("if(between(t,%.3f,%.3f),%.3f,%.3f)",
						startSec, endSec, seg.Volume, defaultVolume)
				} else {
					expr = fmt.Sprintf("if(between(t,%.3f,%.3f),%.3f,%s)",
						startSec, endSec, seg.Volume, expr)
				}
			}
		} else {
			// Instant transition
			if expr == "" {
				expr = fmt.Sprintf("if(between(t,%.3f,%.3f),%.3f,%.3f)",
					startSec, endSec, seg.Volume, defaultVolume)
			} else {
				expr = fmt.Sprintf("if(between(t,%.3f,%.3f),%.3f,%s)",
					startSec, endSec, seg.Volume, expr)
			}
		}
	}

	return expr
}

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
			fmt.Printf("DEBUG: Adding intro: %s\n", *v.introPath)
		}
		videoPaths = append(videoPaths, v.filepath)
		fmt.Printf("DEBUG: Adding main: %s\n", v.filepath)
		if hasOutro {
			videoPaths = append(videoPaths, *v.outroPath)
			fmt.Printf("DEBUG: Adding outro: %s\n", *v.outroPath)
		}

		fmt.Printf("DEBUG: Total video paths: %d\n", len(videoPaths))

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
			// Determine which index is the main video
			mainVideoIndex := 0
			if hasIntro {
				mainVideoIndex = 1 // Main video is after intro
			}

			// Scale all videos to the main video's resolution
			for i := 0; i < len(videoPaths); i++ {
				videoFilter := fmt.Sprintf("[%d:v]", i)
				audioFilter := fmt.Sprintf("[%d:a:0]", i)

				// Apply trim to main video only if start/end are set
				if i == mainVideoIndex && (v.start > 0 || v.end < v.duration) {
					// Trim video stream
					trimDuration := (v.end - v.start).Seconds()
					videoFilter += fmt.Sprintf("trim=start=%.3f:duration=%.3f,setpts=PTS-STARTPTS,", v.start.Seconds(), trimDuration)
					// Trim audio stream
					audioFilter += fmt.Sprintf("atrim=start=%.3f:duration=%.3f,asetpts=PTS-STARTPTS,", v.start.Seconds(), trimDuration)
				}

				// Apply scaling
				videoFilter += fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease,setsar=1,pad=%d:%d:(ow-iw)/2:(oh-ih)/2[v%d];",
					v.width, v.height, v.width, v.height, i)
				audioFilter += fmt.Sprintf("[a%d];", i)

				filterComplex += videoFilter + audioFilter
			}

			// Build concat part with video and audio streams
			for i := 0; i < len(videoPaths); i++ {
				filterComplex += fmt.Sprintf("[v%d][a%d]", i, i)
			}

			// Add concat filter
			filterComplex += fmt.Sprintf("concat=n=%d:v=1:a=1[outv][videoa];", len(videoPaths))
			videoOutputLabel = "[outv]"
			audioOutputLabel = "[videoa]"
		} else {
			// Single video, use direct stream mapping for video (no filter needed)
			videoOutputLabel = "0:v"
			audioOutputLabel = "[0:a:0]"
		}

		// Build background music filter chain if needed
		if hasMusic {
			music := v.backgroundMusic

			// Calculate total video duration including intro and outro
			// Use trimmed duration for main video (v.end - v.start)
			mainVideoDuration := v.end - v.start
			totalDuration := mainVideoDuration

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

			// Reset timestamps after looping so volume filter time references are correct
			// This is critical for time-based volume control to work properly
			musicFilter += "asetpts=PTS-STARTPTS,"

			// Apply volume adjustment (dynamic or static)
			if len(music.VolumeSegments) > 0 {
				// Build dynamic volume expression from segments
				fmt.Printf("DEBUG: Building volume expression from %d segments\n", len(music.VolumeSegments))
				for i, seg := range music.VolumeSegments {
					fmt.Printf("  Segment %d: %.1fs-%.1fs, vol=%.3f, type=%s\n",
						i, seg.StartTime.Seconds(), seg.EndTime.Seconds(), seg.Volume, seg.TransitionType)
				}
				volumeExpr := buildVolumeExpression(music.VolumeSegments)
				fmt.Printf("DEBUG: Volume expression: %s\n", volumeExpr)
				// Use volume filter with eval=frame to evaluate expression for each frame
				musicFilter += fmt.Sprintf("volume='%s':eval=frame,", volumeExpr)
			} else {
				// Use static volume
				fmt.Printf("DEBUG: Using static volume: %.3f\n", music.Volume)
				musicFilter += fmt.Sprintf("volume=%.3f,", music.Volume)
			}

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

			fmt.Printf("DEBUG: Full music filter: %s\n", musicFilter)
			filterComplex += musicFilter

			// Mix music with video audio
			filterComplex += fmt.Sprintf("%s[music]amix=inputs=2:duration=first:dropout_transition=2[outa]", audioOutputLabel)
			audioOutputLabel = "[outa]"
		}

		// Add filter_complex if we built one
		if filterComplex != "" {
			fmt.Printf("DEBUG: Final filter_complex:\n%s\n", filterComplex)
			cmdline = append(cmdline, "-filter_complex", filterComplex)
		}

		// Map outputs
		fmt.Printf("DEBUG: Video output: %s, Audio output: %s\n", videoOutputLabel, audioOutputLabel)
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
