package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filePath)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run ffprobe: %w", err)
	}

	type ffprobeOut struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	var result ffprobeOut
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return "", fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	if len(result.Streams) == 0 {
		return "", fmt.Errorf("no video streams found")
	}

	width := result.Streams[0].Width
	height := result.Streams[0].Height

	if width == 0 || height == 0 {
		return "", fmt.Errorf("invalid dimensions: width=%d height=%d", width, height)
	}

	if width/height == 16/9 {
		return "landscape", nil
	} else if height/width == 16/9 {
		return "portrait", nil
	}

	return "other", nil
}
