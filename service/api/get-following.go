package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/Simone0401/WASAPhoto/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// getFollowing allows to get a specif following.
// If the request following user doesn't follow uid user, 404 error will be reported.
// If the user id doesn't exist, the request will fail.
// If the following user id doesn't exist, the request will fail.
// If the user in not authorized, the request will fail.
// If the request is valid, it will return the User{} object about followed user
func (rt *_router) getFollowing(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing uid for getting followed request", err.Error())
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "Something wrong in the server",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// The Followed User ID in the path is a 64-bit unsigned integer. Let's parse it.
	fuid, err := strconv.ParseUint(params.ByName("fuid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing fuid in getting followed request", err.Error())
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "Something wrong in the server",
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

	// check if the uid user has followed the fuid user
	isFollowed, err := rt.db.HasFollowed(uid, fuid)
	if err != nil {
		context.Logger.Error("Something wrong retrieving followed information ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
	}

	if !isFollowed {
		context.Logger.Info("The user with specified uid doesn't follow the user with specified fuid")
		http.Error(w, "uid doesn't follow fuid", http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var user User
	err = user.FromDatabase(fuserdb)
	if err != nil {
		context.Logger.Error("Error parsing fuid user from DB ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(fuserdb)
}
