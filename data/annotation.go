package data

type AnnotationsResponse struct {
	Response struct {
		Referent struct {
			Annotations []Annotation `json:"annotations"`
		} `json:"referent"`
	} `json:"response"`
}

type Annotation struct {
	Body struct {
		HTML string `json:"html"`
	} `json:"body"`
	State    string `json:"state"`
	Verified bool   `json:"verified"`
}
