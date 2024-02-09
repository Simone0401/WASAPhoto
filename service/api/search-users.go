package api

import (
	"encoding/json"
	"fmt"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// getUsers allows searching for users' profile information passing a username (or part of it).
// If the user is not authorized, the request will fail.
// If the user id doesn't exist, the request will fail.
// The return values will be the user's profile information for each user that matches the search criteria.
func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {

	// check if the Bearer Authorization Token is set
	if !rt.isAuthorized(r.Header) {
		context.Logger.Error("The bearer format token is not valid in getting profile request!")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		response := map[string]string{
			"error": "log to do action",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the current user is authorized
	check, err := rt.db.CheckExistsByUID(context.Uid)
	if err != nil {
		context.Logger.Error("Error retrieving information on UID that makes getting profile request!")
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in getting profile request! User that makes request doesn't exist!")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// The User username (or part of it) in the query is a string. Let's parse it.
	username := r.URL.Query().Get("search")

	// Prepare return statement
	var users []User

	usersDb, err := rt.db.SearchUserByUsername(username)
	if err != nil {
		context.Logger.Error("Error getting []Users from username in getting profiles request!\nDetail: ", err.Error())
		http.Error(w, "Something wrong retrieving user profiles", http.StatusInternalServerError)
		return
	}

	// Append each user to the list of users
	for i, user := range usersDb {

		// Check if the user is banned
		banned, err := rt.db.HasBanned(user.Userid, context.Uid)
		if err != nil {
			mess := fmt.Sprintf("Error getting ban information for the user number %d in getting users request\nDetail: ", i)
			context.Logger.Error(mess, err.Error())
			http.Error(w, "Something wrong retrieving profiles", http.StatusInternalServerError)
			return
		}

		if !banned {
			var userAPI User
			err = userAPI.FromDatabase(user)
			if err != nil {
				mess := fmt.Sprintf("Error parsing userDB to userAPI for user number %d in getting users request\nDetail: ", i)
				context.Logger.Error(mess, err.Error())
				http.Error(w, "Something wrong retrieving profiles", http.StatusInternalServerError)
				return
			}
			users = append(users, userAPI)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&map[string]interface {
	}{"users": users})

}
