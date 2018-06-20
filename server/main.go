package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	api "github.com/bogem/mircoservices-demo/api"
	"google.golang.org/grpc"
)

const audioFilePath = "data/output.wav"

func main() {
	http.HandleFunc("/audio", handleAudio)
	http.HandleFunc("/convertTextToSpeech", handleConvertTextToSpeech)
	http.HandleFunc("/", handleIndex)
	http.ListenAndServe(":8081", nil)
}

func handleAudio(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "audio/wav")
	if err := copyFile(w, audioFilePath); err != nil {
		log.Fatalln(err)
	}
}

func handleConvertTextToSpeech(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	log.Printf("Converting %q\n", text)

	conn, err := grpc.Dial("192.168.0.3:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to tts server: %v", err)
	}
	defer conn.Close()

	client := api.NewTextToSpeechClient(conn)

	res, err := client.Say(context.Background(), &api.Text{Text: text})
	if err != nil {
		log.Fatalf("could not say %s: %v", text, err)
	}

	if err := ioutil.WriteFile(audioFilePath, res.Audio, 0666); err != nil {
		log.Fatalf("could not write to %s: %v", audioFilePath, err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	loadIndexFile(w)
}

func loadIndexFile(w io.Writer) {
	if err := copyFile(w, "data/index.html"); err != nil {
		log.Fatalln(err)
	}
}

func copyFile(w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	return err
}
