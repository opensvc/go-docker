package image_test

import (
	"context"
	"testing"

	"github.com/cpuguy83/go-docker/image"
	"gotest.tools/assert"
)

func TestList(t *testing.T) {
	ctx := context.Background()
	s := newTestService(t)

	// First create a few images that we can list later.
	_, err := s.Pull(ctx, func(config *image.PullConfig) {
		config.Image = "busybox"
		config.Tag = "latest"
		config.Repo = "dockerhub.io"
	})
	assert.NilError(t, err, "expected pulling busybox to succeed")

	images, err := s.List(ctx, func(config *image.ListConfig) {
		config.Filter.Reference = append(config.Filter.Reference, "busybox:latest")
	})
	assert.NilError(t, err, "expected listing images with no options to succeed")
	assert.Assert(t, len(images) == 1, "expected created images to be listed")

	img, err := s.FindImage(ctx, "busybox:latest")
	assert.NilError(t, err, "expected FindImage to succeed")
	assert.Assert(t, img != nil, "expected image to be found")

}
