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
		Html string `json:"html"`
	} `json:"body"`
}
