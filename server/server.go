package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	proto "github.com/EmanuelFeij/MinderaPractice/protos/company"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

/*
Host	localhost
Port	5432
User	your system user name
Database	same as user
Password	none
Connection URL	postgresql://localhost
*/

const (
	host   = "localhost"
	port   = 5432
	user   = "emanuelfeijo"
	dbname = "emanuelfeijo"

	myServer = "localhost:8080"
)

func createDbConnection() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatalf("cannot connect to db :%v", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS COMPANY3(
		NAME           TEXT    NOT NULL,
		ID 				INT      NOT NULL,
		PROFESSION     TEXT,
		AGE            INT     NOT NULL
	 );`)
	if err != nil {
		log.Fatalf("error creating table :%v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Not connected: %v", err)
	} else {
		log.Println("connected")
	}
	return db
}

type MyServer struct {
	db *sql.DB
	proto.UnimplementedCompanyServer
}

func NewMyServer() *MyServer {
	return &MyServer{db: createDbConnection()}
}

func (s *MyServer) CloseDbConection() {
	s.db.Close()
}

func NewUser(username string, id int32, profession string, age int32) *proto.User {
	return &proto.User{
		Username:   &proto.UserName{Name: username},
		Id:         &proto.UserID{Id: id},
		Profession: profession,
		Age:        age,
	}
}

func (s *MyServer) GetAllUsers(em *proto.EmptyMessage, stream proto.Company_GetAllUsersServer) error {
	log.Println("here")
	rows, err := s.db.Query(`SELECT * FROM "company3" `)
	if err != nil {
		log.Fatalf("error getting Users %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var id int32
		var profession string
		var age int32
		rows.Scan(&name, &id, &profession, &age)
		stream.Send(NewUser(name, id, profession, age))
	}
	return nil
}
func (s *MyServer) GetUserByName(*proto.UserName, proto.Company_GetUserByNameServer) error {
	return nil
}
func (s *MyServer) GetUserByID(*proto.UserID, proto.Company_GetUserByIDServer) error {
	return nil
}
func (s *MyServer) AddUser(ctx context.Context, user *proto.User) (*proto.Error, error) {
	insertDynStmt := `insert into "company3"("name", "id", "profession", "age") values($1, $2, $3, $4)`
	_, err := s.db.Exec(insertDynStmt, user.Username.Name, user.Id.Id, user.Profession, user.Age)
	if err != nil {
		log.Fatalf("erro adicionar user: %v", err)
	}
	return &proto.Error{Yes: false, No: false}, nil
}
func (s *MyServer) AddUserSeveralUsers(proto.Company_AddUserSeveralUsersServer) error {
	return nil
}
func (s *MyServer) DeleteUser(context.Context, *proto.User) (*proto.Error, error) {
	return &proto.Error{Yes: false, No: false}, nil
}

func main() {
	s := NewMyServer()
	defer s.CloseDbConection()
	cancel := make(chan error)

	go func(s *MyServer) {
		lis, err := net.Listen("tcp", myServer)
		if err != nil {
			log.Fatalf("Cannot listen at %v with the error %v", myServer, err)
		}
		grpcServer := grpc.NewServer()

		proto.RegisterCompanyServer(grpcServer, s)
		cancel <- grpcServer.Serve(lis)

	}(s)
	log.Println(<-cancel)
}
