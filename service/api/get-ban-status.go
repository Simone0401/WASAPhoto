package api

import (
	"encoding/json"
	"fmt"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// getMuted Allows getting information on ban status for a specific uid relatives to a ban uid.
// If the user id doesn't exist, the request will fail.
// If the banned user id doesn't exist, the request will fail.
// If the user in not authorized, the request will fail.
// If the request is correct, it will return 200 status and the muted_user_info{} object
func (rt *_router) getMuted(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing uid for getting ban request", err.Error())
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "Something wrong in the server",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// The ban User ID in the path is a 64-bit unsigned integer. Let's parse it.
	muteduid, err := strconv.ParseUint(params.ByName("muteduid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing muteduid for getting ban request", err.Error())
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
		context.Logger.Error("User is not authorizated in getting ban status request")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// check if the uid exists
	_, err = rt.db.GetUserByID(uid)
	if err != nil {
		context.Logger.Error("The user with specified uid seems not exist in getting ban request ", err.Error())
		http.Error(w, "The user with specified uid seems not exist ", http.StatusNotFound)
		return
	}

	// check if the muteduid exists
	_, err = rt.db.GetUserByID(muteduid)
	if err != nil {
		context.Logger.Error("The user with specified muteduid seems not exist in getting ban request ", err.Error())
		http.Error(w, "The user with specified muteduid seems not exist", http.StatusNotFound)
		return
	}

	// check if the specified uid has banned muteduid user
	isBanned, err := rt.db.HasBanned(uid, muteduid)
	if err != nil {
		context.Logger.Error("Something wrong retrieving ban information in getting ban request", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !isBanned {
		infoMessage := fmt.Sprintf("User %d hasn't banned %d user.", uid, muteduid)
		context.Logger.Info(infoMessage)
		http.Error(w, "Specified muteduid seems not exist", http.StatusNotFound)
		return
	}

	// Get out the user information
	userDB, err := rt.db.GetUserByID(muteduid)

	if err != nil {
		context.Logger.Error("Error retrieving information about banned user\nDetail: ", err.Error())
		http.Error(w, "Error during getting mute action", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var user User
	err = user.FromDatabase(userDB)

	if err != nil {
		context.Logger.Error("Error converting userdb struct to user API struct\nDetail: ", err.Error())
		http.Error(w, "Error during getting mute action", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]User{"muted_user_info": user})
}
