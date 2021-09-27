package tgit

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type CommitsService struct {
	client *Client
}

type Commit struct {
	ID             string     `json:"id"`
	ShortID        string     `json:"short_id"`
	Title          string     `json:"title"`
	AuthorName     string     `json:"author_name"`
	AuthorEmail    string     `json:"author_email"`
	// AuthoredDate   *time.Time `json:"authored_date"`
	CommitterName  string     `json:"committer_name"`
	CommitterEmail string     `json:"committer_email"`
	// CommittedDate  *time.Time `json:"committed_date"`
	// CreatedAt      *time.Time `json:"created_at"`
	Message        string     `json:"message"`
	ParentIDs      []string   `json:"parent_ids"`
}

func (c Commit) String() string {
	return Stringify(c)
}

type ListCommitsOptions struct {
	ListOptions
	RefName *string    `json:"ref_name,omitempty" json:"ref_name,omitempty"`
	Since   *time.Time `url:"since,omitempty" json:"since,omitempty"`
	Until   *time.Time `url:"until,omitempty" json:"until,omitempty"`
	Path    *string    `url:"path,omitempty" json:"path,omitempty"`
}

func (s *CommitsService) ListCommits(pid interface{}, opts *ListCommitsOptions) ([]*Commit, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/commits", pathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, opts)
	if err != nil {
		return nil, nil, err
	}

	var c []*Commit
	resp, err := s.client.Do(req, &c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, err
}

type CommitRef struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type GetCommitRefsOptions struct {
	ListOptions
	Type *string `url:"type,omitempty" json:"type,omitempty"`
}

func (s *CommitsService) ListCommitRefs(pid interface{}, sha string, opts *GetCommitRefsOptions) ([]*CommitRef, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/commits/%s/refs", pathEscape(project), pathEscape(sha))

	req, err := s.client.NewRequest(http.MethodGet, u, opts)
	if err != nil {
		return nil, nil, err
	}

	var cs []*CommitRef
	resp, err := s.client.Do(req, &cs)
	if err != nil {
		return nil, resp, err
	}

	return cs, resp, err
}

func (s *CommitsService) GetCommit(pid interface{}, sha string) (*Commit, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	if sha == "" {
		return nil, nil, fmt.Errorf("SHA must be a non-empty string")
	}
	u := fmt.Sprintf("projects/%s/repository/commits/%s", pathEscape(project), url.PathEscape(sha))

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	c := new(Commit)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, err
}
