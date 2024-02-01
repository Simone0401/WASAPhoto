package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// getLike allows to check if a user has put like to a post.
// If the user in not authorized, the request will fail.
// If the user id doesn't exist, the request will fail.
// If the post id doesn't exist, the request will fail.
// If the request is OK, it will return Like resource.
func (rt *_router) getLike(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing uid in getting like request")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "not correct format for uid",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the Bearer Authorization Token is set
	if !rt.isAuthorized(r.Header) {
		context.Logger.Error("The bearer format token is not valid in getting like request!")
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
		context.Logger.Error("Error retrieving the current uid that makes getting like request")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// check if the user exists
	check, err := rt.db.CheckExistsByUID(uid)
	if err != nil {
		context.Logger.Error("Error retrieving information on UID for getting like request!")
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in getting like request! User doesn't exist!")
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
		context.Logger.Error("Error retrieving information on postid for getting like!\nDetail: ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in getting like request! Post doesn't exist")
		http.Error(w, "Post seems not exist.", http.StatusNotFound)
		return
	}

	// Check like on table
	check, err = rt.db.CheckLike(postid, uid)

	if err != nil {
		context.Logger.Error("Error getting like from table.\nDetail: ", err.Error())
		http.Error(w, "Something wrong getting likes from the post.", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Info("Like doesn't found on post!")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		response := map[string]interface{}{
			"errno": 1,
			"error": "like doesn't found on post",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// Like status correctly recovered
	// Now return Like resource to the client
	likeapi := map[string]uint64{
		"user_id": uid,
		"postid":  postid,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(likeapi)

}
