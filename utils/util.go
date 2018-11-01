package utils

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
)

func Message(status uint, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func RespondWithMessage(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func RespondWithError(w http.ResponseWriter, errorMessage string, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": errorMessage})
}

func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func RespondWithFile(w http.ResponseWriter, file *os.File, fileName string) {
	io.Copy(w, file)
	contentType, _ := GetFileContentType(file)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", contentType)
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func CreateStorageDirectoryIfNotExists() {
	storagePath := os.Getenv("storage_path")
	if exists, err := FileExists(storagePath); !exists {
		err = os.MkdirAll(storagePath, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func GetFileContentType(out *os.File) (string, error) {

	buffer := make([]byte, 512)
	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func GetPathVar(name string, r *http.Request) (value string) {
	return mux.Vars(r)[name]
}
