package ffmpeg

import (
	"bytes"
	"fmt"
	osUtils "github.com/sabriboughanmi/go_utils/os"
	"os"
	"os/exec"
)


//LoadVideoFromReEncodedFragmentsIgnoreRotation returns a merged Video that can be operated on.
//Note! path and Fragments need to be already Existing.
//Note! this function will ReEncode all videos to fit the lowest resolution and Rotation will be ignored.
func LoadVideoFromReEncodedFragmentsIgnoreRotation(path string, fragmentsPath ...string) (*Video, error) {

	//////////////////////////////////////Rotate the Fragment \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	if len(fragmentsPath) < 2 {
		return nil, fmt.Errorf("at least 2 fragments must be passed")
	}

	var noRotationProvided = false

	var rotations = make([]int, len(fragmentsPath))
	for i, path := range fragmentsPath {
		video, err := LoadVideo(path)
		if err != nil {
			return nil, err
		}
		if video.GetVideoRotate() == nil {
			noRotationProvided = true
			break
		}
		rotations[i] = *video.GetVideoRotate()
	}

	//Videos do not require rotation
	if noRotationProvided == true {
		return LoadVideoFromReEncodedFragments(path, fragmentsPath...)
	}

	var finalTsFiles []string

	var filesToDelete []string
	//ffmpeg -i 0.mp4 -c copy -metadata:s:v:0 rotate=0 md_0.mp4
	for i, fragmentPath := range fragmentsPath {

		originalFragmentPath, err := osUtils.CreateTempFile(fmt.Sprintf("_md_%d.mp4", i), nil)
		if err != nil {
			return nil, err
		}
		filesToDelete = append(filesToDelete, originalFragmentPath)

		//Construct Command to Rotate the Fragment
		cmdline := []string{
			"ffmpeg",
			"-y",
			"-i",
			fragmentPath,
			"-c",
			"copy",
			"-metadata:s:v:0",
			fmt.Sprintf("rotate=%d", rotations[i]-90),
			originalFragmentPath,
		}

		//Execute Command
		cmd := exec.Command(cmdline[0], cmdline[1:]...)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		cmd.Stdout = nil

		if err = cmd.Run(); err != nil {
			//Remove files
			defer osUtils.RemovePathsIfExists(filesToDelete...) //Ignore highlighted Leak as the look will break here with the return statement
			defer osUtils.RemovePathsIfExists(finalTsFiles...)  //Ignore highlighted Leak as the look will break here with the return statement
			return nil, fmt.Errorf(stderr.String())
		}

		//////////////////////////////////////Transpose the Fragment \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

		transposedFragmentPath, err := osUtils.CreateTempFile(fmt.Sprintf("_rt_%d.mp4", i), nil)
		if err != nil {
			return nil, err
		}
		filesToDelete = append(filesToDelete, transposedFragmentPath)

		//Construct Command to Rotate the Fragment
		cmdline = []string{
			"ffmpeg",
			"-y",
			"-i",
			originalFragmentPath,
			"-vf",
			"transpose=1",
			transposedFragmentPath,
		}
		//Execute Command
		cmd = exec.Command(cmdline[0], cmdline[1:]...)
		cmd.Stderr = &stderr
		cmd.Stdout = nil

		if err = cmd.Run(); err != nil {
			//Remove files
			defer osUtils.RemovePathsIfExists(filesToDelete...) //Ignore highlighted Leak as the look will break here with the return statement
			defer osUtils.RemovePathsIfExists(finalTsFiles...)  //Ignore highlighted Leak as the look will break here with the return statement
			return nil, fmt.Errorf(stderr.String())
		}

		///////////////////////////////////Get Fragment TS\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

		TsFragmentPath, err := osUtils.CreateTempFile(fmt.Sprintf("_%d.ts", i), nil)
		if err != nil {
			return nil, err
		}
		finalTsFiles = append(finalTsFiles, TsFragmentPath)
		// ffmpeg -i rt_0.mp4 -c copy -bsf:v h264_mp4toannexb -f mpegts 0.ts
		//Construct Command to Rotate the Fragment
		cmdline = []string{
			"ffmpeg",
			"-y",
			"-i",
			transposedFragmentPath,
			"-c",
			"copy",
			"-bsf:v",
			"h264_mp4toannexb",
			"-f",
			"mpegts",
			TsFragmentPath,
		}
		//Execute Command
		cmd = exec.Command(cmdline[0], cmdline[1:]...)
		cmd.Stderr = &stderr
		cmd.Stdout = nil

		if err = cmd.Run(); err != nil {
			//Remove files
			defer osUtils.RemovePathsIfExists(filesToDelete...) //Ignore highlighted Leak as the look will break here with the return statement
			defer osUtils.RemovePathsIfExists(finalTsFiles...)  //Ignore highlighted Leak as the look will break here with the return statement
			return nil, fmt.Errorf(stderr.String())
		}

	}

	defer osUtils.RemovePathsIfExists(filesToDelete...)
	defer osUtils.RemovePathsIfExists(finalTsFiles...)

	var concatinatedPaths = "concat:"

	for i, tsPath := range finalTsFiles {
		if i < len(finalTsFiles)-1 {
			concatinatedPaths += tsPath + "|"
		} else {
			concatinatedPaths += tsPath
		}

	}

	//Construct Command to Rotate the Fragment
	cmdline := []string{
		"ffmpeg",
		"-y",
		"-i",
		concatinatedPaths,
		"-c",
		"copy",
		"-bsf:a",
		"aac_adtstoasc",
		path,
	}

	//Execute Command
	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	var stderr2 bytes.Buffer
	cmd.Stderr = &stderr2
	cmd.Stdout = nil

	if err := cmd.Run(); err != nil {
		//Remove files
		defer osUtils.RemovePathsIfExists(filesToDelete...) //Ignore highlighted Leak as the look will break here with the return statement
		defer osUtils.RemovePathsIfExists(finalTsFiles...)  //Ignore highlighted Leak as the look will break here with the return statement
		return nil, fmt.Errorf(stderr2.String())
	}

	return LoadVideo(path)
}


//LoadVideoFromReEncodedFragments returns a merged Video that can be operated on.
//Note! path and Fragments need to be already Existing.
//Note! this function will ReEncode all videos to fit the lowest resolution.
func LoadVideoFromReEncodedFragments(path string, fragmentsPath ...string) (*Video, error) {

	if len(fragmentsPath) < 2 {
		return nil, fmt.Errorf("at least 2 fragments must be passed")
	}

	importList := `# Fragments Paths`
	for _, p := range fragmentsPath {
		importList += fmt.Sprintf("\nfile '%s'", p)
	}

	listPath, err := osUtils.CreateTempFile("list.txt", []byte(importList))
	if err != nil {
		return nil, err
	}
	defer os.Remove(listPath)

	cmdline := []string{
		"ffmpeg",
		"-y",
	}

	//get the lowest video resolution
	var lowestRes = VideoResolution(8000)

	//Add Inputs
	for i := 0; i < len(fragmentsPath); i++ {
		cmdline = append(cmdline, "-i", fragmentsPath[i])
		video, err := LoadVideo(fragmentsPath[i])
		if err != nil {
			return nil, err
		}

		//Get the Lowest Resolution available, otherwise ffmpeg will throw an error while merging the videos
		currentVideoRes := video.GetVideoResolution()
		if lowestRes > currentVideoRes {
			lowestRes = currentVideoRes
		}
	}

	//Construct Fragments Resolutions
	var filterComplex = ""
	for i := 0; i < len(fragmentsPath); i++ {
		filterComplex += fmt.Sprintf("[%d]scale=%d:-2:force_original_aspect_ratio=decrease,setsar=1[v%d];", i, lowestRes, i)
	}

	for i := 0; i < len(fragmentsPath); i++ {
		filterComplex += fmt.Sprintf("[v%d][%d:a:0]", i, i)
	}

	filterComplex += fmt.Sprintf("concat=n=%d:v=1:a=1[v][a]", len(fragmentsPath))

	//Add -filter_complex
	cmdline = append(cmdline, "-filter_complex", filterComplex)

	// Add the Output
	cmdline = append(cmdline, "-map", "[v]", "-map", "[a]", path)

	//fmt.Println(cmdline)

	cmd := exec.Command(cmdline[0], cmdline[1:]...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = nil

	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf(stderr.String())
	}
	return LoadVideo(path)
}
