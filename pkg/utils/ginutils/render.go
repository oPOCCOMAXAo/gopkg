package ginutils

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	pkgerr "github.com/pkg/errors"
)

func Render(ctx *gin.Context, status int, component templ.Component) {
	ctx.Render(status, &templRenderer{
		Context:   ctx.Request.Context(),
		Component: component,
	})
}

type templRenderer struct {
	Context   context.Context //nolint:containedctx
	Component templ.Component
}

func (t templRenderer) Render(w http.ResponseWriter) error {
	t.WriteContentType(w)

	if t.Component != nil {
		err := t.Component.Render(t.Context, w)
		if err != nil {
			return pkgerr.WithStack(err)
		}
	}

	return nil
}

func (t templRenderer) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
