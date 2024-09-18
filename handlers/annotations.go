package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rramiachraf/dumb/data"
	"github.com/rramiachraf/dumb/utils"
	"github.com/rramiachraf/dumb/views"
)

func annotations(l *utils.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["annotation-id"]
		if a, err := getCache[data.Annotation]("annotation:" + id); err == nil {
			encoder := json.NewEncoder(w)

			w.Header().Set("content-type", "application/json")
			if err = encoder.Encode(&a); err != nil {
				l.Error(err.Error())
			}

			return
		}

		url := fmt.Sprintf("https://genius.com/api/referents/%s?text_format=html", id)
		resp, err := utils.SendRequest(url)

		if err != nil {
			l.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderTemplate(w, views.ErrorPage(500, "cannot reach genius servers"), l)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			w.WriteHeader(http.StatusNotFound)
			utils.RenderTemplate(w, views.ErrorPage(404, "page not found"), l)
			return
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			l.Error("Error paring genius api response: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderTemplate(w, views.ErrorPage(500, "something went wrong"), l)
			return
		}

		var data data.AnnotationsResponse
		err = json.Unmarshal(buf.Bytes(), &data)
		if err != nil {
			l.Error("could not unmarshal json: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.RenderTemplate(w, views.ErrorPage(500, "something went wrong"), l)
			return
		}

		annotation := data.Response.Referent.Annotations[0]
		annotation.Body.HTML = utils.CleanBody(annotation.Body.HTML)

		w.Header().Set("content-type", "application/json")
		encoder := json.NewEncoder(w)

		if err = encoder.Encode(&annotation); err != nil {
			l.Error("Error sending response: %s", err.Error())
			return
		}

		if err = setCache("annotation:"+id, annotation); err != nil {
			l.Error(err.Error())
		}
	}
}
