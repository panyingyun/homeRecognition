package main

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"gocv.io/x/gocv"
)

const KeyEsc int = 27

var (
	ColorBlue = color.RGBA{B: 255}
)

type CVCapture struct {
	webcam     *gocv.VideoCapture
	window     *gocv.Window
	hr         *HomeRecognizer
	lastHrTime time.Time
	mutex      sync.Mutex
	cf         gocv.CascadeClassifier
	img        gocv.Mat
	cancel     context.CancelFunc
}

func NewCVCapture(camera CameraConfig, window *gocv.Window, hr *HomeRecognizer, cancel context.CancelFunc) (*CVCapture, error) {

	//Create Video Capture
	webcam, err := gocv.OpenVideoCapture(camera.DeviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", camera.DeviceID)
		return nil, errors.New("opening video capture device fail")
	}
	webcam.Set(gocv.VideoCaptureFPS, float64(30))
	webcam.Set(gocv.VideoCaptureFrameWidth, float64(camera.Width))
	webcam.Set(gocv.VideoCaptureFrameHeight, float64(camera.Height))
	img := gocv.NewMat()
	// create video write
	if ok := webcam.Read(&img); !ok {
		fmt.Printf("Device closed: %v\n", camera.DeviceID)
		return nil, errors.New("read image fail, device maybe closed")
	}
	fmt.Printf("Start reading device: %v\n", camera.DeviceID)
	classifier := gocv.NewCascadeClassifier()
	xmlFile := filepath.Join(dataDir, "facedetect", "haarcascade_frontalface_default.xml")
	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return nil, errors.New("classifier init error")
	}
	return &CVCapture{
		window: window,
		webcam: webcam,
		hr:     hr,
		cf:     classifier,
		img:    img,
		cancel: cancel,
	}, nil
}

func (cv *CVCapture) Close() {
	if cv.webcam != nil {
		cv.webcam.Close()
		cv.webcam = nil
	}
	cv.img.Close()
	cv.hr.Close()
}

func (cv *CVCapture) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			cv.Close()
			fmt.Println("cancel and quit CVCapture")
			return
		default:
			cv.faceDetectAndRecognition()
		}
	}
}

func (cv *CVCapture) faceDetectAndRecognition() {
	if ok := cv.webcam.Read(&cv.img); !ok {
		fmt.Printf("device closed")
		return
	}
	if cv.img.Empty() {
		return
	}
	// detect faces
	rects := cv.cf.DetectMultiScale(cv.img)
	fmt.Printf("found %d faces\n", len(rects))
	// draw a rectangle around each face on the original image,
	// along with text identifing as "Human"
	for _, r := range rects {
		gocv.Rectangle(&cv.img, r, ColorBlue, 3)

		size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
		pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
		gocv.PutText(&cv.img, "Human", pt, gocv.FontHersheyPlain, 1.2, ColorBlue, 2)
	}

	// show the image in the window, and wait 1 millisecond
	cv.window.IMShow(cv.img)
	if cv.window.WaitKey(1) == KeyEsc {
		cv.window.Close()
		cv.cancel()
	}
	go cv.recognizeFace()
}

func (cv *CVCapture) recognizeFace() {
	if time.Since(cv.lastHrTime).Seconds() < 10 {
		return
	}
	cv.mutex.Lock()
	cv.lastHrTime = time.Now()
	buf, _ := gocv.IMEncode(".jpg", cv.img)
	ret := cv.hr.RecognizeWithSingleImage(buf)
	fmt.Println("ret = ", ret)
	cv.playMp3(ret)
	cv.mutex.Unlock()
}

func (cv *CVCapture) playMp3(name string) {
	if name == "" {
		return
	}
	mp3file := filepath.Join(dataDir, "mp3", name+".mp3")
	f, err := os.Open(mp3file)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}
