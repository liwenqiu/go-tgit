package tgit

import (
	"fmt"
	"net/http"
	"net/url"
)

type TagsService struct {
	client *Client
}

type Tag struct {
	Commit  *Commit `json:"commit"`
	Name    string  `json:"name"`
	Message string  `json:"message"`
}

func (t Tag) String() string {
	return Stringify(t)
}

type ListTagsOptions struct {
	ListOptions
}

func (s *TagsService) ListTags(pid interface{}, opts *ListTagsOptions) ([]*Tag, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/tags", pathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, opts)
	if err != nil {
		return nil, nil, err
	}

	var t []*Tag
	resp, err := s.client.Do(req, &t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, err
}

func (s *TagsService) GetTag(pid interface{}, tag string) (*Tag, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/tags/%s", pathEscape(project), url.PathEscape(tag))

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var t *Tag
	resp, err := s.client.Do(req, &t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, err
}
