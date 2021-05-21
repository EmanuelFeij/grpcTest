package main

import (
	"context"
	"io"
	"log"

	proto "github.com/EmanuelFeij/MinderaPractice/protos/company"
	"google.golang.org/grpc"
)

const (
	myServer = "localhost:8080"
)

func NewUser(username string, id int32, profession string, age int32) *proto.User {
	return &proto.User{
		Username:   &proto.UserName{Name: username},
		Id:         &proto.UserID{Id: id},
		Profession: profession,
		Age:        age,
	}
}

func AddUserClient(c proto.CompanyClient, user *proto.User) error {
	ctx := context.Background()
	_, err := c.AddUser(ctx, user)
	return err
}

func GetAllUsersClient(c proto.CompanyClient) []*proto.User {
	users := []*proto.User{}
	ctx := context.Background()
	stream, err := c.GetAllUsers(ctx, &proto.EmptyMessage{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		user, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users
}

func main() {
	lis, err := grpc.Dial(myServer, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	c := proto.NewCompanyClient(lis)

	emanuel := NewUser("Tone", 12, "maquina", 28)
	lola := NewUser("Lola", 2, "menosmaquina", 28)
	_ = AddUserClient(c, emanuel)
	_ = AddUserClient(c, lola)
	users := GetAllUsersClient(c)
	log.Println(users)
}
