package api

import (
	"encoding/json"
	"net/http"

	"github.com/golangpub/gox/v2"
	"github.com/golangpub/types"
	"github.com/golangpub/wine"
	"github.com/golangpub/wine/mime"
)

// ParseParams unmarshal request body into params. Only support JSON type for now.
func ParseParams(req *wine.Request, params interface{}) error {
	// Unsafe assignment, so ignore error
	data, err := json.Marshal(req.Params())
	if err == nil {
		_ = json.Unmarshal(data, params)
	}
	_ = gox.CopyWithNamer(params, req.Params(), gox.SnakeToCamelNamer)

	if req.ContentType() == mime.JSON {
		err := json.Unmarshal(req.Body(), params)
		if err != nil {
			return types.NewError(http.StatusBadRequest, "unmarshal: %v", err)
		}
	}
	return nil
}
