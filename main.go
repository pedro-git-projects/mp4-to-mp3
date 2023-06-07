package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	filename := os.Args[1]

	output := fmt.Sprintf("%s.mp3", filename)
	ffmpegArgs := []string{
		"-i", fmt.Sprintf("%s", filename),
		"-vn",
		"-acodec", "libmp3lame",
		"-ac", "2",
		"-ab", "160k",
		"-ar", "48000",
		output,
	}

	// Create a pipe to capture ffmpeg output
	reader, writer := io.Pipe()

	// Create a scanner to read from the pipe
	scanner := bufio.NewScanner(reader)

	// Create a log file
	logFile, err := os.Create("ffmpeg.log")
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}
	defer logFile.Close()

	// Create a multi-writer to write to both stdout and log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Set the log output to the multi-writer
	log.SetOutput(multiWriter)

	// Execute the ffmpeg command
	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	cmd.Stdout = writer // Redirect command output to pipe
	cmd.Stderr = writer // Redirect command error output to pipe

	// Start reading from the pipe in a separate goroutine
	go func() {
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()

	err = cmd.Run()
	if err != nil {
		log.Fatal("Error executing command:", err)
	}

	log.Println("Command executed successfully.")
}
