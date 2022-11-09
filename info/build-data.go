package info

type BuildDataPage struct {
	Pathname string `json:"pathname"`
}

type BuildData struct {
	Pages []BuildDataPage `json:"pages"`
}
