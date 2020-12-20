package main

import (
	"fmt"
	"log"
	"path/filepath"

	face "github.com/Kagami/go-face"
)

const dataDir = "data"

var (
	labels = []string{"daughter", "father", "grandmother", "mother"}
)

type HomeRecognizer struct {
	rec *face.Recognizer
}

func NewHomeRecognizer() (*HomeRecognizer, error) {
	// Init the recognizer.
	recog, err := face.NewRecognizer(filepath.Join(dataDir, "models"))
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}

	// Fill known samples. In the real world you would use a lot of images
	// for each person to get better classification results
	// Name the categories, i.e. people on the image.

	var samples []face.Descriptor
	var cats []int32
	for i, label := range labels {
		image := filepath.Join(dataDir, "images", label+".jpg")
		fFace, err := recog.RecognizeSingleFile(image)
		if err != nil {
			log.Fatalf("Can't recognize: %v", err)
		}
		if fFace == nil {
			log.Fatalf("Not a single face on the image")
		}
		samples = append(samples, fFace.Descriptor)
		cats = append(cats, int32(i))
		fmt.Printf("i = %v, min = %v, max = %v\n", i, fFace.Rectangle.Min, fFace.Rectangle.Max)
	}
	// Pass samples to the recognizer.
	recog.SetSamples(samples, cats)
	return &HomeRecognizer{
		rec: recog,
	}, nil
}

func (h *HomeRecognizer) RecognizeWithSingleFile(file string) string {
	// Now let's try to classify some not yet known image.
	image := filepath.Join(dataDir, "images", file)
	fFaces, err := h.rec.RecognizeFile(image)
	if err != nil {
		fmt.Printf("Can't recognize: %v", err)
		return ""
	}
	if len(fFaces) < 1 {
		fmt.Print("no face on the image.\n")
		return ""
	}
	catID := h.rec.Classify(fFaces[0].Descriptor)
	if catID < 0 {
		fmt.Printf("Can't classify.")
		return ""
	}
	// Finally print the classified label. It should be "father".
	fmt.Println(labels[catID])
	return labels[catID]
}

func (h *HomeRecognizer) RecognizeWithSingleImage(imgData []byte) string {
	// Now let's try to classify some not yet known image.
	fFaces, err := h.rec.Recognize(imgData)
	if err != nil {
		fmt.Printf("Can't recognize: %v", err)
		return ""
	}
	if len(fFaces) < 1 {
		fmt.Print("Not a single face on the image.\n")
		return ""
	}

	for _, face := range fFaces {
		catID := h.rec.Classify(face.Descriptor)
		if catID < 0 {
			fmt.Printf("Can't classify.")
			continue
		}
		fmt.Println(labels[catID])
		return labels[catID]
	}
	return ""
}

func (h *HomeRecognizer) Close() {
	// Free the resources when you're finished.
	h.rec.Close()
}
