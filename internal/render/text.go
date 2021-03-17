package render

import (
	"net/http"

	"github.com/ugent-library/momo/internal/ctx"
)

func Text(w http.ResponseWriter, r *http.Request, v string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(ctx.GetStatus(r))
	w.Write([]byte(v))
}
