package utils

import "net/http"

// WriteImage writes data to the response
func WriteImage(w http.ResponseWriter, data []byte) {
	w.Header().Add("Content-Type", "image/png")
	w.Write(data)
}
