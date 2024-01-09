package api

import (
	"encoding/json"
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// getImage allows recovering an image passing the image ID.
// If the image id doesn't exist, the request will fail.
// If the user is not authorized, the request will fail.
// Note: image id is the same of post id.
func (rt *_router) getImage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// The image ID in the path is a 64-bit unsigned integer. Let's parse it.
	imageid, err := strconv.ParseUint(params.ByName("imageid"), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing imageid in request.")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]string{
			"error": "not correct format for imageid",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the Bearer Authorization Token is set
	if !rt.isAuthorized(r.Header) {
		context.Logger.Error("The bearer format token is not valid for getting an image!")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		response := map[string]string{
			"error": "log to do action",
		}

		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// check if the image exist
	fileName := strconv.FormatUint(imageid, 10)
	fileName, err = imageExists(fileName, "media/img")

	if err != nil {
		context.Logger.Error("Image doesn't exists\nDetail: ", err.Error())
		http.Error(w, "Image seems not exist", http.StatusNotFound)
		return
	}

	// File doesn't exist
	if fileName == "" {
		context.Logger.Error("Requested image doesn't exist")
		http.Error(w, "Images seems not exist", http.StatusNotFound)
		return
	}

	// Read the image content and prepare it for sending
	imageFile, err := os.Open(fileName)
	if err != nil {
		context.Logger.Error("Error during image opening")
		http.Error(w, "Something wrong retrieving image", http.StatusInternalServerError)
		return
	}
	defer func(imageFile *os.File) {
		_ = imageFile.Close()
	}(imageFile)

	fileInfo, err := imageFile.Stat()

	// Set Content-Type Header to image/png or image/jpeg
	if check := strings.HasSuffix(fileName, ".png"); check {
		w.Header().Set("Content-Type", "image/png")
	} else {
		w.Header().Set("Content-Type", "image/jpeg")
	}

	// Now return the binary image
	// NOTE: w.WriteHeader(http.StatusOK) is unnecessary because http.ServeContent already set it
	http.ServeContent(w, r, fileName, fileInfo.ModTime(), imageFile)
}
