package core

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func GetMetadata(input string) (int, int, string) {
	command := "ffprobe"
	flags := "-v error -select_streams v:0 -show_entries stream=width,height,avg_frame_rate -of csv=p=0 "
	args := strings.Fields(flags)
	args = append(args, input)

	cmd := exec.Command(command, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Erro: %v\n Command: %v %v\n PATH: %v\n", err, command, args, input)
		log.Fatal("A problem occurred when fetching the input's metadata")
	}

	metadata := strings.Split(string(output), ",")

	width, _ := strconv.Atoi(metadata[0])
	height, _ := strconv.Atoi(metadata[1])
	splitFPS := strings.Split(metadata[2], "/")
	fps := splitFPS[0]

	return width, height, fps
}
