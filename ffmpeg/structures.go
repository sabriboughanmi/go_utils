package ffmpeg

import "time"

type VideoResolution int

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
	filepath       string
	width          int
	height         int
	fps            int
	bitrate        int
	rotate         *int
	start          time.Duration
	end            time.Duration
	duration       time.Duration
	filters        []string
	additionalArgs []string
}
