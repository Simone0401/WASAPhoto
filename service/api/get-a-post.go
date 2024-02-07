package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// getPost allows recovering a post passing its post id.
// If the post id doesn't exist, the request will fail.
// If the user is not authorized, the request will fail.
func (rt *_router) getPost(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The post ID in the path is a 64-bit unsigned integer. Let's parse it.
	postid, err := strconv.ParseUint(params.ByName("postid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing postid in get post request.")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "not correct format for postid",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the Bearer Authorization Token is set
	if !rt.isAuthorized(r.Header) {
		context.Logger.Error("The bearer format token is not valid for getting a post!")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		response := map[string]string{
			"error": "log to do action",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the post exists
	check, err := rt.db.CheckPostByPostid(postid)

	if err != nil {
		context.Logger.Error("Something wrong checking postid\nDetail: ", err.Error())
		http.Error(w, "Something wrong", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Postid requested doesn't exist")
		http.Error(w, "Post seems not exist", http.StatusNotFound)
		return
	}

	// Post exists, prepare return statement
	var PostAPI Post

	// Recover post data
	postDB, err := rt.db.GetPost(postid)
	if err != nil {
		context.Logger.Error("Something wrong recovering post information\nDetail: ", err.Error())
		http.Error(w, "Something wrong retrieving post", http.StatusInternalServerError)
		return
	}
	err = PostAPI.FromDatabase(postDB)
	if err != nil {
		context.Logger.Error("Something wrong casting post structure\nDetail: ", err.Error())
		http.Error(w, "Something wrong retrieving post", http.StatusInternalServerError)
		return
	}

	result := map[string]Post{
		"post": PostAPI,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}
