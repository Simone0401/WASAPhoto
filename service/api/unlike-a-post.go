package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// unlikePost allows a user to remove like from a post.
// If the user in not authorized, the request will fail.
// If the user id doesn't exist, the request will fail.
// If the post id doesn't exist, the request will fail.
// If the request is OK, it will return 204 code status.
// Note: a user can remove like from his own post
func (rt *_router) unlikePost(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing uid in put like request")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "not correct format for uid",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the Bearer Authorization Token is set
	if !rt.isAuthorized(r.Header) {
		context.Logger.Error("The bearer format token is not valid in deleting put like request!")
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
		context.Logger.Error("Error retrieving the current uid that makes delete like request")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// check if the user exists
	check, err := rt.db.CheckExistsByUID(uid)
	if err != nil {
		context.Logger.Error("Error retrieving information on UID for deleting like request!")
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in deleting like request! User doesn't exist!")
		http.Error(w, "User seems not exist.", http.StatusNotFound)
		return
	}

	// The Post ID in the path is a 64-bit unsigned integer. Let's parse it.
	postid, err := strconv.ParseUint(params.ByName("postid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing postid in deleting like request")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "not correct format for postid",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the post exists
	check, err = rt.db.CheckPostByPostid(postid)
	if err != nil {
		context.Logger.Error("Error retrieving information on postid for deleting like!\nDetail: ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in deleting like request! Post doesn't exist")
		http.Error(w, "Post seems not exist.", http.StatusNotFound)
		return
	}

	// Delete the like from table
	err = rt.db.UnlikePost(postid, uid)

	if err != nil {
		context.Logger.Error("Error removing like from table.\nDetail: ", err.Error())
		http.Error(w, "Something wrong removing like from the post.", http.StatusInternalServerError)
		return
	}

	// Like correctly removed
	// Now return 204 status
	w.WriteHeader(http.StatusNoContent)
}
