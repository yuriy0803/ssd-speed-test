package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	testFile     = "testfile.tmp"
	testSizeMB   = 5120             // Test file size in MB (5GB)
	bufferSize   = 4 * 1024 * 1024  // 4MB buffer
	testDuration = 60 * time.Second // Test duration of 1 minute
)

func main() {
	fmt.Println("Starting SSD speed test...")
	endTime := time.Now().Add(testDuration)
	var totalWriteSpeed, totalReadSpeed float64
	var iterations int

	for time.Now().Before(endTime) {
		writeSpeed := measureWriteSpeed()
		readSpeed := measureReadSpeed()

		totalWriteSpeed += writeSpeed
		totalReadSpeed += readSpeed
		iterations++

		fmt.Printf("\rWrite Speed: %.2f MB/s | Read Speed: %.2f MB/s", writeSpeed, readSpeed)
		time.Sleep(1 * time.Second)
	}

	avgWriteSpeed := totalWriteSpeed / float64(iterations)
	avgReadSpeed := totalReadSpeed / float64(iterations)

	fmt.Printf("\nAverage Write Speed: %.2f MB/s | Average Read Speed: %.2f MB/s\n", avgWriteSpeed, avgReadSpeed)
	fmt.Println("SSD speed test completed.")

	// Remove test file after benchmark
	if err := os.Remove(testFile); err != nil {
		fmt.Printf("Warning: Failed to delete test file: %v\n", err)
	}
}

func measureWriteSpeed() float64 {
	data := make([]byte, bufferSize)
	for i := range data {
		data[i] = 0xFF
	}

	file, err := os.Create(testFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return 0
	}
	defer file.Close()

	start := time.Now()
	written := 0
	for time.Since(start) < time.Second {
		if _, err := file.Write(data); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return 0
		}
		written += bufferSize
	}
	duration := time.Since(start).Seconds()
	return float64(written) / (1024 * 1024) / duration
}

func measureReadSpeed() float64 {
	file, err := os.Open(testFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return 0
	}
	defer file.Close()

	buffer := make([]byte, bufferSize)
	start := time.Now()
	read := 0
	for time.Since(start) < time.Second {
		if _, err := file.Read(buffer); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading file: %v\n", err)
			return 0
		}
		read += bufferSize
	}
	duration := time.Since(start).Seconds()
	return float64(read) / (1024 * 1024) / duration
}
