package fluffy

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

type Client struct {
	conn *grpc.ClientConn
	sc pb.ServiceControlClient
}

var ctxTimeout = 10*time.Second

func New(addr string) (client *Client, err error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %+v", err)
		return nil, err
	}
	c := pb.NewServiceControlClient(conn)
	return &Client{conn: conn, sc: c}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) StartService(service string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	r, err := c.sc.StartService(ctx, &pb.ServiceRequest{Name: []string{service+".service"}})
	if err != nil {
		log.Fatalf("could not contact service: %v", err)
	}
	return r.Result[0], nil
}

func (c *Client) StopService(service string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	r, err := c.sc.StopService(ctx, &pb.ServiceRequest{Name: []string{service+".service"}})
	if err != nil {
		log.Fatalf("could not contact service: %v", err)
	}
	return r.Result[0], nil
}

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
