package image

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cpuguy83/go-docker/httputil"
	"github.com/cpuguy83/go-docker/version"
)

// PullConfig holds the options for creating images.
type PullConfig struct {
	Image    string
	Repo     string
	Tag      string
	Platform string
}

type PullResponseMessage struct {
	Status string `json:"status"`
}

// PullOption is used as functional arguments to create images from a registry pull.
// PullOption configure a PullConfig.
type PullOption func(config *PullConfig)

func (s *Service) PullAsync(ctx context.Context, opts ...PullOption) (chan PullResponseMessage, error) {
	c := make(chan PullResponseMessage)
	cfg := PullConfig{}
	for _, o := range opts {
		o(&cfg)
	}

	withPullConfig := func(req *http.Request) error {
		q := req.URL.Query()
		q.Add("fromImage", cfg.Image)
		q.Add("fromSrc", "")
		q.Add("repo", cfg.Repo)
		q.Add("tag", cfg.Tag)
		q.Add("platform", cfg.Platform)
		req.URL.RawQuery = q.Encode()
		return nil
	}

	resp, err := httputil.DoRequest(ctx, func(ctx context.Context) (*http.Response, error) {
		return s.tr.Do(ctx, http.MethodPost, version.Join(ctx, "/images/create"), withPullConfig)
	})
	if err != nil {
		return c, err
	}

	// dockerd reports a stream of messages like
	// {"status":"Pulling from library/busybox","id":"latest"}
	// {"status":"Digest: sha256:f7ca5a32c10d51aeda3b4d01c61c6061f497893d7f6628b92f822f7117182a57"}
	// {"status":"Status: Image is up to date for busybox:latest"}
	//
	// Reading all is required otherwise the docker daemon considers the context
	// as cancelled and the pulling is aborted.
	go func() {
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			msg := PullResponseMessage{}
			s := scanner.Bytes()
			if err := json.Unmarshal(s, &msg); err != nil {
				continue
			}
			c <- msg
		}
	}()
	return c, nil
}

func (s *Service) Pull(ctx context.Context, opts ...PullOption) (*Image, error) {
	cfg := PullConfig{}
	for _, o := range opts {
		o(&cfg)
	}

	withPullConfig := func(req *http.Request) error {
		q := req.URL.Query()
		q.Add("fromImage", cfg.Image)
		q.Add("fromSrc", "")
		q.Add("repo", cfg.Repo)
		q.Add("tag", cfg.Tag)
		q.Add("platform", cfg.Platform)
		req.URL.RawQuery = q.Encode()
		return nil
	}

	resp, err := httputil.DoRequest(ctx, func(ctx context.Context) (*http.Response, error) {
		return s.tr.Do(ctx, http.MethodPost, version.Join(ctx, "/images/create"), withPullConfig)
	})
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Reading all is required otherwise the docker daemon considers the context
	// as cancelled and the pulling is aborted.
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	images, err := s.List(ctx, WithListFilterReference([]string{cfg.Image + ":" + cfg.Tag}))
	if err != nil {
		return nil, err
	}
	if len(images) == 0 {
		return nil, fmt.Errorf("image %s not found after successful pull", cfg.Image+":"+cfg.Tag)
	}
	return &Image{id: images[0].ID, tr: s.tr}, nil
}
