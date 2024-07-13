package utils

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func RenderTemplate(w http.ResponseWriter, t templ.Component, l *Logger) {
	if err := t.Render(context.Background(), w); err != nil {
		l.Error("unable to render template %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte{})
		if err != nil {
			l.Error(err.Error())
		}
	}
}
