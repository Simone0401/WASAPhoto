package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// uploadPost allows to add a photo to the collection of posts.
// If the user in not authorized, the request will fail.
// If the MIME type is not PNG or JPEG the request will fail.
// The function will return the post ID created for new image
func (rt *_router) uploadPost(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The User ID in the path is a 64-bit unsigned integer. Let's parse it.
	uid, err := strconv.ParseUint(params.ByName("uid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing uid in uploading post request")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "not correct format",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the Bearer Authorization Token is set
	if !rt.isAuthorized(r.Header) {
		context.Logger.Error("The bearer format token is not valid in uploading post request!")
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

	// read the body
	body, err := io.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	if err != nil {
		context.Logger.Error("Unable to read binary image for uploading post\nDetail: ", err.Error())
		http.Error(w, "Something wrong uploading the image", http.StatusInternalServerError)
		return
	}

	// check if the Content-Type in the request is correct and if the binary format is correct
	var imageType string
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/png") {
		imageType = detectImageType(body, &context)
		if imageType == "png" {
			context.Logger.Info("Image is correct PNG")
		} else {
			context.Logger.Error("Content-Type is a PNG but the body is not")
			http.Error(w, "File is not supported", http.StatusBadRequest)
			return
		}
	} else if strings.HasPrefix(contentType, "image/jpeg") {
		imageType = detectImageType(body, &context)
		if imageType == "jpeg" {
			context.Logger.Info("Image is a correct JPEG")
		} else {
			context.Logger.Error("Content-Type is a JPEG but the body is not")
			http.Error(w, "File is not supported", http.StatusBadRequest)
			return
		}
	} else {
		context.Logger.Error("File format is not valid")
		http.Error(w, "File is not supported", http.StatusBadRequest)
		return
	}

	// check if media/img folders already exists.
	// Create them if they are not present
	// media/img path folders store all the uploaded photos
	err = createDirs("media/img")
	if err != nil {
		context.Logger.Error("Error creating media/img folders\nDetail: ", err.Error())
		http.Error(w, "Error storing your image", http.StatusInternalServerError)
		return
	}

	// create post and get the post id
	imageId, err := rt.db.AddPost(uid)

	if err != nil {
		message := fmt.Sprintf("Error creating new post for user %d\nDetail: ", uid)
		context.Logger.Error(message, err.Error())
		http.Error(w, "Somenthing wrong adding the post", http.StatusInternalServerError)
		return
	}

	filename := strconv.FormatUint(imageId, 10)
	filename += "." + imageType

	// because r.Body was already read before, we need to put to pointer to the start
	r.Body = io.NopCloser(bytes.NewReader(body))
	err = saveImage(r.Body, "media/img", filename)

	if err != nil {
		context.Logger.Error("Error saving image on the disk\nDetail: ", err.Error())
		http.Error(w, "Somenthing wrong uploading your post", http.StatusInternalServerError)
		return
	}

	// image correctly saved on disk
	// now return the postid to the client
	w.Header().Set("Content-Type", "application/json")

	result := map[string]uint64{
		"postid": imageId,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
}
