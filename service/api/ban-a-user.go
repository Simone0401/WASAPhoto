package api

import (
	"encoding/json"
	"fmt"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// banUser allows to the specified uid user to ban another specified user.
// If the user to mute is already banned, nothing change.
// If the uid doesn't exist, the request will fail.
// If the banned user id doesn't exist, the request will fail.
// If the user in not authorized, the request will fail.
// User cannot ban himself. The request will fail.
// If the request is correct, it will return 201 status and the muted_user_info{} object
// Note: if muteduid user has followed uid and the ban request is ok, the follow will be removed
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing uid for banning request", err.Error())
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": err.Error(),
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// The ban User ID in the path is a 64-bit unsigned integer. Let's parse it.
	muteduid, err := strconv.ParseUint(params.ByName("muteduid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing muteduid for banning request", err.Error())
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
		context.Logger.Error("User is not authorizated to ban")
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

	// check if the muteduid exists
	_, err = rt.db.GetUserByID(muteduid)
	if err != nil {
		context.Logger.Error("The user with specified muteduid seems not exist ", err.Error())
		http.Error(w, "The user with specified muteduid seems not exist", http.StatusNotFound)
		return
	}

	// check if the user is trying to mute himself
	if uid == muteduid {
		context.Logger.Error("User is trying to mute himself")
		http.Error(w, "User cannot mute himself", http.StatusBadRequest)
		return
	}

	// check if the specified muteduid has not banned uid user
	isBanned, err := rt.db.HasBanned(muteduid, uid)
	if err != nil {
		context.Logger.Error("Something wrong retrieving ban information ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if isBanned {
		errorMessage := fmt.Sprintf("User %d has banned %d user, already muted.", muteduid, uid)
		context.Logger.Error(errorMessage)
		http.Error(w, "Specified muteduid seems not exist", http.StatusNotFound)
		return
	}

	// check if the uid user already muted the muteduid user
	isMuted, err := rt.db.HasMuted(uid, muteduid)
	if err != nil {
		context.Logger.Error("Something wrong retrieving ban information ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	if !isMuted {
		// all the checks are performed, the uid user is able to mute muteduid user
		// Ban action involves unfollow action
		_, err := rt.db.BanUser(uid, muteduid)
		if err != nil {
			context.Logger.Error("Something wrong banning user\nDetail: ", err.Error())
			http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
			return
		}

		_, err = rt.db.UnfollowUser(muteduid, uid)
		if err != nil {
			context.Logger.Error("Something wrong banning user\nDetail: ", err.Error())
			http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	var user User
	userdb, err := rt.db.GetUserByID(muteduid)
	if err != nil {
		context.Logger.Error("Error retrieving information about just banned user\nDetail: ", err.Error())
		http.Error(w, "Error during mute action", http.StatusInternalServerError)
		return
	}

	err = user.FromDatabase(userdb)

	if err != nil {
		context.Logger.Error("Error converting userdb struct to user API struct\nDetail: ", err.Error())
		http.Error(w, "Error during mute action", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]User{"muted_user_info": user})
}
