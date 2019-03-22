package seed

type Source struct {
	Bangumi string        `json:"bangumi"`
	Path    string        `json:"path"`
	Poster  string        `json:"poster"`
	Role    []interface{} `json:"role"`
	Publish string        `json:"publish"`
}

func ReadJSON(path string) {

}
