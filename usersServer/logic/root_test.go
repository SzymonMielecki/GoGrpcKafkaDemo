package logic

import (
	"context"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/usersServer/cache"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/usersServer/persistance"
	pb "github.com/SzymonMielecki/GoGrpcKafkaDemo/usersService"
	"reflect"
	"testing"
)

func TestServer_CheckUser(t *testing.T) {
	type fields struct {
		UnimplementedUsersServiceServer pb.UnimplementedUsersServiceServer
		db                              persistance.IDB[types.User]
		c                               cache.ICache[types.User]
	}
	type args struct {
		ctx context.Context
		in  *pb.CheckUserRequest
	}
	user := types.User{
		ID:           1,
		Username:     "user1",
		Email:        "user1@example.com",
		PasswordHash: "test",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.CheckUserResponse
		wantErr bool
	}{
		{name: "CheckUserExistsTrue", fields: fields{
			UnimplementedUsersServiceServer: pb.UnimplementedUsersServiceServer{},
			db: persistance.NewMockDBFromData([]*types.User{
				&user,
			}),
			c: cache.NewMockCache(),
		}, args: args{
			ctx: context.Background(),
			in: &pb.CheckUserRequest{
				Id:           1,
				PasswordHash: "test",
			},
		}, want: &pb.CheckUserResponse{
			Success: true,
			User:    user.ToProto(),
			Message: "User found",
		}, wantErr: false},
		{name: "CheckUserExistsFalse", fields: fields{
			UnimplementedUsersServiceServer: pb.UnimplementedUsersServiceServer{},
			db: persistance.NewMockDBFromData([]*types.User{
				&user,
			}),
			c: cache.NewMockCache(),
		}, args: args{
			ctx: context.Background(),
			in: &pb.CheckUserRequest{
				Id:           2,
				PasswordHash: "test",
			},
		}, want: &pb.CheckUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "User not found",
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedUsersServiceServer: tt.fields.UnimplementedUsersServiceServer,
				db:                              tt.fields.db,
				c:                               tt.fields.c,
			}
			got, err := s.CheckUser(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_LoginUser(t *testing.T) {
	type fields struct {
		UnimplementedUsersServiceServer pb.UnimplementedUsersServiceServer
		db                              persistance.IDB[types.User]
		c                               cache.ICache[types.User]
	}
	type args struct {
		ctx context.Context
		in  *pb.LoginUserRequest
	}
	user := types.User{
		ID:           1,
		Username:     "user1",
		Email:        "user1@example.com",
		PasswordHash: "test",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.LoginUserResponse
		wantErr bool
	}{
		{name: "LoginUserWorking", fields: fields{
			UnimplementedUsersServiceServer: pb.UnimplementedUsersServiceServer{},
			db: persistance.NewMockDBFromData([]*types.User{
				&user,
			}),
			c: cache.NewMockCache(),
		}, args: args{
			ctx: context.Background(),
			in: &pb.LoginUserRequest{
				Username:     "user1",
				Email:        "user1@example.com",
				PasswordHash: "test",
			},
		}, want: &pb.LoginUserResponse{
			Success: true,
			User:    user.ToProto(),
			Message: "Logged in as user1",
		}, wantErr: false},
		{name: "LoginUserNotWorking", fields: fields{
			UnimplementedUsersServiceServer: pb.UnimplementedUsersServiceServer{},
			db: persistance.NewMockDBFromData([]*types.User{
				&user,
			}),
			c: cache.NewMockCache(),
		}, args: args{
			ctx: context.Background(),
			in: &pb.LoginUserRequest{
				Username:     "user2",
				Email:        "user2@example.com",
				PasswordHash: "test",
			},
		}, want: &pb.LoginUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "User not found",
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedUsersServiceServer: tt.fields.UnimplementedUsersServiceServer,
				db:                              tt.fields.db,
				c:                               tt.fields.c,
			}
			got, err := s.LoginUser(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoginUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_RegisterUser(t *testing.T) {
	type fields struct {
		UnimplementedUsersServiceServer pb.UnimplementedUsersServiceServer
		db                              persistance.IDB[types.User]
		c                               cache.ICache[types.User]
	}
	type args struct {
		ctx context.Context
		in  *pb.RegisterUserRequest
	}
	user := types.User{
		ID:           1,
		Username:     "user1",
		Email:        "user1@example.com",
		PasswordHash: "test",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.RegisterUserResponse
		wantErr bool
	}{
		{name: "RegisterUserWorking", fields: fields{
			UnimplementedUsersServiceServer: pb.UnimplementedUsersServiceServer{},
			db: persistance.NewMockDBFromData([]*types.User{
				&user,
			}),
			c: cache.NewMockCache(),
		}, args: args{
			ctx: context.Background(),
			in: &pb.RegisterUserRequest{
				Username:     "user2",
				Email:        "user2@example.com",
				PasswordHash: "test",
			},
		}, want: &pb.RegisterUserResponse{
			Success: true,
			User: &pb.User{
				Id:           2,
				Username:     "user2",
				Email:        "user2@example.com",
				PasswordHash: "test",
			},
			Message: "Registered",
		}, wantErr: false},
		{name: "RegisterUserNotWorkingEmailExists", fields: fields{
			UnimplementedUsersServiceServer: pb.UnimplementedUsersServiceServer{},
			db: persistance.NewMockDBFromData([]*types.User{
				&user,
			}),
			c: cache.NewMockCache(),
		}, args: args{
			ctx: context.Background(),
			in: &pb.RegisterUserRequest{
				Username:     "user2",
				Email:        "user1@example.com",
				PasswordHash: "test",
			},
		}, want: &pb.RegisterUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "Email already exists",
		}, wantErr: false},
		{name: "RegisterUserNotWorkingUsernameExists", fields: fields{
			UnimplementedUsersServiceServer: pb.UnimplementedUsersServiceServer{},
			db: persistance.NewMockDBFromData([]*types.User{
				&user,
			}),
			c: cache.NewMockCache(),
		}, args: args{
			ctx: context.Background(),
			in: &pb.RegisterUserRequest{
				Username:     "user1",
				Email:        "user2@example.com",
				PasswordHash: "test",
			},
		}, want: &pb.RegisterUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "Username already exists",
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedUsersServiceServer: tt.fields.UnimplementedUsersServiceServer,
				db:                              tt.fields.db,
				c:                               tt.fields.c,
			}
			got, err := s.RegisterUser(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetUser(t *testing.T) {
	type fields struct {
		UnimplementedUsersServiceServer pb.UnimplementedUsersServiceServer
		db                              persistance.IDB[types.User]
		c                               cache.ICache[types.User]
	}
	type args struct {
		ctx context.Context
		in  *pb.GetUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GetUserResponse
		wantErr bool
	}{
		{
			name: "GetUserWorking",
			fields: fields{
				UnimplementedUsersServiceServer: pb.UnimplementedUsersServiceServer{},
				db: persistance.NewMockDBFromData([]*types.User{
					&types.User{
						ID:           1,
						Username:     "user1",
						Email:        "user1@example.com",
						PasswordHash: "test",
					},
				}),
				c: cache.NewMockCache(),
			},
			args: args{
				ctx: context.Background(),
				in: &pb.GetUserRequest{
					Id: 1,
				},
			},
			want: &pb.GetUserResponse{
				Success: true,
				User: &pb.User{
					Id:           1,
					Username:     "user1",
					Email:        "user1@example.com",
					PasswordHash: "test",
				},
				Message: "User found",
			},
			wantErr: false,
		},
		{
			name: "GetUserNotWorking",
			fields: fields{
				UnimplementedUsersServiceServer: pb.UnimplementedUsersServiceServer{},
				db: persistance.NewMockDBFromData([]*types.User{
					&types.User{
						ID:           1,
						Username:     "user1",
						Email:        "user1@example.com",
						PasswordHash: "test",
					},
				}),
				c: cache.NewMockCache(),
			},
			args: args{
				ctx: context.Background(),
				in: &pb.GetUserRequest{
					Id: 2,
				},
			},
			want: &pb.GetUserResponse{
				Success: false,
				User:    &pb.User{},
				Message: "User not found",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedUsersServiceServer: tt.fields.UnimplementedUsersServiceServer,
				db:                              tt.fields.db,
				c:                               tt.fields.c,
			}
			got, err := s.GetUser(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
