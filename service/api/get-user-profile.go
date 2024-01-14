package api

import (
	"encoding/json"
	"fmt"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// getUserProfile allows getting user's profile information passing the uid.
// If the user in not authorized, the request will fail.
// If the user id doesn't exist, the request will fail.
// The return values will be all the user information and his upload post stream in reverse chronological order
// The stream consists in an array of post. (Check API documentation for detail)
// Note: for getting a binary image it's necessary using the 'Get Image API'
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing uid in get profile request")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "not correct format for uid",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the Bearer Authorization Token is set
	if !rt.isAuthorized(r.Header) {
		context.Logger.Error("The bearer format token is not valid in getting profile request!")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		response := map[string]string{
			"error": "log to do action",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the user exists
	check, err := rt.db.CheckExistsByUID(uid)
	if err != nil {
		context.Logger.Error("Error retrieving information on UID for getting profile request!")
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in getting profile request! User doesn't exist!")
		http.Error(w, "User seems not exist.", http.StatusNotFound)
		return
	}

	// check if the current user is authorized
	check, err = rt.db.CheckExistsByUID(context.Uid)
	if err != nil {
		context.Logger.Error("Error retrieving information on UID that makes getting profile request!")
		http.Error(w, "Something wrong in the server", http.StatusInternalServerError)
		return
	}

	if !check {
		context.Logger.Error("Error in getting profile request! User that makes request doesn't exist!")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Prepare return statement
	var profileInfo ProfileInfo
	var uploadedPost []Post

	profileInfo.User.Userid = uid

	userDb := profileInfo.User.ToDatabase()
	userDb, err = rt.db.GetUserByID(userDb.Userid)
	if err != nil {
		context.Logger.Error("Error getting User from uid for getting profile request!\nDetail: ", err.Error())
		http.Error(w, "Something wrong retrieving user profile", http.StatusInternalServerError)
		return
	}

	// Get Profile information
	profileDB, err := rt.db.GetProfileInfo(userDb.Userid)
	if err != nil {
		context.Logger.Error("Error getting ProfileDB for getting profile request!\nDetail: ", err.Error())
		http.Error(w, "Something wrong retrieving user profile", http.StatusInternalServerError)
		return
	}

	_ = profileInfo.FromDatabase(profileDB)

	// Get profile posts in reverse chronological order
	listPost, err := rt.db.GetProfilePosts(userDb.Userid)
	if err != nil {
		context.Logger.Error("Error retrieving user profile posts during getting profile request!\nDetail: ", err.Error())
		http.Error(w, "Something wrong retrieving profile", http.StatusInternalServerError)
		return
	}

	// Append each post to the list
	for i, post := range listPost {
		var postAPI Post
		err = postAPI.FromDatabase(post)
		if err != nil {
			mess := fmt.Sprintf("Error parsing postDB to postAPI for post number %d in getting profile request\nDetail: ", i)
			context.Logger.Error(mess, err.Error())
			http.Error(w, "Something wrong retrieving your profile", http.StatusInternalServerError)
			return
		}
		// Change datetime format for each comment
		for i := 0; i < len(postAPI.Comments); i++ {
			postAPI.Comments[i].Datetime, _ = formatDatetime(postAPI.Comments[i].Datetime)
		}
		postAPI.Datetime, _ = formatDatetime(postAPI.Datetime)
		uploadedPost = append(uploadedPost, postAPI)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&map[string]interface {
	}{"profile_info": profileInfo,
		"uploaded_post": uploadedPost})

}
