package image

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cpuguy83/go-docker/httputil"

	"github.com/cpuguy83/go-docker/version"

	"github.com/cpuguy83/go-docker/image/imageapi"
	"github.com/cpuguy83/go-docker/transport"
	"github.com/pkg/errors"
)

// DefaultInspectDecodeLimitBytes is the default value used for limit how much data is read from the inspect response.
const DefaultInspectDecodeLimitBytes = 64 * 1024

// InspectConfig holds the options for inspecting a image
type InspectConfig struct {
	// Allows callers of `Inspect` to unmarshal to any object rather than only the built-in types.
	// This is useful for anyone wrapping the API and providing more metadata (e.g. classic swarm)
	// To must be a pointer or it may cause a panic.
	// If `To` is provided, `Inspect`'s returned image object may be empty.
	To interface{}
}

// InspectOption is used as functional arguments to inspect an image
// InspectOptions configure an InspectConfig.
type InspectOption func(config *InspectConfig)

// Inspect fetches detailed information about an image
func (s *Service) Inspect(ctx context.Context, name string, opts ...InspectOption) (imageapi.ImageInspect, error) {
	return handleInspect(ctx, s.tr, name, opts...)
}

func handleInspect(ctx context.Context, tr transport.Doer, name string, opts ...InspectOption) (imageapi.ImageInspect, error) {
	cfg := InspectConfig{}
	for _, o := range opts {
		o(&cfg)
	}

	var c imageapi.ImageInspect

	resp, err := httputil.DoRequest(ctx, func(ctx context.Context) (*http.Response, error) {
		return tr.Do(ctx, http.MethodGet, version.Join(ctx, "/images/"+name+"/json"))
	})
	if err != nil {
		return c, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c, nil
	}

	if cfg.To != nil {
		if err := json.Unmarshal(data, cfg.To); err != nil {
			return c, errors.Wrap(err, "error unmarshalling to requested type")
		}
		return c, nil
	}

	if err := json.Unmarshal(data, &c); err != nil {
		return c, errors.Wrap(err, "error unmarshalling image json")
	}

	return c, nil
}

// Inspect fetches detailed information about the image.
func (c *Image) Inspect(ctx context.Context, opts ...InspectOption) (imageapi.ImageInspect, error) {
	return handleInspect(ctx, c.tr, c.id, opts...)
}
