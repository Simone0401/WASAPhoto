package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// uncommentPost allows a user to remove an own comment under a post.
// If the user in not authorized, the request will fail.
// If the user id doesn't exist, the request will fail.
// If the post id doesn't exist, the request will fail.
// If the request is OK, it will return 204 scode tatus.
func (rt *_router) uncommentPost(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {

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

	// check if the user exists
	currentUid := context.Uid
	check, err := rt.db.CheckExistsByUID(currentUid)
	if err != nil {
		context.Logger.Error("Error retrieving information on UID for putting like request!")
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}
	if !check {
		context.Logger.Error("Error in removing comment request! User doesn't exist!")
		http.Error(w, "User seems not exist.", http.StatusNotFound)
		return
	}

	// The Post ID in the path is a 64-bit unsigned integer. Let's parse it.
	postid, err := strconv.ParseUint(params.ByName("postid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing postid in deleting comment request")
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
		context.Logger.Error("Error retrieving information on postid for deleting comment!\nDetail: ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in deleting comment request! Post doesn't exist")
		http.Error(w, "Post seems not exist.", http.StatusNotFound)
		return
	}

	// The Comment ID in the path is a 64-bit unsigned integer. Let's parse it.
	commentid, err := strconv.ParseUint(params.ByName("commentid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing comment in deleting comment request")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "not correct format for commentid",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the comment exists
	check, err = rt.db.CheckCommentByCommentid(commentid)
	if err != nil {
		context.Logger.Error("Error retrieving information on commentid for deleting comment!\nDetail: ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in deleting comment request! Comment doesn't exist")
		http.Error(w, "Comment seems not exist.", http.StatusNotFound)
		return
	}

	// check if the current user is authorized
	if check, err = rt.db.CheckCommentOwner(commentid, currentUid); err != nil {
		context.Logger.Error("Error checking comment owner in deleting comment request\nDetail: ", err.Error())
		http.Error(w, "Something wrong in deleting comment request", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("User is not the owner of comment in delete comment request")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Delete comment from table
	err = rt.db.DeleteComment(commentid)

	if err != nil {
		context.Logger.Error("Error deleting comment from table\nDetail: ", err.Error())
		http.Error(w, "Something wrong deleting your comment", http.StatusInternalServerError)
		return
	}

	// Comment correctly removed
	w.WriteHeader(http.StatusNoContent)

}
