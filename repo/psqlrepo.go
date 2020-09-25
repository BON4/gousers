package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	pb "user-grpc/pkg/api"
	utils "user-grpc/utils"

	"github.com/golang/protobuf/ptypes/empty"
)

//PsqlRepository - repo for psql database
type PsqlRepository struct {
	DB *sql.DB
}

//NewPsqlREpository - create new Psql repo
func NewPsqlREpository(database *sql.DB) Repository {
	//TODO setup connection

	//TODO move database init someware else
	err := initTestDB(database)
	if err != nil {
		log.Fatal(err)
	}

	return &PsqlRepository{DB: database}
}

func initTestDB(db *sql.DB) error {
	_, err := db.Exec("create table gousers(id serial not null, email varchar unique, password bytea)")
	if err != nil {
		return err
	}
	return nil
}

//ListUsers - get list of users from psql gousers table with certain boarders. PageToken >= list >= PageSize
func (r *PsqlRepository) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	var sqlStr string = "SELECT id, email FROM (select *, rank() over(order by ID) as my_rank from gousers) as x where my_rank>=$1 and my_rank <= $2"
	selectedPage, err := strconv.Atoi(req.GetPageToken())

	listUsers := pb.ListUsersResponse{}

	if err != nil {
		return nil, errors.New("Invalid page token")
	}

	if selectedPage <= 0 {
		return nil, errors.New("Invalid page token")
	}

	//Check selected page if it is in bound
	var maxPageToken int32 = 0
	err = r.DB.QueryRowContext(ctx, "select count(*) from gousers").Scan(&maxPageToken)
	if err != nil {
		return nil, err
	}

	pageTokenLimit := 1 + (maxPageToken-1)/req.GetPageSize()

	//If pageToken will be last nextPageToken will be equal to previous pageToken
	if int32(selectedPage) > pageTokenLimit {
		return nil, errors.New("Invalid page token")
	} else if int32(selectedPage) == pageTokenLimit {
		listUsers.NextPageToken = fmt.Sprintf("%d", selectedPage)
	}

	//First page = pageToken * pageSize - pageSize
	//Last page = pageToken * pageSize
	rows, err := r.DB.QueryContext(ctx, sqlStr, (int32(selectedPage)*req.GetPageSize())-req.GetPageSize(), int32(selectedPage)*req.GetPageSize())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var user pb.User

		if err := rows.Scan(&user.Id, &user.Email); err != nil {
			return nil, err
		}
		listUsers.Users = append(listUsers.Users, &user)
	}

	return &listUsers, nil
}

//GetUser - get user with gettuserrequest
//In futurue password checking need to be done
func (r *PsqlRepository) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	var newUser pb.User
	err := r.DB.QueryRowContext(ctx, "select * from gousers where id = $1", req.GetId()).Scan(&newUser.Id, &newUser.Email)

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("Not found")
	case err != nil:
		return nil, err
	default:
		return &newUser, nil
	}
}

//CreateUser - create user
func (r *PsqlRepository) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	var newUser pb.User

	hasedPass, err := utils.HashPassword([]byte(req.User.GetPassword()))

	if err != nil {
		return nil, err
	}

	err = r.DB.QueryRowContext(ctx, "INSERT INTO gousers(email, password) values($1, $2) returning id, email", req.User.GetEmail(), hasedPass).Scan(&newUser.Id, &newUser.Email)

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("Not found")
	case err != nil:
		return nil, err
	default:
		return &newUser, nil
	}
}

//UpdateUser - update user with given id
func (r *PsqlRepository) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	var newUser pb.User
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	hasedPass, err := utils.HashPassword([]byte(req.User.GetPassword()))

	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(ctx, "update gousers set email = $1, password = $2 where id = $3 returning id, name", req.User.GetEmail(), hasedPass, req.User.GetId()).Scan(&newUser.Id, &newUser.Email)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("Not found")
	case err != nil:
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("update failed: %v, unable to rollback: %v", err, rollbackErr)
		}
		return nil, fmt.Errorf("update failed: %v", err)
	default:
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &newUser, nil
	}
}

//DeleteUser - delete user from db by id
func (r *PsqlRepository) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*empty.Empty, error) {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	_, execErr := tx.ExecContext(ctx, "Delete from gousers WHERE id = $1", req.GetId())
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("delete failed: %v, unable to rollback: %v", err, rollbackErr)
		}
		return nil, fmt.Errorf("delete failed: %v", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (r *PsqlRepository) PasswordCheck(ctx context.Context, req *pb.SessionAuthUserRequest) (bool, error) {
	var hashpassword []byte
	err := r.DB.QueryRowContext(ctx, "Select id, password from gousers where email = $1", req.User.GetEmail()).Scan(&req.User.Id, &hashpassword)

	//Compering passwords
	ok, _ := utils.VerifyPassword([]byte(req.User.GetPassword()), hashpassword)

	switch {
	case err == sql.ErrNoRows:
		return false, errors.New("Not found")
	case err != nil:
		return false, err
	default:
		if !ok {
			return false, errors.New("Not found")
		}
		return true, nil
	}
}
