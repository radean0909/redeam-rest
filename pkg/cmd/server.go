package cmd

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/radean0909/redeam-rest/pkg/protocol/grpc"
	"github.com/radean0909/redeam-rest/pkg/protocol/rest"
	"github.com/radean0909/redeam-rest/pkg/service/v1"
)

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"127.0.0.1",
		5432,
		"postgres-dev",
		"sn34kyp4$$w0rD",
		"redeam-library")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := v1.NewBookServiceServer(db)

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, "9090", "8080")
	}()

	return grpc.RunServer(ctx, v1API, "9090")
}
