package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	_ "order/order-service/basic/initializer"
	"order/order-service/handler/payment"
	__ "order/proto/payment"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go payment.StartStockConsumer()
	s := grpc.NewServer()
	//pb.RegisterGreeterServer(s, &server{})
	__.RegisterPaymentServiceServer(s, &payment.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
