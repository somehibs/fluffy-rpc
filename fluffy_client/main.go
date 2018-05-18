package main

import (
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/somehibs/fluffy-rpc/fluffy"
)

const (
	address     = "localhost:7893"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewServiceControlClient(conn)

	// Contact the server and print out its response.
	name := "sphinxsearch.service"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.StatusService(ctx, &pb.ServiceRequest{Name: []string{name}})
	if err != nil {
		log.Fatalf("could not contact service: %v", err)
	}
	log.Printf("service: %s state: %s", name, r.States)
	s, err := c.StopService(ctx, &pb.ServiceRequest{Name: []string{name}})
	if err != nil {
		log.Fatalf("could not contact service: %v", err)
	}
	log.Printf("stop service: %s result: %s", name, s.Result)
	ss, err := c.StartService(ctx, &pb.ServiceRequest{Name: []string{name}})
	if err != nil {
		log.Fatalf("could not contact service: %v", err)
	}
	log.Printf("start service: %s result: %s", name, ss.Result)
}
