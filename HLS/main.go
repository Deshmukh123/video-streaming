package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/bluenviron/gohlslib"
)

func main() {
	log.Println("Starting HLS video streaming application...")

	// Create HLS directory if it doesn't exist
	log.Println("Creating HLS directory...")
	err := os.MkdirAll("hls", 0755)
	if err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}
	log.Println("HLS directory created successfully.")

	// Set up HLS muxer
	log.Println("Initializing HLS muxer...")
	muxer := &gohlslib.Muxer{
		SegmentCount:    7,               // At least 7 segments required for Low-Latency HLS
		SegmentDuration: 4 * time.Second, // Duration of each segment
		Directory:       "hls",           // Directory where the HLS files will be saved
	}
	log.Println("HLS muxer initialized.")

	// Start the muxer and handle errors
	log.Println("Starting HLS muxer...")
	err = muxer.Start()
	if err != nil {
		log.Fatalf("Failed to start muxer: %v", err)
	}
	log.Println("HLS muxer started successfully.")

	// Serve the HLS files
	log.Println("Setting up HTTP server to serve HLS files...")
	http.Handle("/hls/", http.StripPrefix("/hls/", http.FileServer(http.Dir("hls"))))

	// Start the HTTP server
	go func() {
		log.Println("Starting HTTP server on http://localhost:8080/hls/index.m3u8...")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Simulate the video packet data (use real video data in production)
	log.Println("Preparing simulated video data...")
	videoData := [][]byte{
		{0x00, 0x00, 0x01, 0x09, 0x10}, // Example video data (NAL unit)
	}
	log.Println("Simulated video data ready.")

	// Use a WaitGroup to ensure the goroutine starts after muxer is initialized
	var wg sync.WaitGroup
	wg.Add(1)

	// Write video data packets to the muxer with timestamps
	go func() {
		log.Println("Waiting for muxer initialization to complete...")
		wg.Wait() // Wait for the muxer to be initialized
		log.Println("Muxer initialization confirmed. Starting to write video data...")

		for {
			log.Println("Writing video data to muxer...")
			err := muxer.WriteH26x(time.Now(), 4*time.Second, videoData) // Pass timestamp, duration, and NAL units
			if err != nil {
				log.Printf("Write error: %v\n", err)
			} else {
				log.Println("Video data written successfully.")
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Signal that the muxer is initialized
	log.Println("Muxer initialization complete. Signaling goroutine to proceed...")
	wg.Done()

	log.Println("Application is running. Press Ctrl+C to exit.")
	select {} // Keep running
}
