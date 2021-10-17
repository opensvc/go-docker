package image

import (
	"context"
	"net/http"

	"github.com/cpuguy83/go-docker/httputil"
	"github.com/cpuguy83/go-docker/version"
)

// ImportConfig holds the options for creating images.
type ImportConfig struct {
	Source   string
	Repo     string
	Tag      string
	Platform string
}

// ImportOption is used as functional arguments to create images from a local source file.
// ImportOption configure a ImportConfig.
type ImportOption func(config *ImportConfig)

func (s *Service) Import(ctx context.Context, opts ...ImportOption) error {
	cfg := ImportConfig{}
	for _, o := range opts {
		o(&cfg)
	}

	withImportConfig := func(req *http.Request) error {
		q := req.URL.Query()
		q.Add("fromImage", "")
		q.Add("fromSrc", cfg.Source)
		q.Add("repo", cfg.Repo)
		q.Add("tag", cfg.Tag)
		q.Add("platform", cfg.Platform)
		req.URL.RawQuery = q.Encode()
		return nil
	}

	resp, err := httputil.DoRequest(ctx, func(ctx context.Context) (*http.Response, error) {
		return s.tr.Do(ctx, http.MethodPost, version.Join(ctx, "/images/create"), withImportConfig)
	})
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
