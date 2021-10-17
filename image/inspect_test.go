package image_test

import (
	"context"
	"testing"

	"github.com/cpuguy83/go-docker/image"
	"gotest.tools/assert"
)

const imageID = "sha256:16ea53ea7c652456803632d67517b78a4f9075a10bfdc4fc6b7b4cbf2bc98497"

func TestInspect(t *testing.T) {
	ctx := context.Background()
	s := newTestService(t)

	// First create a few images that we can list later.
	img, err := s.Pull(ctx, func(config *image.PullConfig) {
		config.Image = "busybox"
		config.Tag = "latest"
		config.Repo = "dockerhub.io"
	})
	assert.NilError(t, err, "expected pulling busybox to succeed")

	data, err := img.Inspect(ctx)
	assert.NilError(t, err, "expected inspect image on ID to succeed")
	assert.Assert(t, data.ID == imageID, "expected inspect data ID field to be"+imageID)
}
