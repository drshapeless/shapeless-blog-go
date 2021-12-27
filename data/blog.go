package data

type PostMetadata struct {
	// ID is implied in the file path.
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Created  string   `json:"created"`
	Updated  string   `json:"updated"`
	Category []string `json:"category"`
}

type Blog struct {
	Metadata PostMetadata `json:"metadata"`
	Content  string       `json:"content"`
}
