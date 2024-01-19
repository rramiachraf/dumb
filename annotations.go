package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type annotationsResponse struct {
	Response struct {
		Referent struct {
			Annotations []struct {
				Body struct {
					Html string `json:"html"`
				} `json:"body"`
			} `json:"annotations"`
		} `json:"referent"`
	} `json:"response"`
}

func annotationsHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if data, err := getCache(id); err == nil {

		response, err := json.Marshal(data)

		if err != nil {
			logger.Errorf("could not marshal json: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			render("error", w, map[string]string{
				"Status": "500",
				"Error":  "Could not parse genius api response",
			})
			return
		}
		w.Header().Set("content-type", "application/json")
		_, err = w.Write(response)
		if err != nil {
			logger.Errorln("Error sending response: ", err)
		}
		return
	}

	url := fmt.Sprintf("https://genius.com/api/referents/%s?text_format=html", id)
	resp, err := sendRequest(url)

	if err != nil {
		logger.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		render("error", w, map[string]string{
			"Status": "500",
			"Error":  "cannot reach genius servers",
		})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		render("error", w, map[string]string{
			"Status": "404",
			"Error":  "page not found",
		})
		return
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		logger.Errorln("Error paring genius api response", err)
		w.WriteHeader(http.StatusInternalServerError)
		render("error", w, map[string]string{
			"Status": "500",
			"Error":  "Parsing error",
		})
		return
	}

	var data annotationsResponse
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		logger.Errorf("could not unmarshal json: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		render("error", w, map[string]string{
			"Status": "500",
			"Error":  "Could not parse genius api response",
		})
		return
	}

	w.Header().Set("content-type", "application/json")
	body := data.Response.Referent.Annotations[0].Body
	body.Html = cleanBody(body.Html)
	response, err := json.Marshal(body)

	if err != nil {
		logger.Errorf("could not marshal json: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		render("error", w, map[string]string{
			"Status": "500",
			"Error":  "Could not parse genius api response",
		})
		return
	}

	setCache(id, body)
	_, err = w.Write(response)
	if err != nil {
		logger.Errorln("Error sending response: ", err)
	}
}

func cleanBody(body string) string {
	var withCleanedImageLinks = strings.Replace(body, "https://images.rapgenius.com/", "/images/", -1)

	var re = regexp.MustCompile(`https?:\/\/[a-z]*.?genius.com`)
	var withCleanedLinks = re.ReplaceAllString(withCleanedImageLinks, "")

	return withCleanedLinks
}
