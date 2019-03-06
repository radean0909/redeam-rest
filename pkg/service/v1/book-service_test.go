package v1

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"

	"github.com/radean0909/redeam-rest/pkg/api/v1"
)

func connectToDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"db",
		5432,
		"postgres-dev",
		"sn34kyp4ssw0rD",
		"redeam-library")
	db, err := sql.Open("postgres", dsn)
	return db, err
}

func clearTable(db *sql.DB) {
	db.Exec("DELETE FROM book")
	db.Exec("ALTER SEQUENCE book_id_seq RESTART WITH 1")
}

func addEntries(num int) {
	ctx := context.Background()
	// Get the DB
	db, err := connectToDB()

	if err != nil {
		fmt.Errorf("Couldn't connect to DB.")
		return
	}
	// Clear the table
	clearTable(db)

	// Start the server
	s := NewBookServiceServer(db)

	now, _ := time.Parse("RFC3339", "2002-10-02T10:00:00")
	publishDate, _ := ptypes.TimestampProto(now)

	for i := 1; i <= num; i++ {
		s.Create(ctx, &v1.CreateRequest{
			Api: "v1",
			Book: &v1.Book{
				Title:       "title" + strconv.Itoa(i),
				Author:      "author" + strconv.Itoa(i),
				Publisher:   "publisher",
				PublishDate: publishDate,
				Rating:      1.0,
				Status:      2,
			},
		})
	}

}
func Test_bookServiceServer_Create(t *testing.T) {
	ctx := context.Background()
	// Get the DB
	db, err := connectToDB()

	if err != nil {
		t.Errorf("Couldn't connect to DB.")
	}
	// Clear the table
	clearTable(db)

	// Start the server
	s := NewBookServiceServer(db)

	// run tests
	type args struct {
		ctx context.Context
		req *v1.CreateRequest
	}

	now, _ := time.Parse("RFC3339", "2002-10-02T10:00:00")
	publishDate, _ := ptypes.TimestampProto(now)

	tests := []struct {
		name    string
		s       v1.BookServiceServer
		args    args
		want    *v1.CreateResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1",
					Book: &v1.Book{
						Title:       "title",
						Author:      "author",
						Publisher:   "publisher",
						PublishDate: publishDate,
						Rating:      2.0,
						Status:      1,
					},
				},
			},
			want: &v1.CreateResponse{
				Api: "v1",
				Id:  1,
			},
			wantErr: false,
		},
		{
			name: "Invalid API Version",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v10000",
					Book: &v1.Book{
						Title:       "title",
						Author:      "author",
						Publisher:   "publisher",
						PublishDate: publishDate,
						Rating:      2.0,
						Status:      1,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.Create() = %v, want %v", got, tt.want)
				return
			}
		})
	}

	defer db.Close()

}

func Test_bookServiceServer_Read(t *testing.T) {
	ctx := context.Background()
	// Get the DB
	db, err := connectToDB()

	if err != nil {
		t.Errorf("Couldn't connect to DB.")
	}

	// Start the server
	s := NewBookServiceServer(db)

	// run tests
	type args struct {
		ctx context.Context
		req *v1.ReadRequest
	}
	now, _ := time.Parse("RFC3339", "2002-10-02T10:00:00")
	publishDate, _ := ptypes.TimestampProto(now)

	tests := []struct {
		name    string
		s       v1.BookServiceServer
		args    args
		want    *v1.ReadResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v1",
					Id:  1,
				},
			},
			want: &v1.ReadResponse{
				Api: "v1",
				Book: &v1.Book{
					Id:          1,
					Title:       "title",
					Author:      "author",
					Publisher:   "publisher",
					PublishDate: publishDate,
					Rating:      2.0,
					Status:      1,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid API Version",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v10000",
					Id:  1,
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Id value",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v1",
					Id:  1000,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Read(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.Read() = %v, want %v", got, tt.want)
				return
			}
		})
	}

	defer db.Close()
}

func Test_bookServiceServer_Update(t *testing.T) {
	ctx := context.Background()
	// Get the DB
	db, err := connectToDB()

	if err != nil {
		t.Errorf("Couldn't connect to DB.")
	}

	// Start the server
	s := NewBookServiceServer(db)

	// run tests
	type args struct {
		ctx context.Context
		req *v1.UpdateRequest
	}

	now, _ := time.Parse("RFC3339", "2002-10-02T10:00:00")
	publishDate, _ := ptypes.TimestampProto(now)

	tests := []struct {
		name    string
		s       v1.BookServiceServer
		args    args
		want    *v1.UpdateResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					Api: "v1",
					Book: &v1.Book{
						Id:          1,
						Title:       "title (UPDATED)",
						Author:      "author",
						Publisher:   "publisher",
						PublishDate: publishDate,
						Rating:      2.0,
						Status:      1,
					},
				},
			},
			want: &v1.UpdateResponse{
				Api:     "v1",
				Updated: 1,
			},
			wantErr: false,
		},
		{
			name: "Invalid API Version",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					Api: "v10000",
					Book: &v1.Book{
						Id:          1,
						Title:       "title (UPDATED)",
						Author:      "author",
						Publisher:   "publisher",
						PublishDate: publishDate,
						Rating:      2.0,
						Status:      1,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Id value",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					Api: "v1",
					Book: &v1.Book{
						Id:          1000,
						Title:       "title (UPDATED)",
						Author:      "author",
						Publisher:   "publisher",
						PublishDate: publishDate,
						Rating:      2.0,
						Status:      1,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Update(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.Read() = %v, want %v", got, tt.want)
				return
			}
		})
	}

	defer db.Close()
}

func Test_bookServiceServer_Delete(t *testing.T) {
	ctx := context.Background()
	// Get the DB
	db, err := connectToDB()

	if err != nil {
		t.Errorf("Couldn't connect to DB.")
	}

	// Start the server
	s := NewBookServiceServer(db)

	// run tests
	type args struct {
		ctx context.Context
		req *v1.DeleteRequest
	}

	tests := []struct {
		name    string
		s       v1.BookServiceServer
		args    args
		want    *v1.DeleteResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.DeleteRequest{
					Api: "v1",
					Id:  1,
				},
			},
			want: &v1.DeleteResponse{
				Api:     "v1",
				Deleted: 1,
			},
			wantErr: false,
		},
		{
			name: "Invalid API Version",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.DeleteRequest{
					Api: "v1000",
					Id:  1,
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Id value",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.DeleteRequest{
					Api: "v1",
					Id:  1,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Delete(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.Read() = %v, want %v", got, tt.want)
				return
			}
		})
	}

	defer db.Close()
}

func Test_bookServiceServer_ReadAll(t *testing.T) {
	ctx := context.Background()
	// Get the DB
	db, err := connectToDB()

	if err != nil {
		t.Errorf("Couldn't connect to DB.")
	}

	// Start the server
	s := NewBookServiceServer(db)

	// Add some entries
	addEntries(4)

	// run tests
	type args struct {
		ctx context.Context
		req *v1.ReadAllRequest
	}

	now, _ := time.Parse("RFC3339", "2002-10-02T10:00:00")
	publishDate, _ := ptypes.TimestampProto(now)

	tests := []struct {
		name    string
		s       v1.BookServiceServer
		args    args
		want    *v1.ReadAllResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadAllRequest{
					Api: "v1",
				},
			},
			want: &v1.ReadAllResponse{
				Api: "v1",
				Books: []*v1.Book{
					{
						Id:          1,
						Title:       "title1",
						Author:      "author1",
						Publisher:   "publisher",
						PublishDate: publishDate,
						Rating:      1.0,
						Status:      2,
					},
					{
						Id:          2,
						Title:       "title2",
						Author:      "author2",
						Publisher:   "publisher",
						PublishDate: publishDate,
						Rating:      1.0,
						Status:      2,
					},
					{
						Id:          3,
						Title:       "title3",
						Author:      "author3",
						Publisher:   "publisher",
						PublishDate: publishDate,
						Rating:      1.0,
						Status:      2,
					},
					{
						Id:          4,
						Title:       "title4",
						Author:      "author4",
						Publisher:   "publisher",
						PublishDate: publishDate,
						Rating:      1.0,
						Status:      2,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid API Version",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadAllRequest{
					Api: "v1000",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ReadAll(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.Read() = %v, want %v", got, tt.want)
				return
			}
		})
	}

	defer db.Close()
}
