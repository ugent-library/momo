package render

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ugent-library/momo/web/ctx"
)

func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ctx.GetStatus(r))
	w.Write(buf.Bytes())
}
