package tgit

import (
	"fmt"
	"net/http"
)

type RepositoriesService struct {
	client *Client
}

type Diff struct {
	OldPath     string `json:"old_path,omitempty"`
	NewPath     string `json:"new_file,omitempty"`
	AMode       int    `json:"a_mode,omitempty"`
	BMode       int    `json:"b_mode,omitempty"`
	Diff        string `json:"diff,omitempty"`
	NewFile     bool   `json:"new_file,omitempty"`
	RenamedFile bool   `json:"renamed_file,omitempty"`
	DeletedFile bool   `json:"deleted_file,omitempty"`
	IsTooLarge  bool   `json:"is_too_large,omitempty"`
	IsCollapse  bool   `json:"is_collapse,omitempty"`
	Additions   int    `json:"additions,omitempty"`
	Deletions   int    `json:"deletions,omitempty"`
}

type Compare struct {
	Commit         *Commit   `json:"commit,omitempty"`
	Commits        []*Commit `json:"commits,omitempty"`
	Diffs          []*Diff   `json:"diffs,omitempty"`
	CompareTimeout bool      `json:"compare_timeout,omitempty"`
	CompareSameRef bool      `json:"compare_same_ref,omitempty"`
	OverFlow       bool      `json:"over_flow,omitempty"`
	FilesTotal     int       `json:"files_total,omitempty"`
	CommitsTotal   int       `json:"commits_total,omitempty"`
}

type CompareOptions struct {
	From string `url:"from,omitempty"`
	To   string `url:"to,omitempty"`
}

func (s *RepositoriesService) Compare(pid interface{}, opts *CompareOptions) (*Compare, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/compare", pathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, opts)
	if err != nil {
		return nil, nil, err
	}

	c := new(Compare)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, err
}
