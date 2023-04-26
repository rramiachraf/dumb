package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type annotationResponse struct {
	Response struct {
		Referent struct {
			Classification string
			Annotations    []annotation
		}
	}
}

type annotation struct {
	Body struct {
		HTML string
	}
}

type annotationResult struct {
	Classification string
	Annotations    []annotation
}

func annotationHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	u := fmt.Sprintf("https://genius.com/api/referents/%s?text_format=html", v["id"])

	res, err := sendRequest(u)
	if err != nil {
		//TODO handle err
	}

	defer res.Body.Close()

	var data annotationResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&data)
	if err != nil {
		//TODO handle err
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	result := annotationResult{data.Response.Referent.Classification, data.Response.Referent.Annotations}
	encoder := json.NewEncoder(w)
	encoder.Encode(&result)
}
