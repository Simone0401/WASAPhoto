package api

import (
	"encoding/json"
	"fmt"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/Simone0401/WASAPhoto/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// followUser allows a user to follow another user.
// If the followed user is already followed, nothing change.
// If the user id doesn't exist, the request will fail.
// If the followed user id doesn't exist, the request will fail.
// If the user in not authorized, the request will fail.
// User cannot follow himself.
// If the request is valid, it will return the User{} object about just followed user
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing uid for follow request", err.Error())
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": err.Error(),
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// The Follow User ID in the path is a 64-bit unsigned integer. Let's parse it.
	fuid, err := strconv.ParseUint(params.ByName("fuid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing fuid for follow request", err.Error())
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": err.Error(),
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
			"error": "login to perform the action",
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

	// check if the uid exists
	_, err = rt.db.GetUserByID(uid)
	if err != nil {
		context.Logger.Error("The user with specified uid seems not exist ", err.Error())
		http.Error(w, "The user with specified uid seems not exist ", http.StatusNotFound)
		return
	}

	// check if the fuid exists
	var fuserdb database.User
	fuserdb, err = rt.db.GetUserByID(fuid)
	if err != nil {
		context.Logger.Error("The user with specified fuid seems not exist ", err.Error())
		http.Error(w, "The user with specified fuid seems not exist", http.StatusNotFound)
		return
	}

	// check if the user is trying to follow himself
	if uid == fuid {
		context.Logger.Error("User is trying to follow himself")
		http.Error(w, "User cannot follow himself", http.StatusBadRequest)
		return
	}

	// check if the specified fuid has not banned uid user
	isBanned, err := rt.db.HasBanned(fuid, uid)
	if err != nil {
		context.Logger.Error("Something wrong retrieving ban information ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if isBanned {
		errorMessage := fmt.Sprintf("User %d has banned %d user. Cannot follow.", fuid, uid)
		context.Logger.Error(errorMessage)
		http.Error(w, "Cannot follow the user", http.StatusForbidden)
		return
	}

	// check if the uid user already follow the fuid user
	isFollowed, err := rt.db.HasFollowed(uid, fuid)
	if err != nil {
		context.Logger.Error("Something wrong retrieving follow information ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	if !isFollowed {
		// all the checks are performed, the uid user is able to follow fuid user
		_, err := rt.db.FollowUser(uid, fuid)
		if err != nil {
			context.Logger.Error("Something wrong following user\nDetail: ", err.Error())
			http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	var user User
	err = user.FromDatabase(fuserdb)
	if err != nil {
		context.Logger.Error("Error parsing fuid user from DB ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(fuserdb)

}
