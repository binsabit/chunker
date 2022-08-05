package data

const BufferSize = 10485760

type Chunk struct {
	UploadID int    `json:"connection_id"`
	Username string `json:"username"`
	Size     int    `json:"size"`
	ID       int    `json:"chunk_id"`
	Hash     string `json:"hash"`
	Content  string `json:"content"`
	Last     bool   `json:"last"`
	Filename string `json:"filename"`
}

type InitResponse struct {
	OK       bool   `json:"ok"`
	UploadID int    `json:"upload_id"`
	Username string `json:"username"`
	Filename string `json:"filename"`
}

type InitRequest struct {
	Username string `json:"username"`
	Filename string `json:"filename"`
}

type User struct {
	ID       int    `json:"user_id"`
	Username string `json:"username"`
}
