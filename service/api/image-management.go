// Package api
/* This file consists in all the function used to manage binary image */
package api

import (
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// detectImageType allows to determine if the MIME type is a PNG or a JPEG.
// If the type is neither PNG nor JPEG it will return ""
func detectImageType(data []byte, context *reqcontext.RequestContext) string {
	dataType := http.DetectContentType(data)
	context.Logger.Info("Detected dataType: ", dataType)
	switch dataType {
	case "image/jpeg":
		return "jpeg"
	case "image/png":
		return "png"
	default:
		return ""
	}
}

// createDir allows to create a complete path
// If the folders in the path are already created, nothing happens
func createDirs(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// saveImage allows to save an image file in a specific directory
func saveImage(body io.Reader, directory, filename string) error {

	// Make complete path file
	filePath := filepath.Join(directory, filename)

	// Make or open file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Copy image content on file
	_, err = io.Copy(file, body)
	if err != nil {
		return err
	}

	return nil
}

// imageExists checks if an image exists. The function check both PNG and JPEG format.
// if the image exists, function will return the full namepath, otherwise return ""
func imageExists(fileName string, filePath string) (string, error) {

	checkImage := func(filePath string) (bool, error) {
		_, err := os.Stat(filePath)

		if err == nil {
			// File exists
			return true, nil
		}

		if os.IsNotExist(err) {
			// Image doesn't exist
			return false, nil
		}

		// There's something wrong, check permissions
		return false, err
	}

	pngName := fileName + ".png"
	jpegName := fileName + ".jpeg"

	// First check PNG
	fullNamePath := filePath + "/" + pngName
	if exists, err := checkImage(fullNamePath); err != nil {
		return "", err
	} else if exists {
		return fullNamePath, nil
	}

	// Check for JPEG
	fullNamePath = filePath + "/" + jpegName
	if exists, err := checkImage(fullNamePath); err != nil {
		return "", err
	} else if exists {
		return fullNamePath, nil
	}

	return "", nil
}

// deleteImage allows to remove a file owned by webserver.
// Function will return nil if file is correctly removed, an error otherwise
func deleteImage(pathFile string, nameFile string) error {
	fullPath := pathFile + nameFile
	err := os.Remove(fullPath)
	return err
}
