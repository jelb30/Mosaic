package main

import (
	"context"
	"log"
	"time"

	pb "github.com/jelb30/Mosaic/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, _ := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	c := pb.NewKVStoreClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// c.Put(ctx, &pb.PutRequest{Key: "name", Value: "Jill"})
	resp, _ := c.Get(ctx, &pb.GetRequest{Key: "name"})
	log.Printf("value=%q found=%v", resp.Value, resp.Found)
}
