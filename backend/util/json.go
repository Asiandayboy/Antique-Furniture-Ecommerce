package util

import (
	"encoding/json"
	"net/http"
)

/*
This function returns an error if a problem occurred
decoding the json-encoded string request body
*/
func ReadJSONReq[T any](r *http.Request, dest *T) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	return decoder.Decode(dest)
}
