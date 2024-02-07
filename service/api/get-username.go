package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// getUsername allows getting username of a user from his uid
func (rt *_router) getUsername(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)
	if err != nil {
		// The value was not uint64, reject the action indicating an error on the client side.
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authId, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)
	if err != nil {
		// The value was not uint64, reject the action indicating an error on the client side.
		context.Logger.Error()
		context.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if the user is authorized to perform the action
	if authId != context.Uid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username, err := rt.db.GetUsername(uid)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result := map[string]string{
		"username": username,
	}

	w.Header().Set("Content-Type", "application-json")
	_ = json.NewEncoder(w).Encode(result)
}
