package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func SplitVideo(inputVideo, outputDir string, fragmentDuration int) ([]string, error) {
	var fragments []string

	// Ensure the output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating output directory: %v", err)
	}

	// FFmpeg command to split the video
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputVideo, // Input video
		"-c", "copy", // Copy codec (no re-encoding)
		"-map", "0", // Map all streams
		"-f", "segment", // Use segment muxer
		"-segment_time", strconv.Itoa(fragmentDuration), // Fragment duration
		"-reset_timestamps", "1", // Reset timestamps for each fragment
		filepath.Join(outputDir, "fragment_%03d.mp4"), // Output fragment names
	)

	// Run the FFmpeg command
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error splitting video: %v", err)
	}

	// Collect the fragment paths
	matches, _ := filepath.Glob(filepath.Join(outputDir, "fragment_*.mp4"))
	fragments = append(fragments, matches...)

	return fragments, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: ./split_video <input_video> <output_directory> <fragment_duration>")
		return
	}

	inputVideo := os.Args[1]
	outputDir := os.Args[2]
	fragmentDuration, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid fragment duration:", err)
		return
	}

	// Split the video into fragments
	fragments, err := SplitVideo(inputVideo, outputDir, fragmentDuration)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Video split into fragments:")
	for _, fragment := range fragments {
		fmt.Println(fragment)
	}
}
