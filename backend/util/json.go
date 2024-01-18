package util

import (
	"encoding/json"
	"net/http"
)

func ReadJSONReq[T any](r *http.Request, dest *T) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	return decoder.Decode(dest)
}
