package reciever

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/binsabit/chunker/internal/data"
	"github.com/binsabit/chunker/pkg/helpers"
)

func Recieve(w http.ResponseWriter, r *http.Request) {
	var chunk data.Chunk

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&chunk)
	if err != nil {
		log.Println("could not read response", err)
		return
	}

	if isDamaged(chunk.Hash, helpers.ToBase64(helpers.HashData(helpers.FromBase64(chunk.Content)))) {
		log.Println("file dameged")
		return
	}
	filepath := path.Join(chunk.Username, strconv.Itoa(chunk.UploadID), chunk.Filename)
	log.Println(filepath)
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println("could not write chunk", err)
		return
	}
	_, _ = file.WriteString(string(helpers.FromBase64(chunk.Content)))
	defer file.Close()
}

func InitUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("req")
	var input struct {
		Username string `json:"username"`
		Filename string `json:"filename"`
	}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(input.Filename)
	UploadID := make([]byte, 16)
	rand.Read(UploadID)
	fmt.Println(UploadID)
	response := data.InitResponse{
		OK:       true,
		UploadID: int(binary.LittleEndian.Uint32(UploadID)),
		Username: input.Username,
	}
	_ = os.Mkdir(input.Username, os.ModePerm)
	_ = os.Mkdir(path.Join(input.Username, strconv.Itoa(response.UploadID)), os.ModePerm)
	_, _ = os.Create(path.Join(input.Username, strconv.Itoa(response.UploadID), input.Filename))
	// os.Create(path.Join(input.Username, strconv.Itoa(response.UploadID)))
	js, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error in writing json", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(js)

}

func isDamaged(hash, content string) bool {
	return hash != content
}
