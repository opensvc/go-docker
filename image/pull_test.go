package image_test

import (
	"context"
	"testing"

	"github.com/cpuguy83/go-docker/image"
	"gotest.tools/assert"
)

func TestPull(t *testing.T) {
	ctx := context.Background()
	s := newTestService(t)

	_, err := s.Pull(ctx)
	assert.Assert(t, err != nil, "expected pull with no options to fail")

	_, err = s.Pull(ctx, func(config *image.PullConfig) {
		config.Image = "busybox"
		config.Tag = "latest"
		config.Repo = "dockerhub.io"
	})
	assert.NilError(t, err, "expected pulling busybox to succeed")
}

func TestPullAsync(t *testing.T) {
	ctx := context.Background()
	s := newTestService(t)

	c, err := s.PullAsync(ctx, func(config *image.PullConfig) {
		config.Image = "busybox"
		config.Tag = "latest"
		config.Repo = "dockerhub.io"
	})
	assert.NilError(t, err, "expected pulling busybox to succeed")
	msg := <-c
	assert.Assert(t, msg.Status != "")
}
