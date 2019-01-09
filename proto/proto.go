package proto

type Filecache struct {
	Names []string `json:"names, omitempty"`
}

type PostRequest struct {
	Env string `json:"env"`
	Filecache Filecache `json:"filecache, omitempty"`
}