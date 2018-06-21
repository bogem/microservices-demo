package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"

	api "github.com/bogem/mircoservices-demo/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("p", 8080, "port to listen to")
	flag.Parse()

	log.Printf("Start auf Port %d", *port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("could not listen to port %d: %v", *port, err)
	}

	s := grpc.NewServer()
	api.RegisterTextToSpeechServer(s, server{})
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("could not serve: %v", err)
	}
}

type server struct{}

func (server) Say(_ context.Context, text *api.Text) (*api.Speech, error) {
	log.Printf("Empfang (gRPC): Abfrage an Audio mit Text: %q\n", text.Text)
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("could not create tmp file: %v", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("could not close %s: %v", f.Name(), err)
	}

	cmd := exec.Command("flite", "-t", text.Text, "-o", f.Name())
	if data, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("flite failed: %s", data)
	}

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("could not read tmp file: %v", err)
	}
	log.Println("Absendung (gRPC): Audio")
	return &api.Speech{Audio: data}, nil
}
