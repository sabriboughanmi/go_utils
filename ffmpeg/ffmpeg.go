package ffmpeg

import (
	"bytes"
	"cloud.google.com/go/storage"
	Vision "cloud.google.com/go/vision/apiv1"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	storageUtils "github.com/sabriboughanmi/go_utils/firebase/storage"
	osUtils "github.com/sabriboughanmi/go_utils/os"
	"google.golang.org/protobuf/reflect/protoreflect"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

// GetVideoOrientation returns the video Screen Orientation
func (v *Video) GetVideoOrientation() ScreenOrientation {
	if v.width > v.height {
		return Landscape
	}
	return Portrait
}

// GetVideoRotate returns the video Screen Orientation
func (v *Video) GetVideoRotate() *int {
	return v.rotate
}

//GetEditableVideoResolution returns the lowest value between  width and height
func (v *EditableVideo) GetEditableVideoResolution() VideoResolution {
	if v.width > v.height {
		return VideoResolution(v.height)
	}
	return VideoResolution(v.width)
}

//GetVideoResolution returns the lowest value between  width and height
func (v *Video) GetVideoResolution() VideoResolution {
	if v.width > v.height {
		return VideoResolution(v.height)
	}
	return VideoResolution(v.width)
}

//GetAspectRatio returns the Aspect Ratio
func (v *EditableVideo) GetAspectRatio() float32 {
	if v.width > v.height {
		return float32(v.width) / float32(v.height)
	}
	return float32(v.height) / float32(v.width)
}

//GetDuration return Video Duration in Seconds
func (v *Video) GetDuration() float64 {
	return v.duration.Seconds()
}

// GetThumbnailAtSec Creates a Thumbnail at path for a given time
func (v *Video) GetThumbnailAtSec(outputPath string, second float64) error {

	cmds := []string{
		"ffmpeg",
		"-y",
		"-i", v.filepath,
		"-vframes", "1", "-an",
		"-s", fmt.Sprintf("%dx%d", v.width, v.height),
		"-ss", strconv.FormatFloat(second, 'f', -1, 64),
		outputPath,
	}

	cmd := exec.Command(cmds[0], cmds[1:]...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = nil

	err := cmd.Run()
	if err != nil {
		return errors.New("Video.Render: ffmpeg failed: " + stderr.String())
	}
	return nil
}

// ModerateVideo verify if a video contain forbidden content
func (v *Video) ModerateVideo(sequenceDuration float64, ctx context.Context, tolerance int32, tempStorageObject *temporaryStorageObjectRef  ) error {
	errorChannel := make(chan error)
	var duration, moderateDuration float64
	moderateDuration = v.GetDuration()
	var wg sync.WaitGroup
	duration = 0


	for {
		wg.Add(1)

		//moderate frame every x sec
		go func(waitGroup *sync.WaitGroup, errorChan chan error) {
			defer wg.Done()
			//Done: path = creation temp file
			path, err := osUtils.CreateTempFile("pic.png", nil)
			if err != nil {
				errorChan <- fmt.Errorf("CreateTempFile: , Error:  %v", err)
				return
			}
			defer os.Remove(path)

			if err := v.GetThumbnailAtSec(path, duration); err != nil {
				errorChan <- fmt.Errorf("GetThumbnailAtSec %f : , Error:  %v", duration, err)
			}
			/*
			if _, err := ModerateVideoFrame(path, ctx, tolerance, client, tempStorageObject); err != nil {
				errorChan <- err
				return
			}*/

		}(&wg, errorChannel)

		duration += sequenceDuration

		if duration >= moderateDuration {
			break
		}
	}

	wg.Wait()

	return nil
}

type temporaryStorageObjectRef struct {
	Client *storage.Client
	Bucket string
}

//GetTemporaryStorageObjectRef is used to send necessary data to ModerateVideoFrame
func GetTemporaryStorageObjectRef(client *storage.Client, bucket string) temporaryStorageObjectRef {
	return temporaryStorageObjectRef{
		Client: client,
		Bucket: bucket,
	}
}

// ModerateVideoFrame verify if an extended frame contain forbidden content
func ModerateVideoFrame(localPath string, ctx context.Context, tolerance int32, client *Vision.ImageAnnotatorClient, tempStorageObject *temporaryStorageObjectRef) (bool, error) {
	storageFileURI := tempStorageObject.Bucket + localPath
	// create image in  storage
	if err := storageUtils.CreateStorageFileFromLocal(tempStorageObject.Bucket, localPath, localPath, nil, tempStorageObject.Client, ctx); err != nil {
		return false, fmt.Errorf("CreateStorageFileFromLocal : , Error:  %v", err)
	}
	// remove image
	defer func() {
		if err := storageUtils.RemoveFile(tempStorageObject.Bucket, localPath, tempStorageObject.Client, ctx); err != nil {
			fmt.Printf("ModerateVideoFrame : Error deleting Temp file %v \n", err)
			return
		}
	}()

	image := Vision.NewImageFromURI(storageFileURI)

	props, err := client.DetectSafeSearch(ctx, image, nil)
	if err != nil {
		return false, fmt.Errorf("DetectSafeSearch : , Error:  %v", err)
	}
	var tolr = protoreflect.EnumNumber(tolerance)

	if props.Adult.Number() > tolr || props.Violence.Number() > tolr {
		return false, errors.New("frame contain forbidden content")
	}

	return true, nil
}

// LoadVideo gives you a Video that can be operated on. Load does not open the file
// or load it into memory. Apply operations to the Video and call Render to
// generate the output video file.
func LoadVideo(path string) (*Video, error) {
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return nil, errors.New("cinema.Load: ffprobe was not found in your PATH " +
			"environment variable, make sure to install ffmpeg " +
			"(https://ffmpeg.org/) and add ffmpeg, ffplay and ffprobe to your " +
			"PATH")
	}

	if _, err := os.Stat(path); err != nil {
		return nil, errors.New("cinema.Load: unable to load file: " + err.Error())
	}

	cmdArgs := []string{"ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		path}

	fmt.Println(cmdArgs)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = nil
	out, err := cmd.Output()

	if err != nil {
		return nil, errors.New("Load: ffprobe failed with Error: " + stderr.String() + stdout.String())
	}

	type description struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
			Tags   struct {
				// Rotation is optional -> use a pointer.
				Rotation *json.Number `json:"rotate"`
			} `json:"tags"`
		} `json:"streams"`
		Format struct {
			DurationSec json.Number `json:"duration"`
			Bitrate     json.Number `json:"bit_rate"`
		} `json:"format"`
	}
	var desc description
	if err := json.Unmarshal(out, &desc); err != nil {
		return nil, errors.New("cinema.Load: unable to parse JSON output " +
			"from ffprobe: " + err.Error())
	}
	if len(desc.Streams) == 0 {
		return nil, errors.New("cinema.Load: ffprobe does not contain stream " +
			"data, make sure the file " + path + " contains a valid video.")
	}

	secs, err := desc.Format.DurationSec.Float64()
	if err != nil {
		return nil, errors.New("cinema.Load: ffprobe returned invalid duration: " +
			err.Error())
	}
	bitrate, err := desc.Format.Bitrate.Int64()
	if err != nil {
		return nil, errors.New("cinema.Load: ffprobe returned invalid duration: " +
			err.Error())
	}

	// Round seconds (floating point value) up to time.Duration. seconds will
	// be >= 0 so adding 0.5 rounds to the right integer Duration value.
	duration := time.Duration(secs*float64(time.Second) + 0.5)

	dsIndex := 0
	for index, v := range desc.Streams {
		if v.Width != 0 && v.Height != 0 {
			dsIndex = index
			break
		}
	}

	width := desc.Streams[dsIndex].Width
	height := desc.Streams[dsIndex].Height
	if desc.Streams[dsIndex].Tags.Rotation != nil {
		// If the video is rotated by -270, -90, 90 or 270 degrees, we need to
		// flip the width and height because they will be reported in unrotated
		// coordinates while cropping etc. works on the rotated dimensions.
		rotation, err := desc.Streams[dsIndex].Tags.Rotation.Int64()
		if err != nil {
			return nil, errors.New("cinema.Load: ffprobe returned invalid " +
				"rotation: " + err.Error())
		}
		flipCount := rotation / 90
		if flipCount%2 != 0 {
			width, height = height, width
		}
	}

	return &Video{
		filepath: path,
		width:    width,
		height:   height,
		fps:      30,
		bitrate:  int(bitrate),
		rotate: func() *int {
			if desc.Streams[dsIndex].Tags.Rotation == nil {
				return nil
			}
			rotation, err := desc.Streams[dsIndex].Tags.Rotation.Int64()
			if err != nil {
				return nil
			}
			var rotationInt = int(rotation)
			return &rotationInt
		}(),
		start:    0,
		end:      duration,
		duration: duration,
	}, nil
}

// LoadVideoFromFragments returns a Video that can be operated on. Load does not open the file or load it into memory.
//Note! path and Fragments need to be already Existing
func LoadVideoFromFragments(path string, fragmentsPath ...string) (*Video, error) {

	importList := `# this is a comment`
	for _, p := range fragmentsPath {
		importList += fmt.Sprintf("\nfile '%s'", p)
	}

	listPath, err := osUtils.CreateTempFile("list.txt", []byte(importList))
	if err != nil {
		return nil, err
	}
	defer os.Remove(listPath)

	//fmt.Printf("listPath: %s\n",listPath)

	cmdline := []string{
		"ffmpeg",
		"-y",
		"-f",
		"concat",
		"-safe",
		"0",
		"-i",
		listPath,
		"-c",
		"copy",
		path,
	}

	cmd := exec.Command(cmdline[0], cmdline[1:]...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = nil

	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf(stderr.String())
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

//GetEditableVideo returns an EditableVideo instance than can be used to safely modify a Video
func (v *Video) GetEditableVideo() *EditableVideo {
	var eVideo = EditableVideo(*v)

	eVideo.filters = make([]string, len(v.filters))
	copy(eVideo.filters, v.filters)

	eVideo.additionalArgs = make([]string, len(v.additionalArgs))
	copy(eVideo.additionalArgs, v.additionalArgs)

	return &eVideo
}

//AddWaterMark Adds a Water mark to a video
func (v *EditableVideo) AddWaterMark(videoPath, iconPath, outputPath string, widthSize, heightSize int) error {

	cmdline := []string{
		"ffmpeg",
		"-y",
		"-i", videoPath,
		"-i", iconPath,
		"-vcodec", "libx264",
	}
	cmdline = append(cmdline, v.additionalArgs...)
	cmdline = append(cmdline, "-filter_complex")
	cmdline = append(cmdline, fmt.Sprintf("[1]scale=%d:%d[wm];[0][wm]overlay=10:10", widthSize, heightSize))
	cmdline = append(cmdline, outputPath)

	//fmt.Println(cmdline)
	cmd := exec.Command(cmdline[0], cmdline[1:]...)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr
	cmd.Stdout = nil

	err := cmd.Run()
	if err != nil {
		return errors.New("Video.Render: ffmpeg failed: " + stderr.String())
	}
	return nil
}

// ConvertFromTo Converts any media file type to another
func ConvertFromTo(inputPath, outputPath string) error {
	cmds := []string{
		"ffmpeg",
		"-y",
		"-i", inputPath,
		outputPath,
	}

	cmd := exec.Command(cmds[0], cmds[1:]...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = nil

	err := cmd.Run()
	if err != nil {
		return errors.New("Video.Render: ffmpeg failed: " + stderr.String())
	}
	return nil
}

// GetThumbnail Creates a Thumbnail at path for a given time
func (v *EditableVideo) GetThumbnail(outputPath string, second float64) error {

	/*var width, height int
	if v.width > v.height {
		width = v.width
		height = v.height
	} else {
		width = v.height
		height = v.width
	}*/

	//ffmpeg -i InputFile.FLV -vframes 1 -an -s 400x222 -ss 30 OutputFile.jpg

	cmds := []string{
		"ffmpeg",
		"-y",
		"-i", v.filepath,
		"-vframes", "1", "-an",
		"-s", fmt.Sprintf("%dx%d", v.width, v.height),
		"-ss", strconv.FormatFloat(second, 'f', -1, 64),
		outputPath,
	}

	cmd := exec.Command(cmds[0], cmds[1:]...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = nil

	err := cmd.Run()
	if err != nil {
		return errors.New("Video.Render: ffmpeg failed: " + stderr.String())
	}
	return nil
}

// Render applies all operations to the Video and creates an output video file
// of the given name. This method won't return anything on stdout / stderr.
// If you need to read ffmpeg's outputs, use RenderWithStreams
func (v *EditableVideo) Render(output string) error {
	return v.RenderWithStreams(output, nil, nil)
}

// RenderInBackground applies all operations to the Video and creates an output video file
// of the given name. This method won't return anything on stdout / stderr.
// If you need to read ffmpeg's outputs, use RenderWithStreams
func (v *EditableVideo) RenderInBackground(output string) (*exec.Cmd, error) {
	return v.RenderWithStreamsInBackground(output, nil)
}

// RenderWithStreamsInBackground applies all operations to the Video and creates an output video file
// of the given name. By specifying an output stream and an error stream, you can read
// ffmpeg's stdout and stderr.
func (v *EditableVideo) RenderWithStreamsInBackground(output string, os io.Writer) (*exec.Cmd, error) {
	line := v.commandLine(output)
	//fmt.Println(line)

	cmd := exec.Command(line[0], line[1:]...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = os

	err := cmd.Start()
	if err != nil {
		return nil, errors.New("cinema.Video.Render: ffmpeg failed: " + stderr.String())
	}
	return cmd, nil
}

// RenderWithStreams applies all operations to the Video and creates an output video file
// of the given name. By specifying an output stream and an error stream, you can read
// ffmpeg's stdout and stderr.
func (v *EditableVideo) RenderWithStreams(output string, os io.Writer, es io.Writer) error {
	line := v.commandLine(output)
	fmt.Println(line)

	cmd := exec.Command(line[0], line[1:]...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = os

	err := cmd.Run()
	if err != nil {
		return errors.New("Video.Render: ffmpeg failed: " + stderr.String())
	}
	return nil
}

// Mute mutes the video
func (v *Video) Mute() {
	v.additionalArgs = append(v.additionalArgs, "-an")
}

// Trim sets the start and end time of the output video. It is always relative
// to the original input video. start must be less than or equal to end or
// nothing will change.
func (v *Video) Trim(start, end time.Duration) {
	if start <= end {
		v.SetStart(start)
		v.SetEnd(end)
	}
}

// Start returns the start of the video .
func (v *Video) Start() time.Duration {
	return v.start
}

// SetStart sets the start time of the output video. It is always relative to
// the original input video.
func (v *Video) SetStart(start time.Duration) {
	v.start = v.clampToDuration(start)
	if v.start > v.end {
		// keep c.start <= v.end
		v.end = v.start
	}
}

// End returns the end of the video.
func (v *Video) End() time.Duration {
	return v.end
}

// SetEnd sets the end time of the output video. It is always relative to the
// original input video.
func (v *Video) SetEnd(end time.Duration) {
	v.end = v.clampToDuration(end)
	if v.end < v.start {
		// keep c.start <= v.end
		v.start = v.end
	}
}

// SetFPS sets the framerate (frames per second) of the output video.
func (v *Video) SetFPS(fps int) {
	v.fps = fps
}

// SetBitrate sets the bitrate of the output video.
func (v *Video) SetBitrate(bitrate int) {
	v.bitrate = bitrate
}

// SetSize sets the width and height of the output video.
func (v *EditableVideo) SetSize(width int, height int) {
	v.width = width
	v.height = height
	v.additionalArgs = append(v.additionalArgs, "-s")
	v.additionalArgs = append(v.additionalArgs, fmt.Sprintf("%dx%d", width, height))
}

// SetPreset  defines the Quality Compression and Speed
func (v *EditableVideo) SetPreset(preset ConversionPreset) {
	v.additionalArgs = append(v.additionalArgs, "-preset")
	v.additionalArgs = append(v.additionalArgs, string(preset))
}

// SetConstantRateFactor The range of the CRF scale is 0–51, where 0 is lossless, 23 is the default, and 51 is worst quality possible. A lower value generally leads to higher quality, and a subjectively sane range is 17–28. Consider 17 or 18 to be visually lossless or nearly so; it should look the same or nearly the same as the input but it isn't technically lossless. The range is exponential, so increasing the CRF value +6 results in roughly half the bitrate / file size, while -6 leads to roughly twice the bitrate. Choose the highest CRF value that still provides an acceptable quality. If the output looks good, then try a higher value. If it looks bad, choose a lower value.
func (v *EditableVideo) SetConstantRateFactor(value int) {
	v.additionalArgs = append(v.additionalArgs, "-crf")
	v.additionalArgs = append(v.additionalArgs, strconv.Itoa(value))
}

//GetResolutions returns the video (Width,Height) tuple for a specific VideoResolution
func (v *EditableVideo) GetResolutions(res VideoResolution) (int, int) {
	aspectRatio := v.GetAspectRatio()
	maxSize := toEvenNumber(int(float32(res) * aspectRatio))

	if v.width > v.height {
		return maxSize, int(res)
	}

	return int(res), maxSize
}

// GetFilePath returns the path of the input video.
func (v *EditableVideo) GetFilePath() string {
	return v.filepath
}

// SetResolution sets the  Resolution respecting the Aspect Ratio of the Original Video.
func (v *EditableVideo) SetResolution(res VideoResolution) {
	aspectRatio := v.GetAspectRatio()
	maxSize := toEvenNumber(int(float32(res) * aspectRatio))

	if v.width > v.height {
		v.SetSize(maxSize, int(res))
	} else {
		v.SetSize(int(res), maxSize)
	}
}

//SetFilePath set the filepath for a video
func (v *Video) SetFilePath(p string) {
	v.filepath = p
}

// Crop makes the output video a sub-rectangle of the input video. (0,0) is the
// top-left of the video, x goes right, y goes down.
func (v *Video) Crop(x, y, width, height int) {
	v.width = width
	v.height = height
	v.filters = append(
		v.filters,
		fmt.Sprintf("crop=%d:%d:%d:%d", width, height, x, y),
	)
}

// Filepath returns the path of the input video.
func (v *Video) Filepath() string {
	return v.filepath
}
