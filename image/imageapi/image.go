package imageapi

import "github.com/cpuguy83/go-docker/container/containerapi"

// Image represents an image from the docker HTTP API.
type Image struct {
	ID          string `json:"Id,omitempty"`
	ParentID    string `json:"ParentId,omitempty"`
	RepoTags    []string
	RepoDigests []string
	Created     int64
	Size        int64
	SharedSize  int64
	VirtualSize int64
	Labels      map[string]string
	Containers  int64
}

// ImageInspect is newly used struct along with MountPoint
type ImageInspect struct {
	ID            string              `json:"Id,omitempty"`
	RepoTags      []string            `json:"RepoTags,omitempty"`
	RepoDigests   []string            `json:"RepoDigests,omitempty"`
	Parent        string              `json:"Parent,omitempty"`
	Comment       string              `json:"Comment,omitempty"`
	Created       string              `json:"Created,omitempty"`
	Container     string              `json:"Container,omitempty"`
	Config        containerapi.Config `json:"Config"`
	DockerVersion string              `json:"DockerVersion,omitempty"`
	Author        string              `json:"Author,omitempty"`
	Architecture  string              `json:"Architecture,omitempty"`
	OS            string              `json:"OS,omitempty"`
	Size          int64               `json:"Size,omitempty"`
	VirtualSize   int64               `json:"VirtualSize,omitempty"`
	GraphDriver   GraphDriver         `json:"GraphDriver,omitempty"`
	RootFS        RootFS              `json:"RootFS,omitempty"`
	Metadata      Metadata            `json:"Metadata,omitempty"`
}

type GraphDriver struct {
	Data *string `json:"Data,omitempty"`
	Name string  `json:"Name,omitempty"`
}

type RootFS struct {
	Type   string   `json:"Type,omitempty"`
	Layers []string `json:"Layers,omitempty"`
}

type Metadata struct {
	LastTagTime string `json:"LastTagTime,omitempty"`
}
