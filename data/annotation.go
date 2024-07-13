package data

type AnnotationsResponse struct {
	Response struct {
		Referent struct {
			Annotations []struct {
				Body Annotation `json:"body"`
			} `json:"annotations"`
		} `json:"referent"`
	} `json:"response"`
}

type Annotation struct {
	HTML string `json:"html"`
}
