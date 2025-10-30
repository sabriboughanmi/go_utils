package ffmpeg

import (
	"testing"
	"time"
)

// ============================================
// Configuration - Update these paths and settings before running
// ============================================
var (
	// Input video paths
	mainVideoPath  = "/Users/sabriboughanmi/Documents/Smash/mainvideo3.mp4"
	introVideoPath = "/Users/sabriboughanmi/Documents/Smash/Intro.mp4"
	outroVideoPath = "/Users/sabriboughanmi/Documents/Smash/Outro.mp4"
	musicFilePath  = "/Users/sabriboughanmi/Documents/Smash/Musique.MP3"

	// Output path
	outputVideoPath = "/Users/sabriboughanmi/Documents/Smash/output_complete.mp4"

	// Video settings
	targetResolution = VideoResolution(720)
	preset           = Medium
	crf              = 23

	// Background music settings
	musicVolume          = float32(0.5) // 25% volume (quiet background)
	musicLoop            = false
	musicStartTime       = 0 * time.Second
	musicFadeInDuration  = 0 * time.Second
	musicFadeOutDuration = 0 * time.Second
)

// TestCompleteWorkflow tests using all features together
func TestCompleteWorkflow(t *testing.T) {
	// Load the main video
	video, err := LoadVideo(mainVideoPath)
	if err != nil {
		t.Fatalf("Failed to load main video: %v", err)
	}

	// Get editable version
	editable := video.GetEditableVideo()

	// Add intro and outro
	editable.SetIntro(introVideoPath)
	editable.SetOutro(outroVideoPath)

	// Add background music with all options
	editable.SetBackgroundMusic(BackgroundMusicOptions{
		FilePath:        musicFilePath,
		Volume:          musicVolume,
		Loop:            musicLoop,
		StartTime:       musicStartTime,
		FadeInDuration:  musicFadeInDuration,
		FadeOutDuration: musicFadeOutDuration,
	})

	// Apply additional transformations
	editable.SetResolution(targetResolution)
	editable.SetPreset(preset)
	editable.SetConstantRateFactor(crf)
	editable.SetStreamable()

	// Render final output
	err = editable.Render(outputVideoPath)
	if err != nil {
		t.Fatalf("Failed to render complete video: %v", err)
	}

	t.Log("Successfully created complete video with intro, outro, background music, and transformations")
}
