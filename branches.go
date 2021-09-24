package tgit

import (
	"fmt"
	"net/http"
	"net/url"
)

type BranchesService struct {
	client *Client
}

type Branch struct {
	Commit             *Commit `json:"commit"`
	Name               string  `json:"name"`
	Protected          bool    `json:"protected"`
	DevelopersCanPush  bool    `json:"developers_can_push"`
	DevelopersCanMerge bool    `json:"developers_can_merge"`
}

func (b Branch) String() string {
	return Stringify(b)
}

type ListBranchesOptions struct {
	ListOptions
}

func (s *BranchesService) ListBranches(pid interface{}, opts *ListBranchesOptions) ([]*Branch, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/branches", pathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, opts)
	if err != nil {
		return nil, nil, err
	}

	var b []*Branch
	resp, err := s.client.Do(req, &b)
	if err != nil {
		return nil, resp, err
	}

	return b, resp, err
}

func (s *BranchesService) GetBranch(pid interface{}, branch string) (*Branch, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/branches/%s", pathEscape(project), url.PathEscape(branch))

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	b := new(Branch)
	resp, err := s.client.Do(req, b)
	if err != nil {
		return nil, resp, err
	}

	return b, resp, err
}
