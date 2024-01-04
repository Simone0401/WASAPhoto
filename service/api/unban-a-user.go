package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// unbanUser allows to the specified uid user to unban another specified user.
// If the user to mute is already unbanned, nothing change.
// If the uid doesn't exist, the request will fail.
// If the banned user id doesn't exist, the request will fail.
// If the user in not authorized, the request will fail.
// User cannot unban himself. The request will fail.
// If the request is correct, it will return 204 status
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing uid for unbanning request", err.Error())
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
		context.Logger.Error("Error parsing muteduid for unbanning request", err.Error())
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

	// check if the user is trying to unmute himself
	if uid == muteduid {
		context.Logger.Error("User is trying to unmute himself")
		http.Error(w, "User cannot mute himself", http.StatusBadRequest)
		return
	}

	// check if the uid user already unmuted the muteduid user
	isMuted, err := rt.db.HasMuted(uid, muteduid)
	if err != nil {
		context.Logger.Error("Something wrong retrieving ban information ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
	}

	if isMuted {
		// all the checks are performed, the uid user is able to unmute muteduid user
		_, err := rt.db.UnbanUser(uid, muteduid)
		if err != nil {
			context.Logger.Error("Something wrong unbanning user\nDetail: ", err.Error())
			w.Header().Add("Content-Type", "application/json")
			http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
			return
		}

	}

	w.WriteHeader(http.StatusNoContent)
}
