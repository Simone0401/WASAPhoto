package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) setUsername(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing uid")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "not correct format",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the Bearer Authorization Token is set
	if !rt.isAuthorized(r.Header) {
		context.Logger.Error("The bearer format token is not valid!")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		response := map[string]string{
			"error": "log to do action",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the current user is authorized
	currentUid := context.Uid
	if currentUid != uid {
		context.Logger.Error("Error retrieving the current uid that makes request")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var user User

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// The body was not a parseable JSON, reject it
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !user.IsValid() {
		// Here we validated the user structure content (e.g., username has correct format)
		// Note: the IsValid() function skips the ID check.
		context.Logger.Error("Error, User structure is not valid!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.Userid = uid

	userdb := user.ToDatabase()
	// try change user username, if username is already taken the request will fail
	err = rt.db.SetUsername(userdb.Userid, user.Username)
	if err != nil {
		context.Logger.Error("Error setting the new username, already taken")
		http.Error(w, "Username already taken. Username must be unique", http.StatusBadRequest)
		return
	}

	returnUser, err := rt.db.GetUserByID(userdb.Userid)
	if err != nil {
		context.Logger.Error("Error retrieving the update user: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// convert database User to API User
	err = user.FromDatabase(returnUser)
	if err != nil {
		context.Logger.Error("Error parsing User Database to User API")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if the new username has been set, reply a success message and
	// return the new user object
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		context.Logger.Error("Error encoding User API to JSON Object")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
