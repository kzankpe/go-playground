package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type user struct {
	ID        string `json: "ID"`
	LastName  string `json: "LastName"`
	FirstName string `json: "FirstName"`
	Type      string `json: "Type"`
}

type allUsers []user

var users = allUsers{
	{
		ID:        "1",
		LastName:  "Client",
		FirstName: "Jon",
		Type:      "Manager",
	},
	{
		ID:        "2",
		LastName:  "Sub",
		FirstName: "Client",
		Type:      "User",
	},
	{
		ID:        "3",
		LastName:  "Tech",
		FirstName: "Nician",
		Type:      "Contractor",
	},
	{
		ID:        "4",
		LastName:  "Super",
		FirstName: "Power",
		Type:      "Admin",
	},
}

func HomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "HomePage !!")
}

func getOneUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	for _, singleUser := range users {
		if singleUser.ID == userID {
			json.NewEncoder(w).Encode(singleUser)
		}
	}
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	var updatedUser user
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "Please Enter correct informations to update the user")
	}

	json.Unmarshal(reqBody, &updatedUser)
	for i, singleUser := range users {
		if singleUser.ID == userID {
			singleUser.LastName = updatedUser.LastName
			singleUser.FirstName = updatedUser.FirstName
			singleUser.Type = updatedUser.Type
			users = append(users[:i], singleUser)
			json.NewEncoder(w).Encode(singleUser)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HomeLink)
	router.HandleFunc("/users/{id}", getOneUser).Methods(http.MethodGet)
	router.HandleFunc("/users", getAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", updateUser).Methods(http.MethodPatch)
	log.Fatal(http.ListenAndServe(":8080", router))
}
