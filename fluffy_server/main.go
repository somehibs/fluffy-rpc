package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/somehibs/fluffy-rpc/fluffy"
	//"google.golang.org/grpc/status"
	//"google.golang.org/grpc/codes"
	"github.com/coreos/go-systemd/dbus"
	//"google.golang.org/grpc/reflection"
)

const (
	port = ":7893"
)

type FluffyService struct{}

func (s *FluffyService) StartService(ctx context.Context, in *pb.ServiceRequest) (*pb.ServiceReply, error) {
	// StartUnit(name, "fail" (or replace if you force))
	systemd, err := dbus.New()
	if err != nil {
			return nil, err
	}
	defer systemd.Close()
	reply := make(chan string)
	replies := make([]string, len(in.Name))
	for i, name := range in.Name {
		_, err := systemd.StartUnit(name, "fail", reply)
		if err != nil {
			return nil, err
		}
		replies[i] = <-reply
	}
	return &pb.ServiceReply{Result: replies}, nil
}

func (s *FluffyService) StopService(ctx context.Context, in *pb.ServiceRequest) (*pb.ServiceReply, error) {
	// StopUnit should always force "replace
	systemd, err := dbus.New()
	if err != nil {
			return nil, err
	}
	defer systemd.Close()
	reply := make(chan string)
	replies := make([]string, len(in.Name))
	for i, name := range in.Name {
		_, err := systemd.StopUnit(name, "replace", reply)
		if err != nil {
			return nil, err
		}
		replies[i] = <-reply
	}
	return &pb.ServiceReply{Result: replies}, nil
	//return &pb.ServiceReply{}, status.Error(codes.Unimplemented,"")
}

func (s *FluffyService) StatusService(ctx context.Context, in *pb.ServiceRequest) (*pb.ServiceStatusReply, error) {
	// quickly get state
	// ListUnitsByNames
	systemd, err := dbus.New()
	if err != nil {
		return nil, err
	}
	defer systemd.Close()
	status, err := systemd.ListUnitsByNames(in.Name)
	if err != nil {
		return nil, err
	}
	states := make(map[string]string, len(status))
	for _, state := range status {
		states[state.Name] = state.ActiveState
	}
	return &pb.ServiceStatusReply{States: states}, nil
	//return &pb.ServiceStatusReply{State: pb.ServiceStatusReply_UNKNOWN, UnixLastStarted: 0}, status.Error(codes.Unimplemented, "")
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServiceControlServer(s, &FluffyService{})
	// Register debug reflection service on gRPC FluffyService.
	//reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
