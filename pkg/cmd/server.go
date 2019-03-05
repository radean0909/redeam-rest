package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/radean0909/redeam-rest/pkg/protocol/grpc"
	"github.com/radean0909/redeam-rest/pkg/service/v1"
)

type Config struct {
	GRPCPort string

	PgDBHost     string
	PgDBUser     string
	PgDBPassword string
	PgDBSchema   string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.PgDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.PgDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.PgDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.PgDBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.PgDBUser,
		cfg.PgDBPassword,
		cfg.PgDBHost,
		cfg.PgDBSchema,
		param)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := v1.NewBookServiceServer(db)

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
