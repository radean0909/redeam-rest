package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	"github.com/radean0909/redeam-rest/pkg/api/v1"
)

const (
	// sanity check
	apiVersion = "v1"
)

func main() {
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewBookServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3000) //  3 second timeout should be good
	defer cancel()

	t := time.Now().In(time.UTC)
	publishDate, _ := ptypes.TimestampProto(t)
	uVal := t.Format(time.RFC3339Nano) //  Get a unique value for identification of titles

	// Create
	req1 := v1.CreateRequest{
		Api: apiVersion,
		Book: &v1.Book{
			Title:       "title (" + uVal + ")",
			Author:      "author (" + uVal + ")",
			Publisher:   "publisher (" + uVal + ")",
			PublishDate: publishDate,
			Rating:      2.0,
			Status:      1,
		},
	}
	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res1)

	id := res1.Id

	// Read
	req2 := v1.ReadRequest{
		Api: apiVersion,
		Id:  id,
	}
	res2, err := c.Read(ctx, &req2)
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}
	log.Printf("Read result: <%+v>\n\n", res2)

	// Update
	req3 := v1.UpdateRequest{
		Api: apiVersion,
		Book: &v1.Book{
			Id:          res2.Book.Id,
			Title:       res2.Book.Title,
			Author:      res2.Book.Author + " + updated",
			Publisher:   res2.Book.Publisher,
			PublishDate: res2.Book.PublishDate,
			Rating:      res2.Book.Rating,
			Status:      res2.Book.Status,
		},
	}
	res3, err := c.Update(ctx, &req3)
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	log.Printf("Update result: <%+v>\n\n", res3)

	// ReadAll
	req4 := v1.ReadAllRequest{
		Api: apiVersion,
	}
	res4, err := c.ReadAll(ctx, &req4)
	if err != nil {
		log.Fatalf("ReadAll failed: %v", err)
	}
	log.Printf("ReadAll result: <%+v>\n\n", res4)

	// Delete
	req5 := v1.DeleteRequest{
		Api: apiVersion,
		Id:  id,
	}
	res5, err := c.Delete(ctx, &req5)
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	log.Printf("Delete result: <%+v>\n\n", res5)
}
