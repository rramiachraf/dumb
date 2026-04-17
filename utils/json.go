package utils

import (
	"encoding/json"
	"net/http"
)

func EncodeJSON(w http.ResponseWriter, data any, l *Logger) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		l.Errorf("unable to render json %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte{})
		if err != nil {
			l.Error(err.Error())
		}
	}
}
