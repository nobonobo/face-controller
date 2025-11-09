package main

import (
	"flag"
	"image"
	"image/color"
	"log"
	"time"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

const (
	haarCascadeFile = "haarcascade_frontalface_default.xml"
	//haarCascadeFile = "haarcascade_eye_tree_eyeglasses.xml"
	//haarCascadeFile = "haarcascade_mcs_nose.xml"
	//haarCascadeFile = "haarcascade_upperbody.xml"
	//haarCascadeFile = "haarcascade_mcs_eyepair_big.xml"
)

func main() {
	capture := 1
	port := ""
	view := false
	flag.BoolVar(&view, "view", view, "show window")
	flag.IntVar(&capture, "capture", capture, "capture device index")
	flag.StringVar(&port, "port", port, "serial port name")
	flag.Parse()
	webcam, err := gocv.OpenVideoCapture(capture)
	if err != nil {
		log.Fatalf("Error opening video capture device: %v\n", err)
	}
	defer webcam.Close()
	classifier := gocv.NewCascadeClassifier()
	if !classifier.Load(haarCascadeFile) {
		log.Fatalf("Error reading cascade file: %v\n", haarCascadeFile)
	}
	defer classifier.Close()
	tracker := contrib.NewTrackerCSRT()
	defer tracker.Close()
	tracking := false // トラッキング状態のフラグ
	var trackRect image.Rectangle
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()
	dst := gocv.NewMat()
	service, err := NewJoyStickService(port)
	if err != nil {
		log.Fatalf("Error opening serial port: %v\n", err)
	}
	defer service.Close()
	toggle := false
	ticker := time.NewTicker(time.Second / 30)
	for range ticker.C {
		v := window.WaitKey(1)
		if v > 0 {
			log.Println("key:", v)
		}
		switch v {
		default:
		case 0x20:
			tracking = false
		case 27, 113:
			return
		case 97:
			toggle = !toggle
			service.SetButton(0, toggle)
		}
		if ok := webcam.Read(&img); !ok || img.Empty() {
			continue
		}
		if view {
			dst = img.Clone()
		} else {
			if dst.Empty() {
				img.CopyTo(&dst)
			}
			gocv.Rectangle(&dst, image.Rect(0, 0, dst.Cols(), dst.Rows()), color.RGBA{G: 255, A: 255}, -1)
		}
		if !tracking {
			rects := classifier.DetectMultiScale(img)
			if len(rects) > 0 {
				// 1番大きい顔を追跡対象に選ぶ例
				maxIdx := 0
				maxArea := 0
				for i, r := range rects {
					area := r.Dx() * r.Dy()
					if area > maxArea {
						maxArea = area
						maxIdx = i
					}
				}
				trackRect = rects[maxIdx]

				// トラッカー初期化
				tracker.Init(img, trackRect)
				tracking = true
			}
		} else {
			// トラッキング更新
			newRect, ok := tracker.Update(img)
			if ok {
				trackRect = newRect
				gocv.Rectangle(&dst, trackRect, (color.RGBA{0, 0, 255, 0}), 3)
			} else {
				// トラッキング失敗時 リセット
				tracking = false
			}
		}
		const K = 64.0
		window.IMShow(dst)
		dx := K * (float64(trackRect.Max.X+trackRect.Min.X)/2 - float64(img.Size()[1]/2))
		dy := K * (float64(trackRect.Max.Y+trackRect.Min.Y)/2 - float64(img.Size()[0]/2))
		//log.Println(dx, dy)
		if err := service.SetAxis(2, int(dx)); err != nil {
			log.Println(err)
		}
		if err := service.SetAxis(3, int(dy)); err != nil {
			log.Println(err)
		}
		if err := service.SendState(); err != nil {
			log.Println(err)
		}
	}
}
