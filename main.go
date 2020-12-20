package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gocv.io/x/gocv"
)

func main() {
	fmt.Printf("gocv version: %s\n", gocv.Version())
	fmt.Printf("opencv lib version: %s\n", gocv.OpenCVVersion())
	InitConfig()
	// parse args
	deviceID := config.Camera.DeviceID

	fmt.Printf("deviceID = %v\n", deviceID)

	var window *gocv.Window
	//Create Windows
	if config.Window.Enable {
		window = gocv.NewWindow(config.Window.Title)
		window.ResizeWindow(config.Window.Width, config.Window.Height)
		window.SetWindowProperty(gocv.WindowPropertyOpenGL, gocv.WindowFlag(0x00001000))
	}
	ctx, cancel := context.WithCancel(context.Background())
	homeRec, _ := NewHomeRecognizer()
	defer homeRec.Close()
	cv, err := NewCVCapture(config.Camera, window, homeRec, cancel)
	if err != nil {
		fmt.Printf("create opencv window fail: %v\n", err)
		return
	}
	cv.Run(ctx)
	defer cv.Close()
	waitForSignal(context.WithCancel(context.Background()))
}

func waitForSignal(ctx context.Context, cancelFunc context.CancelFunc) {
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals
	cancelFunc()
	//time sleep for out "fmt.Println("cancel and quit doLoop")"
	time.Sleep(1 * time.Second)
}
