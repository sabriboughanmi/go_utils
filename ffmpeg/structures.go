package ffmpeg

import "time"

type VideoResolution int

// VolumeTransitionType defines how volume transitions between segments
type VolumeTransitionType string

const (
	// VolumeTransitionInstant changes volume immediately at the boundary
	VolumeTransitionInstant VolumeTransitionType = "instant"
	// VolumeTransitionLinear gradually transitions volume over the transition duration
	VolumeTransitionLinear VolumeTransitionType = "linear"
)

// VolumeSegment defines a time-based volume control segment for background music
type VolumeSegment struct {
	StartTime          time.Duration        // Start time of this volume segment
	EndTime            time.Duration        // End time of this volume segment
	Volume             float32              // Volume level (0.0 to 1.0) for this segment
	TransitionType     VolumeTransitionType // How to transition into this volume level
	TransitionDuration time.Duration        // Duration of transition (only used if TransitionType is Linear)
}

// BackgroundMusicOptions contains all configuration options for background music
type BackgroundMusicOptions struct {
	FilePath         string          // Path to the audio file
	Volume           float32         // Volume as percentage (0.0 to 1.0) of original music volume (ignored if VolumeSegments is set)
	Loop             bool            // Loop music if shorter than video duration
	StartTime        time.Duration   // Offset in the music file to start from
	FadeInDuration   time.Duration   // Fade in duration (0 = no fade in)
	FadeOutDuration  time.Duration   // Fade out duration (0 = no fade out)
	VolumeSegments   []VolumeSegment // Dynamic volume control over time (if set, overrides Volume field)
}

// EditableVideo and Editable Video representation which  contains information about a video file and all the operations that
// need to be applied to it. Call Load to initialize a Video from file. Call the
// transformation functions to generate the desired output. Then call Render to
// generate the final output video file.
type EditableVideo Video

// Video contains information about a video file and all the operations that
// need to be applied to it. Call Load to initialize a Video from file. Call the
// transformation functions to generate the desired output. Then call Render to
// generate the final output video file.
type Video struct {
	filepath        string
	width           int
	height          int
	fps             int
	bitrate         int
	rotate          *int
	start           time.Duration
	end             time.Duration
	duration        time.Duration
	filters         []string
	additionalArgs  []string
	introPath       *string
	outroPath       *string
	backgroundMusic *BackgroundMusicOptions
}
