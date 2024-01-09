package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// likePost allows a user to put like to a post.
// If the user in not authorized, the request will fail.
// If the user id doesn't exist, the request will fail.
// If the post id doesn't exist, the request will fail.
// If the request is OK, it will return User{} object who has put the like.
// Note: a user can put like to his own post
func (rt *_router) likePost(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
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
		context.Logger.Error("Error retrieving the current uid that makes put like request")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// check if the user exists
	check, err := rt.db.CheckExistsByUID(uid)
	if err != nil {
		context.Logger.Error("Error retrieving information on UID for putting like request!")
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in putting like request! User doesn't exist!")
		http.Error(w, "User seems not exist.", http.StatusNotFound)
		return
	}

	// The Post ID in the path is a 64-bit unsigned integer. Let's parse it.
	postid, err := strconv.ParseUint(params.ByName("postid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing postid in putting like request")
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
		context.Logger.Error("Error retrieving information on postid for putting like!\nDetail: ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in putting like request! Post doesn't exist")
		http.Error(w, "Post seems not exist.", http.StatusNotFound)
		return
	}

	// Put the like on table
	err = rt.db.LikePost(postid, uid)

	if err != nil {
		context.Logger.Error("Error adding like to table.\nDetail: ", err.Error())
		http.Error(w, "Something wrong adding like to the post.", http.StatusInternalServerError)
		return
	}

	// Like correctly added
	// Now return User{} struct
	userdb, err := rt.db.GetUserByID(uid)

	if err != nil {
		context.Logger.Error("Error retrieving user structure in like request.\nDetail: ", err.Error())
		http.Error(w, "Something wrong adding like to the post.", http.StatusInternalServerError)
		return
	}

	var userapi User
	err = userapi.FromDatabase(userdb)

	if err != nil {
		context.Logger.Error("Error parsing using structure in like request.\nDetail: ", err.Error())
		http.Error(w, "Something wrong adding like to the post.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userapi)

}
