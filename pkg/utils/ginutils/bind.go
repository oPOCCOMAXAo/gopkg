package ginutils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	pkgerr "github.com/pkg/errors"
)

// formMemory is the max memory used to parse multipart forms.
const formMemory = 32 << 20

// BindOptions selects request sources for [BindAll].
type BindOptions struct {
	// Uri is path params.
	Uri bool

	// Form is x-www-urlencoded/multipart/querystring.
	// Multipart file fields are not mapped.
	Form bool

	// JSON body.
	JSON bool
}

// BindAll maps all selected request sources onto objectRef
// and validates the result once at the end.
//
// Returns false and attaches the error to ctx (gin.ErrorTypeBind) on failure.
//
// Notes:
//   - Validation runs even when no source is enabled,
//     so zero-value fields fail "required".
//   - Form consumes the request body, so do not combine it with JSON.
//     The exception is a JSON body, where Form only maps the URL query.
func BindAll(
	ctx *gin.Context,
	objectRef any,
	opts BindOptions,
) bool {
	mappers := []struct {
		enabled bool
		fn      func(*gin.Context, any) error
	}{
		{opts.Uri, MapUri},
		{opts.Form, MapForm},
		{opts.JSON, MapJSON},
	}

	for _, mapper := range mappers {
		if !mapper.enabled {
			continue
		}

		err := mapper.fn(ctx, objectRef)
		if err != nil {
			ctx.Error(err).SetType(gin.ErrorTypeBind)

			return false
		}
	}

	if binding.Validator != nil {
		err := binding.Validator.ValidateStruct(objectRef)
		if err != nil {
			ctx.Error(pkgerr.WithStack(err)).SetType(gin.ErrorTypeBind)

			return false
		}
	}

	return true
}

// MapUri maps path parameters onto fields with the "uri" tag,
// without validation.
func MapUri(
	ctx *gin.Context,
	objectRef any,
) error {
	params := make(map[string][]string, len(ctx.Params))
	for _, v := range ctx.Params {
		params[v.Key] = append(params[v.Key], v.Value)
	}

	err := binding.MapFormWithTag(objectRef, params, "uri")
	if err != nil {
		return pkgerr.WithStack(err)
	}

	return nil
}

// MapForm maps the request form (urlencoded or multipart body,
// plus URL query) onto fields with the "form" tag, without validation.
// Multipart file fields are not mapped.
func MapForm(
	ctx *gin.Context,
	objectRef any,
) error {
	err := ctx.Request.ParseForm()
	if err != nil {
		return pkgerr.WithStack(err)
	}

	err = ctx.Request.ParseMultipartForm(formMemory)
	if err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return pkgerr.WithStack(err)
	}

	err = binding.MapFormWithTag(objectRef, ctx.Request.Form, "form")
	if err != nil {
		return pkgerr.WithStack(err)
	}

	return nil
}

// MapJSON decodes the request body as JSON onto objectRef,
// without validation.
func MapJSON(
	ctx *gin.Context,
	objectRef any,
) error {
	err := json.NewDecoder(ctx.Request.Body).Decode(objectRef)
	if err != nil {
		return pkgerr.WithStack(err)
	}

	return nil
}
