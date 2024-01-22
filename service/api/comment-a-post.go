package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/Simone0401/WASAPhoto/service/database"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
)

// commentPost allows a user to add a comment under a post.
// If the user in not authorized, the request will fail.
// If the user id doesn't exist, the request will fail.
// If the post id doesn't exist, the request will fail.
// If the request is OK, it will return Comment{} object.
func (rt *_router) commentPost(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {

	comment := map[string]Comment{
		"comment": {},
	}

	// Read the body content and parse it into Comment{} struct
	bodyContent, err := io.ReadAll(r.Body)
	if err != nil {
		context.Logger.Error("Error retrieving request body in comment request.\nDetail: ", err.Error())
		http.Error(w, "Something wrong uploading your comment", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bodyContent, &comment)

	if err != nil {
		context.Logger.Error("Error parsing comment into structure.\nDetail: ", err.Error())
		http.Error(w, "Something wrong uploading your comment", http.StatusInternalServerError)
		return
	}

	commentApi := comment["comment"]
	uid := commentApi.Userid

	// check if the Bearer Authorization Token is set
	if !rt.isAuthorized(r.Header) {
		context.Logger.Error("The bearer format token is not valid in adding comment request!")
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
		context.Logger.Error("Error retrieving the current uid that makes comment post request")
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

	// Check message validity
	if !commentApi.IsValid() {
		context.Logger.Error("Content message for comment is not valid!")
		http.Error(w, "Your message cannot be uploaded. Check its format!", http.StatusBadRequest)
		return
	}

	// Insert comment into table
	commentApi.Postid = postid
	var commentDb database.Comment
	commentDb = commentApi.ToDatabase()
	commentDb, err = rt.db.AddComment(commentDb.Userid, commentDb.Postid, commentDb.Message)

	if err != nil {
		context.Logger.Error("Error inserting comment into tables\nDetail: ", err.Error())
		http.Error(w, "Something wrong uploading your comment", http.StatusInternalServerError)
		return
	}

	// Message correctly inserted
	err = commentApi.FromDatabase(commentDb)

	if err != nil {
		context.Logger.Error("Error converting CommentDb to CommentAPI\nDetail: ", err.Error())
		http.Error(w, "Something wrong uploading your comment", http.StatusInternalServerError)
		return
	}

	// Change DateTime format
	commentApi.Datetime, err = formatDatetime(commentApi.Datetime)
	if err != nil {
		context.Logger.Warning("Error parsing datetime in adding comment request!\nDetail: ", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	comment["comment"] = commentApi
	_ = json.NewEncoder(w).Encode(comment)
}
