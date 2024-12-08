package main

import (
	user_v1 "github.com/laiker/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":50052")

	if err != nil {
		panic(err)
	}

	g := grpc.NewServer()
	reflection.Register(g)
	user_v1.RegisterNoteV1Server(g, &listener)

	log.Printf("server listening at %v", listener.Addr())

	if err = g.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func CreateUser() {

}

func GetUser() {

}

func UpdateUser() {

}

func DeleteUser() {

}
