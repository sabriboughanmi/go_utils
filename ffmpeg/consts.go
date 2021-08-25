package ffmpeg


// ConversionPreset .
type ConversionPreset string

const (
	Ultrafast ConversionPreset = "ultrafast"
	Superfast ConversionPreset = "superfast"
	Veryfast  ConversionPreset = "veryfast"
	Faster    ConversionPreset = "faster"
	Fast      ConversionPreset = "fast"
	Medium    ConversionPreset = "medium"
	Slow      ConversionPreset = "slow"
	Slower    ConversionPreset = "slower"
	Veryslow  ConversionPreset = "veryslow"
	Placebo   ConversionPreset = "placebo"
)


// ScreenOrientation defines if a video is Portrait or Landscape
type ScreenOrientation byte

const (
	Portrait  ScreenOrientation = 0
	Landscape ScreenOrientation = 1
)

