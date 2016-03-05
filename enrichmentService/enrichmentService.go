package main

import (
	"log"
	"net"

	pb "github.com/garrywright/tibial/messaging"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"time"
)

const (
	port = ":50051"
)

// server is used to implement tibial.SenderServiceServer
type server struct {
	logger *log.Logger
}

// SayHello implements tibial.SenderService
func (s *server) SendMessage(ctx context.Context, in *pb.Message) (*pb.MessageReply, error) {
	t := time.Now()
	s.logger.Printf("Received on %d-%02d-%02d %02d:%02d:%02d : %s \n", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), in.Body)
	return &pb.MessageReply{Reply: "Hello " + in.Body}, nil
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterSenderServiceServer(s, &server{log.New(os.Stdout, "[enrichmentService] ", 0)})

	s.Serve(lis)
}