package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/nfnt/resize"
	"gocv.io/x/gocv"
)

var (
	timeout = 3 * time.Second
	port    = flag.String("p", "8080", "server port")
	host    = flag.String("h", "", "server ip")
	picdir  = flag.String("d", "/tmp/", "images temp dir")
)

func main() {
	flag.Parse()
	http.HandleFunc("/thumb", thumb)
	log.Println("Starting server ...")
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", *host, *port), nil); err != nil {
		log.Fatalln("Start server error:", err)
	}
}

func thumb(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("url")
	width, err := strconv.ParseUint(r.FormValue("w"), 10, 64)
	if err != nil {
		width = 0
	}
	height, err := strconv.ParseUint(r.FormValue("h"), 10, 64)
	if err != nil {
		height = 0
	}
	t, err := strconv.ParseFloat(r.FormValue("t"), 64)
	if err != nil {
		t = 1000
	}

	outname := path.Join(*picdir, time.Now().Format("2006/01/02"), fmt.Sprintf("%s_%d_%d_%0.0f.jpg", stringMd5(filename), width, height, t))

	log.Println(fmt.Sprintf("Get the params:[file=%s,width=%d,height=%d,time=%0.3f],outname:%s", filename, width, height, t, outname))

	if !exist(outname) {
		createDir(path.Dir(outname))
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		img, err := createImageByOpenCV(filename, outname, uint(width), uint(height), t)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}

		w.Header().Set("Content-Type", "image/png")
		http.ServeContent(w, r, outname, time.Time{}, bytes.NewReader(buf.Bytes()))

		return
	}

	log.Println("Get the cache file:", outname)
	http.ServeFile(w, r, outname)
}

func createImageByOpenCV(filename, outname string, width, height uint, time float64) (i image.Image, err error) {
	cv, err := gocv.VideoCaptureFile(filename)
	if err != nil {
		return i, err
	}

	frames := cv.Get(gocv.VideoCaptureFrameCount)
	fps := cv.Get(gocv.VideoCaptureFPS)
	videoWidth := uint(cv.Get(gocv.VideoCaptureFrameWidth))
	if width < 10 || width > videoWidth {
		width = videoWidth
	}
	videoHeight := uint(cv.Get(gocv.VideoCaptureFrameHeight))
	if height < 10 || height > videoHeight {
		height = videoHeight
	}

	duration := frames * 1000 / fps

	frames = (time / duration) * frames

	cv.Set(gocv.VideoCapturePosFrames, frames)

	img := gocv.NewMat()
	defer img.Close()

	cv.Read(&img)

	if img.Empty() {
		return i, fmt.Errorf("this video duration is %0.3f,but you need %0.3f", duration, time)
	}

	imageObject, err := img.ToImage()
	if err != nil {
		return i, err
	}

	imageObject = resize.Resize(width, height, imageObject, resize.Lanczos3)
	if img, err := gocv.ImageToMatRGB(imageObject); err == nil {
		gocv.IMWrite(outname, img)
	}

	return imageObject, nil
}

func createDir(path string) error {
	if !exist(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func stringMd5(text string) string {
	md5hash := md5.New()
	md5hash.Write([]byte(text))
	return hex.EncodeToString(md5hash.Sum(nil))
}
