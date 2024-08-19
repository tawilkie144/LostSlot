package Controllers

import (
	"LostSlot/src/Entities"
	"LostSlot/src/Services"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type UserController struct {
	service *Services.UserService
}

var this UserController

func init() {
	this = UserController{service: &Services.UserService{}}
}

func GetUsersByPage(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("GetUsersByPage", writer, request)
}

func GetUser(writer http.ResponseWriter, request *http.Request) {
	stringId := chi.URLParam(request, "userId")
	id, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		//malformed request
		http.Error(writer, "400", http.StatusBadRequest)
		return
	}
	userCollection, err := this.service.GetUsers(id, 1)
	if err != nil {
		//some other error
		http.Error(writer, "500", http.StatusInternalServerError)
		return
	}
	if len(userCollection) == 0 {
		//No user
		http.Error(writer, "404", http.StatusNotFound)
		return
	}
	rUser := userCollection[0]
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(&rUser)
}

func CreateUser(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("CreateUser", writer, request)
}

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("DeleteUser", writer, request)
}

func GetUsersByPage2(page int, pageSize int) ([]Entities.User, error) {
	startId := int64(page - 1*pageSize)
	rUsers, err := this.service.GetUsers(startId, pageSize)
	if err != nil {
		return nil, err
	}
	if len(rUsers) == 0 {
		return nil, fmt.Errorf("error: no users found")
	}

	return rUsers, nil
}

func GetUser2(id int64) (*Entities.User, error) {
	rUser, err := this.service.GetUsers(id, 1)
	if err != nil {
		return nil, err
	}
	if len(rUser) == 0 {
		return nil, fmt.Errorf("404")
	}
	return &rUser[0], nil
}

func NewUser2(user *Entities.User) error {
	err := this.service.NewUser(user)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser2(id int64) error {
	err := this.service.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
