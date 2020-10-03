package main

import (
	"io"
	"os"

	bpm "github.com/benjojo/bpm"
	"github.com/gin-gonic/gin"
	"github.com/go-audio/wav"
	"github.com/kkdai/youtube/v2"
	"github.com/mattetti/audio/decoder"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/catjam", func(c *gin.Context) {
		videoID := c.Query("video")
		client := youtube.Client{}
		video, err := client.GetVideo(videoID)

		if err != nil {
			c.JSON(500, gin.H{
				"message":    "internal server error",
				"stacktrace": err,
			})
		}

		resp, err := client.GetStream(video, &video.Formats[0])
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		file, err := os.Create("vid.mp3")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			panic(err)
		}
		f, err := os.Open("vid.mp3")
		if err != nil {
			c.JSON(500, gin.H{
				"message":    "internal server error",
				"stacktrace": err,
			})
		}
		defer f.Close()
		var dec decoder.Decoder
		d := wav.NewDecoder(f)
		if d.IsValidFile() {
			dec = d
		}
		buffer, err := dec.FullPCMBuffer()
		floatbuffer := buffer.AsFloat32Buffer().Data
		floatarray := bpm.ReadFloatArray(floatbuffer)
		res := bpm.ScanForBpm(floatarray, 60, 180, 10, 1)

		c.JSON(200, gin.H{
			"message": "success",
			"bpm":     res,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
