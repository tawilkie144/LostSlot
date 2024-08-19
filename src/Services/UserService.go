package Services

import (
	"context"
	"fmt"
	"log"
	"LostSlot/src/Data"
	"LostSlot/src/Data/postgres"
	"LostSlot/src/Entities"
	"strconv"
)

type UserService struct {
	userData Data.DataStore
	config   Entities.Config
}

func NewUserService() *UserService{
	rService := UserService{}
	config, err := Entities.GetConfig()
	if err != nil {
		log.Fatalf("error: could not retrieve config when initializing UserService: %v", err)
	}
	rService.config = config
	dataSource := config.UserServiceDataSource
	rService.userData = 
}

func (service *UserService) GetUsers(startNum int64, NumReturn int) ([]Entities.User, error) {
	conn, err := service.userData.NewConnection(config.PostgresHost, strconv.Itoa(config.PostgresPort), config.PostgresDB, config.PostgresUser, config.PostgresPassword)
	if err != nil {
		return nil, fmt.Errorf("error: could not connect to database: %v", err)
	}
	rows, err := conn.Query(context.TODO(),
		"SELECT * FROM mercury.app_user WHERE user_id >= $1 ORDER BY user_id LIMIT $2",
		map[string]any{"1": startNum, "2": NumReturn})
	if err != nil {
		return nil, err
	}
	defer (*rows).Close()

	rRows := make([]Entities.User, 1, NumReturn)
	i := 0
	for rows.Next() {
		i++
		var thisUser Entities.User
		err = rows.Scan(&thisUser.UserId, &thisUser.GivenName, &thisUser.LastName, &thisUser.Initials, &thisUser.PreferredName, &thisUser.HashedPassword, &thisUser.CreatedAt, &thisUser.UpdatedAt)
		if err != nil {
			log.Printf("error scanning row %d: %s", i, err)
		}

		rRows = append(rRows, thisUser)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error: error while querying database at row: %v", err)
	}
	return rRows, nil
}

func (service *UserService) NewUser(user *Entities.User) error {
	if user.PreferredName == "" || user.GivenName == "" || user.LastName == "" || user.HashedPassword == "" {
		return fmt.Errorf("error: malformed input")
	}
	config, _ := Entities.GetConfig() //ignoring error for now. I think if we can't retrieve the config, we should just crash
	conn, err := postgres.NewConnection(config.PostgresHost, strconv.Itoa(config.PostgresPort), config.PostgresDB, config.PostgresUser, config.PostgresPassword)
	if err != nil {
		return fmt.Errorf("error: could not connect to database: %v", err)
	}
	_, err = conn.Query(context.TODO(),
		"INSERT INTO mercury.app_user (given_name, last_name, initials, preferred_name, hashed_password) VALUES ($1, $2, $3, $4, $5)",
		map[string]any{"1": user.GivenName, "2": user.LastName, "3": user.Initials, "4": user.PreferredName, "5": user.HashedPassword})
	if err != nil {
		return err
	}
	return nil
}

func (service *UserService) DeleteUser(id int64) error {
	config, _ := Entities.GetConfig() //ignoring error for now. I think if we can't retrieve the config, we should just crash
	conn, err := postgres.NewConnection(config.PostgresHost, strconv.Itoa(config.PostgresPort), config.PostgresDB, config.PostgresUser, config.PostgresPassword)
	if err != nil {
		return fmt.Errorf("error: could not connect to database: %v", err)
	}
	_, err = conn.Query(context.TODO(),
		"DELETE FROM mercury.app_user WHERE user_id = $1",
		map[string]any{"1": id})
	if err != nil {
		return err
	}
	return nil
}
