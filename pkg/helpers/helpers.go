package helpers

import (
	"crypto/sha1"
	"encoding/base64"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Envelope map[string]interface{}

func ReadIDParam(r *http.Request) string {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")
	return id
}

func HashData(data []byte) []byte {
	hash := sha1.Sum(data)
	return hash[:]
}

func MakeReuest(method, to string, data []byte) error {

	return nil
}

func ToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func FromBase64(data string) []byte {
	res, _ := base64.StdEncoding.DecodeString(data)
	return res
}
