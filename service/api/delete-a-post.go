package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// deletePost allows deleting a post by a user, if he is the post author.
// If the user id doesn't exist, the request will fail.
// If the user in not authorized, the request will fail.
// If the post id doesn't exist, the request will fail.
// The request will remove all the comments and the image, too.
func (rt *_router) deletePost(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing uid in deleting post request")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "not correct format for uid",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the Bearer Authorization Token is set
	if !rt.isAuthorized(r.Header) {
		context.Logger.Error("The bearer format token is not valid in deleting post request!")
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
		context.Logger.Error("Error retrieving the current uid that makes uploading post request")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// check if the user exists
	check, err := rt.db.CheckExistsByUID(uid)
	if err != nil {
		context.Logger.Error("Error retrieving information on UID for deleting post!")
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in deleting post request! User doesn't exist")
		http.Error(w, "User seems not exist", http.StatusNotFound)
		return
	}

	// The Post ID in the path is a 64-bit unsigned integer. Let's parse it.
	postid, err := strconv.ParseUint(params.ByName("postid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing postid in deleting post request")
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
		context.Logger.Error("Error retrieving information on postid for deleting post!\nDetail: ", err.Error())
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in deleting post request! Post doesn't exist")
		http.Error(w, "Post seems not exist", http.StatusNotFound)
		return
	}

	// Remove all comments under the post
	err = rt.db.RemoveCommentsFromPost(postid)
	if err != nil {
		context.Logger.Info("Error removing comments under post\nDetail: ", err.Error())
		http.Error(w, "Something wrong deleting the post", http.StatusInternalServerError)
		return
	}

	// Remove all likes from the post
	err = rt.db.RemoveLikesFromPost(postid)
	if err != nil {
		context.Logger.Info("Error removing likes under post\nDetail: ", err.Error())
		http.Error(w, "Something wrong deleting the post", http.StatusInternalServerError)
		return
	}

	// Remove the image from folder
	// First try with .png, if not exists try with .jpeg
	pngPath := strconv.FormatUint(postid, 10) + ".png"
	err = deleteImage("media/img/", pngPath)
	if err != nil {
		jpegPath := strconv.FormatUint(postid, 10) + ".jpeg"
		err = deleteImage("media/img/", jpegPath)
		if err != nil {
			context.Logger.Error("Error removing post!\nDetail: ", err.Error())
			http.Error(w, "Something wrong removing post", http.StatusInternalServerError)
			return
		}
	}

	// Remove post from table
	err = rt.db.RemovePost(postid, uid)
	if err != nil {
		context.Logger.Error("Cannot delete post from table.\nDetail: ", err.Error())
		http.Error(w, "Something wrong removing post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
