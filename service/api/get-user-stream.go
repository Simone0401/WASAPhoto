package api

import (
	"encoding/json"
	"fmt"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// getMyStream allows getting user stream photos passing the uid.
// If the user in not authorized, the request will fail.
// If the user id doesn't exist, the request will fail.
// If the post id doesn't exist, the request will fail.
// The stream consists in an array of post. (Check API documentation for detail)
// Note: for getting a binary image it's necessary using the 'Get Image API'
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing uid in get stream request")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "not correct format for uid",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the Bearer Authorization Token is set
	if !rt.isAuthorized(r.Header) {
		context.Logger.Error("The bearer format token is not valid in getting stream request!")
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
		context.Logger.Error("Error retrieving the current uid that makes getting stream request")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// check if the user exists
	check, err := rt.db.CheckExistsByUID(uid)
	if err != nil {
		context.Logger.Error("Error retrieving information on UID for getting stream request!")
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in getting stream request! User doesn't exist!")
		http.Error(w, "User seems not exist.", http.StatusNotFound)
		return
	}

	// Prepare return struct
	posts := map[string][]Post{
		"posts": {},
	}

	// Get the stream
	listPost, err := rt.db.GetUserStream(uid)
	if err != nil {
		context.Logger.Error("Error retrieving post for user during getting stream request!\nDetail: ", err.Error())
		http.Error(w, "Something wrong retrieving your stream", http.StatusInternalServerError)
		return
	}

	// Append each post to the list
	for i, post := range listPost {
		var postAPI Post
		err = postAPI.FromDatabase(post)
		if err != nil {
			mess := fmt.Sprintf("Error parsing postDB to postAPI for post number %d\nDetail: ", i)
			context.Logger.Error(mess, err.Error())
			http.Error(w, "Something wrong retrieving your stream", http.StatusInternalServerError)
			return
		}
		// Change datetime format for each comment
		for i := 0; i < len(postAPI.Comments); i++ {
			postAPI.Comments[i].Datetime, _ = formatDatetime(postAPI.Comments[i].Datetime)
		}
		postAPI.Datetime, _ = formatDatetime(postAPI.Datetime)
		posts["posts"] = append(posts["posts"], postAPI)
	}

	// Prepare the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(posts)
}
