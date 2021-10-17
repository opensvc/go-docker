package image

import (
	"context"
	"fmt"

	"github.com/cpuguy83/go-docker/transport"
)

// Image provides bindings for interacting with a image in Docker
type Image struct {
	id string
	tr transport.Doer
}

// ID returns the image ID
func (c *Image) ID() string {
	return c.id
}

// NewConfig holds the options available for `New`
type NewConfig struct {
}

// NewOption is used as functional parameters to `New`
type NewOption func(*NewConfig)

// New creates a new image object in memory. This function does not interact with the Docker API at all.
// If the image does not exist in Docker, all calls on the Image will fail.
//
// To actually create a image you must call `Pull` or `Import` first (which will return a image object to you).
//
func (s *Service) NewImage(_ context.Context, id string, opts ...NewOption) *Image {
	var cfg NewConfig
	for _, o := range opts {
		o(&cfg)
	}
	return &Image{id: id, tr: s.tr}
}

func (s *Service) FindImage(ctx context.Context, name string) (*Image, error) {
	l, err := s.List(ctx, WithListFilterReference([]string{name}))
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, fmt.Errorf("image not found")
	}
	return &Image{id: l[0].ID, tr: s.tr}, nil
}
