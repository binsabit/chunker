package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/binsabit/chunker/internal/data"
	"github.com/binsabit/chunker/pkg/helpers"
)

func InitConnection(username, filename string) (data.InitResponse, error) {
	user := data.InitRequest{
		Username: username,
		Filename: filename,
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Println("could not init connection", err)
		return data.InitResponse{}, err
	}
	b := bytes.NewReader(userJSON)
	req, err := http.NewRequest("POST", "http://localhost:5000/init", b)
	if err != nil {
		log.Println("could not init connection", err)
		return data.InitResponse{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return data.InitResponse{}, err
	}
	var upload data.InitResponse

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&upload)
	if err != nil {
		return data.InitResponse{}, err
	}
	upload.Filename = user.Filename
	return upload, nil
}

func Send(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	username := queryValues.Get("username")
	filename := queryValues.Get("filename")

	init, err := InitConnection(username, filename)
	if err != nil {
		log.Println("could not initiate conncection", err)
	}
	file, err := os.Open(path.Join("./res", username, filename))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	buffer := make([]byte, data.BufferSize)
	now := time.Now()
	chunkID := 0
	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		last := false
		if bytesread < data.BufferSize {
			last = true
		}
		// pr, pw := io.Pipe()
		// go func() {
		// 	defer pw.Close()
		// 	if _, err := io.WriteString(pw, string(buffer[:bytesread])); err != nil {
		// 		log.Println(err)
		// 	}
		// }()

		// fmt.Println("bytes read: ", bytesread)
		// req, err := http.NewRequest("POST", "http://localhost:5000/recieve", pr)
		// if err != nil {
		// 	log.Println("could not create request: %w", err)
		// 	return
		// }
		chunk := data.Chunk{
			Username: init.Username,
			UploadID: init.UploadID,
			Size:     bytesread,
			ID:       chunkID,
			Hash:     helpers.ToBase64(helpers.HashData(buffer[:bytesread])),
			Content:  helpers.ToBase64(buffer[:bytesread]),
			Last:     last,
			Filename: init.Filename,
		}

		chunkJSON, err := json.Marshal(chunk)
		if err != nil {
			log.Println("could not struct -> json", err)
			return
		}

		b := bytes.NewReader(chunkJSON)
		log.Println(chunk.UploadID, chunk.Username, chunkID)
		req, err := http.NewRequest("POST", "http://localhost:5000/recieve", b)
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		if err != nil {
			log.Println("could not create request: %w", err)
			return
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("bad request done %w", err)
		}
		fmt.Println(res.Status, bytesread)
		chunkID++
	}

	fmt.Println(time.Since(now))
}

func Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1 << 20)

	f, fh, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	filename := fh.Filename

	filepath := path.Join("./res/", r.FormValue("username"))
	_ = os.Mkdir(filepath, os.ModePerm)
	file, err := os.OpenFile(path.Join(filepath, filename), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	io.Copy(file, f)
}

func ShowUplaodPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}
