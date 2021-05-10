package render

import (
	"bytes"
	"encoding/xml"
	"net/http"

	"github.com/ugent-library/momo/internal/ctx"
)

func XML(w http.ResponseWriter, r *http.Request, v interface{}) {
	buf := &bytes.Buffer{}
	enc := xml.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(ctx.GetStatus(r))
	w.Write(buf.Bytes())
}
