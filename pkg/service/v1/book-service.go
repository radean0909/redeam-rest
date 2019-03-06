package v1

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/radean0909/redeam-rest/pkg/api/v1"
)

const (
	apiVersion = "v1" // sanity check
)

type bookServiceServer struct {
	db *sql.DB
}

func NewBookServiceServer(db *sql.DB) v1.BookServiceServer {
	return &bookServiceServer{db: db}
}

// version sanity check
func (s *bookServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: implemented API version '%s'\trequested version '%s'", apiVersion, api)
		}
	}
	return nil
}

// connect to the next DB in the pool, this allows for horizontal scaling
func (s *bookServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database: "+err.Error())
	}
	return c, nil
}

// Create request/response from proto definition
func (s *bookServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// Check for proper timestamp formatting
	fmt.Println(req.Book)
	publish_date, err := ptypes.Timestamp(req.Book.PublishDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "publish_date field has invalid format: "+err.Error())
	}

	res, err := c.ExecContext(ctx, "INSERT INTO Book(Title, Author, Publisher, PublishDate, Rating, Status) VALUES(?, ?, ?, ?, ?, ?)",
		req.Book.Title, req.Book.Author, req.Book.Publisher, publish_date, req.Book.Rating, req.Book.Status)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert: "+err.Error())
	}

	// get ID
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id: "+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

// Read request/response from proto definition
func (s *bookServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT ID, Title, Publisher, PublishDate, Rating, Status FROM Book WHERE ID=?",
		req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "couldn't select: "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data: "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("cannot find ID='%d'",
			req.Id))
	}

	var row v1.Book
	var publish_date time.Time
	if err := rows.Scan(&row.Id, &row.Title, &row.Author, &row.Publisher, &publish_date, &row.Rating, &row.Status); err != nil {
		return nil, status.Error(codes.Unknown, "couldn't retrieve field values: "+err.Error())
	}
	row.PublishDate, err = ptypes.TimestampProto(publish_date)
	if err != nil {
		return nil, status.Error(codes.Unknown, "publish_Date field has invalid format: "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("multiple rows with ID='%d'",
			req.Id))
	}

	return &v1.ReadResponse{
		Api:  apiVersion,
		Book: &row,
	}, nil

}

// Update request/response from proto definition
func (s *bookServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	publish_date, err := ptypes.Timestamp(req.Book.PublishDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "publish_date field has invalid format: "+err.Error())
	}

	res, err := c.ExecContext(ctx, "UPDATE Book SET Title=?, Author=?, Publisher=?, PublishDate=?, Rating=?, Status=? WHERE ID=?",
		req.Book.Title, req.Book.Author, req.Book.Publisher, publish_date, req.Book.Rating, req.Book.Status, req.Book.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update: "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "retrieve rows affected value error: "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ID='%d' not found",
			req.Book.Id))
	}

	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: rows,
	}, nil
}

// Delete request/response from proto definition
func (s *bookServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "DELETE FROM Book WHERE ID=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete: "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value: "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ID='%d' is not found",
			req.Id))
	}

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: rows,
	}, nil
}

// Read all request/response from proto definition
func (s *bookServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT ID, Title, Publisher, PublishDate, Rating, Status FROM Book")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to SELECT: "+err.Error())
	}
	defer rows.Close()

	var publish_date time.Time
	list := []*v1.Book{}
	for rows.Next() {
		row := new(v1.Book)
		if err := rows.Scan(&row.Id, &row.Title, &row.Author, &row.Publisher, &publish_date, &row.Rating, &row.Status); err != nil {
			return nil, status.Error(codes.Unknown, "couldn't retrieve field values: "+err.Error())
		}
		row.PublishDate, err = ptypes.TimestampProto(publish_date)
		if err != nil {
			return nil, status.Error(codes.Unknown, "publish_Date field has invalid format: "+err.Error())
		}
		list = append(list, row)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "couldn't retrieve: "+err.Error())
	}

	return &v1.ReadAllResponse{
		Api:   apiVersion,
		Books: list,
	}, nil
}
