package v1

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/radean0909/redeam-rest/pkg/api/v1"
)

func Test_bookServiceServer_Create(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error when opening stub database connection: '%s'", err)
	}
	defer db.Close()
	s := NewBookServiceServer(db)
	now := time.Now().In(time.UTC)
	publishDate, _ := ptypes.TimestampProto(now)

	type args struct {
		ctx context.Context
		req *v1.CreateRequest
	}

	rows := sqlmock.NewRows([]string{"api", "id"}).AddRow(1, 1)
	tests := []struct {
		name    string
		s       v1.BookServiceServer
		args    args
		mock    func()
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
			mock: func() {
				mock.ExpectQuery("INSERT INTO Book").WithArgs("title", "author", "publisher", now, 2.0, 1).
					WillReturnRows(rows)
			},
			want: &v1.CreateResponse{
				Api: "v1",
				Id:  1,
			},
		},
		{
			name: "Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1000",
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
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Invalid publishDate format",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1",
					Book: &v1.Book{
						Title:     "title",
						Author:    "author",
						Publisher: "publisher",
						PublishDate: &timestamp.Timestamp{
							Seconds: 1,
							Nanos:   -1,
						},
						Rating: 2.0,
						Status: 1,
					},
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "INSERT failed",
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
			mock: func() {
				mock.ExpectQuery("INSERT INTO Book").WithArgs("title", "author", "publisher", now, 2.0, 1).
					WillReturnError(errors.New("INSERT failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bookServiceServer_Read(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error when opening stub database connection: '%s'", err)
	}
	defer db.Close()
	s := NewBookServiceServer(db)
	now := time.Now().In(time.UTC)
	publishDate, _ := ptypes.TimestampProto(now)

	type args struct {
		ctx context.Context
		req *v1.ReadRequest
	}
	tests := []struct {
		name    string
		s       v1.BookServiceServer
		args    args
		mock    func()
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
			mock: func() {
				rows := sqlmock.NewRows([]string{"ID", "Title", "Author", "Publisher", "PublishDate", "Rating", "Status"}).
					AddRow(1, "title", "author", "publisher", now, 2.0, 1)
				mock.ExpectQuery("SELECT (.+) FROM Book").WithArgs(1).WillReturnRows(rows)
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
		},
		{
			name: "Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v1",
					Id:  1,
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "SELECT failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v1",
					Id:  1,
				},
			},
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM Book").WithArgs(1).
					WillReturnError(errors.New("SELECT failed"))
			},
			wantErr: true,
		},
		{
			name: "Not found",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v1",
					Id:  1,
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"ID", "Title", "Author", "Publisher", "PublishDate", "Rating", "Status"})
				mock.ExpectQuery("SELECT (.+) FROM Book").WithArgs(1).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Read(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bookServiceServer_Update(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error when opening stub database connection: '%s'", err)
	}
	defer db.Close()
	s := NewBookServiceServer(db)
	now := time.Now().In(time.UTC)
	publishDate, _ := ptypes.TimestampProto(now)

	type args struct {
		ctx context.Context
		req *v1.UpdateRequest
	}
	tests := []struct {
		name    string
		s       v1.BookServiceServer
		args    args
		mock    func()
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
						Title:       "new title",
						Author:      "new author",
						Publisher:   "new publisher",
						PublishDate: publishDate,
						Rating:      2.0,
						Status:      1,
					},
				},
			},
			mock: func() {
				mock.ExpectExec("UPDATE Book").WithArgs("new title", "new author", "new publisher", now, 2.0, 1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: &v1.UpdateResponse{
				Api:     "v1",
				Updated: 1,
			},
		},
		{
			name: "Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					Api: "v1",
					Book: &v1.Book{
						Id:          1,
						Title:       "new title",
						Author:      "new author",
						Publisher:   "new publisher",
						PublishDate: publishDate,
						Rating:      2.0,
						Status:      1,
					},
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Invalid publishDate field format",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					Api: "v1",
					Book: &v1.Book{
						Title:     "new title",
						Author:    "new author",
						Publisher: "new publisher",
						PublishDate: &timestamp.Timestamp{
							Seconds: 1,
							Nanos:   -1,
						},
						Rating: 2.0,
						Status: 1,
					},
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "UPDATE failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					Api: "v1",
					Book: &v1.Book{
						Id:          1,
						Title:       "new title",
						Author:      "new author",
						Publisher:   "new publisher",
						PublishDate: publishDate,
						Rating:      2.0,
						Status:      1,
					},
				},
			},
			mock: func() {
				mock.ExpectExec("UPDATE Book").WithArgs("new title", "new author", "new publisher", now, 2.0, 1, 1).
					WillReturnError(errors.New("UPDATE failed"))
			},
			wantErr: true,
		},
		{
			name: "RowsAffected failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					Api: "v1",
					Book: &v1.Book{
						Id:          1,
						Title:       "new title",
						Author:      "new author",
						Publisher:   "new publisher",
						PublishDate: publishDate,
						Rating:      2.0,
						Status:      1,
					},
				},
			},
			mock: func() {
				mock.ExpectExec("UPDATE Book").WithArgs("new title", "new author", "new publisher", now, 2.0, 1, 1).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("RowsAffected failed")))
			},
			wantErr: true,
		},
		{
			name: "Not Found",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.UpdateRequest{
					Api: "v1",
					Book: &v1.Book{
						Id:          1,
						Title:       "new title",
						Author:      "new author",
						Publisher:   "new publisher",
						PublishDate: publishDate,
						Rating:      2.0,
						Status:      1,
					},
				},
			},
			mock: func() {
				mock.ExpectExec("UPDATE Book").WithArgs("new title", "new author", "new publisher", now, 2.0, 1, 1).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Update(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bookServiceServer_Delete(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error when opening stub database connection: '%s'", err)
	}
	defer db.Close()
	s := NewBookServiceServer(db)

	type args struct {
		ctx context.Context
		req *v1.DeleteRequest
	}
	tests := []struct {
		name    string
		s       v1.BookServiceServer
		args    args
		mock    func()
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
			mock: func() {
				mock.ExpectExec("DELETE FROM Book").WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: &v1.DeleteResponse{
				Api:     "v1",
				Deleted: 1,
			},
		},
		{
			name: "Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.DeleteRequest{
					Api: "v1",
					Id:  1,
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "DELETE failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.DeleteRequest{
					Api: "v1",
					Id:  1,
				},
			},
			mock: func() {
				mock.ExpectExec("DELETE FROM Book").WithArgs(1).
					WillReturnError(errors.New("DELETE failed"))
			},
			wantErr: true,
		},
		{
			name: "RowsAffected failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.DeleteRequest{
					Api: "v1",
					Id:  1,
				},
			},
			mock: func() {
				mock.ExpectExec("DELETE FROM Book").WithArgs(1).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("RowsAffected failed")))
			},
			wantErr: true,
		},
		{
			name: "Not Found",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.DeleteRequest{
					Api: "v1",
					Id:  1,
				},
			},
			mock: func() {
				mock.ExpectExec("DELETE FROM Book").WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Delete(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bookServiceServer_ReadAll(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error when opening stub database connection: '%s'", err)
	}
	defer db.Close()
	s := NewBookServiceServer(db)
	now1 := time.Now().In(time.UTC)
	publishDate1, _ := ptypes.TimestampProto(now1)
	now2 := time.Now().In(time.UTC)
	publishDate2, _ := ptypes.TimestampProto(now2)

	type args struct {
		ctx context.Context
		req *v1.ReadAllRequest
	}
	tests := []struct {
		name    string
		s       v1.BookServiceServer
		args    args
		mock    func()
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
			mock: func() {
				rows := sqlmock.NewRows([]string{"ID", "Title", "Author", "Publisher", "PublishDate", "Rating", "Status"}).
					AddRow(1, "title1", "author1", "publisher1", now1, 2.0, 1).
					AddRow(2, "title2", "author2", "publisher2", now2, 3.0, 2)
				mock.ExpectQuery("SELECT (.+) FROM Book").WillReturnRows(rows)
			},
			want: &v1.ReadAllResponse{
				Api: "v1",
				Books: []*v1.Book{
					{
						Id:          1,
						Title:       "title1",
						Author:      "author1",
						Publisher:   "publisher1",
						PublishDate: publishDate1,
						Rating:      2.0,
						Status:      1,
					},
					{
						Id:          2,
						Title:       "title2",
						Author:      "author2",
						Publisher:   "publisher2",
						PublishDate: publishDate2,
						Rating:      3.0,
						Status:      2,
					},
				},
			},
		},
		{
			name: "Empty",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadAllRequest{
					Api: "v1",
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"ID", "Title", "Author", "Publisher", "PublishDate", "Rating", "Status"})
				mock.ExpectQuery("SELECT (.+) FROM Book").WillReturnRows(rows)
			},
			want: &v1.ReadAllResponse{
				Api:   "v1",
				Books: []*v1.Book{},
			},
		},
		{
			name: "Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadAllRequest{
					Api: "v1",
				},
			},
			mock:    func() {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.ReadAll(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.ReadAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
