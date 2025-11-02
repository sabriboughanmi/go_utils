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
	mainVideoPath  = "/Users/sabriboughanmi/Documents/Smash/mainvideo4.mp4"
	introVideoPath = "/Users/sabriboughanmi/Documents/Smash/intro.mp4"
	outroVideoPath = "/Users/sabriboughanmi/Documents/Smash/outro.mp4"
	musicFilePath  = "/Users/sabriboughanmi/Documents/Smash/music.MP3"

	// Output path
	outputVideoPath = "/Users/sabriboughanmi/Documents/Smash/output_complete.mp4"

	// Video settings
	targetResolution = VideoResolution(720)
	preset           = Medium
	crf              = 23

	// Background music settings
	musicVolume          = float32(2) // 25% volume (quiet background)
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

// TestDynamicVolumeControl demonstrates dynamic volume control for background music
func TestDynamicVolumeControl(t *testing.T) {
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

	// Add background music
	editable.SetBackgroundMusic(BackgroundMusicOptions{
		FilePath: musicFilePath,
		Loop:     true,
	})

	// Set dynamic volume control using convenience methods
	// High volume for intro
	editable.SetMusicVolumeForIntro(2, VolumeTransitionInstant)

	// Lower volume for main video (where people speak)
	editable.SetMusicVolumeForMain(0.6, VolumeTransitionLinear)

	// Medium volume for outro with fade
	editable.SetMusicVolumeForOutro(2, VolumeTransitionInstant)

	// Apply transformations
	editable.SetResolution(targetResolution)
	editable.SetPreset(preset)
	editable.SetConstantRateFactor(crf)
	editable.SetStreamable()

	// Render output
	outputPath := "/Users/sabriboughanmi/Documents/Smash/output_dynamic_volume.mp4"
	t.Logf("Rendering with intro: %v, outro: %v", editable.IntroPath() != nil, editable.OutroPath() != nil)
	if editable.BackgroundMusic() != nil {
		t.Logf("Background music segments: %d", len(editable.BackgroundMusic().VolumeSegments))
	}

	err = editable.Render(outputPath)
	if err != nil {
		t.Fatalf("Failed to render video with dynamic volume: %v", err)
	}

	t.Log("Successfully created video with dynamic volume control")
}

// TestCustomVolumeSegments demonstrates using custom volume segments
func TestCustomVolumeSegments(t *testing.T) {
	video, err := LoadVideo(mainVideoPath)
	if err != nil {
		t.Fatalf("Failed to load main video: %v", err)
	}

	editable := video.GetEditableVideo()

	// Add background music
	editable.SetBackgroundMusic(BackgroundMusicOptions{
		FilePath: musicFilePath,
		Loop:     true,
	})

	// Add custom volume segments at specific times
	// High volume for first 10 seconds
	editable.AddMusicVolumeSegment(VolumeSegment{
		StartTime:          0,
		EndTime:            10 * time.Second,
		Volume:             0.8,
		TransitionType:     VolumeTransitionInstant,
		TransitionDuration: 0,
	})

	// Lower volume from 10-30 seconds (with 2 second smooth transition)
	editable.AddMusicVolumeSegment(VolumeSegment{
		StartTime:          10 * time.Second,
		EndTime:            30 * time.Second,
		Volume:             0.2,
		TransitionType:     VolumeTransitionLinear,
		TransitionDuration: 2 * time.Second,
	})

	// Back to high volume for the rest (with 1 second transition)
	editable.AddMusicVolumeSegment(VolumeSegment{
		StartTime:          30 * time.Second,
		EndTime:            60 * time.Second,
		Volume:             0.7,
		TransitionType:     VolumeTransitionLinear,
		TransitionDuration: 1 * time.Second,
	})

	// Render
	outputPath := "/Users/sabriboughanmi/Documents/Smash/output_custom_segments.mp4"
	err = editable.Render(outputPath)
	if err != nil {
		t.Fatalf("Failed to render video with custom segments: %v", err)
	}

	t.Log("Successfully created video with custom volume segments")
}
