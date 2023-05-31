package tgit

import (
	"fmt"
	"net/http"
	"time"
)

type MergeRequestsService struct {
	client *Client
}

type ListMergeRequestsOptions struct {
	Id      string `json:"id,omitempty"`
	Iid     int64  `url:"iid,omitempty" json:"iid,omitempty"`
	State   string `url:"state,omitempty" json:"state,omitempty"`
	OrderBy string `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort    string `url:"sort,omitempty" json:"sort,omitempty"`
	Page    int    `url:"page,omitempty" json:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty" json:"per_page,omitempty"`
}

type MergeRequestUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	WebURL    string `json:"web_url"`
	Name      string `json:"name"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
}

type Milestone struct {
	ID          int64  `json:"id"`
	ProjectID   int64  `json:"project_id"`
	Title       string `json:"title"`
	State       string `json:"state"`
	Iid         int64  `json:"iid"`
	DueDate     string `json:"due_date"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Description string `json:"description"`
}

type MergeRequestViewer struct {
	Type           string `json:"type"`
	ReviewState    string `json:"review_state"`
	ReviewDuration int32  `json:"review_duration"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	WebURL         string `json:"web_url"`
	Name           string `json:"name"`
	State          string `json:"state"`
	AvatarURL      string `json:"avatar_url"`
}

type MergeRequest struct {
	Labels              []string              `json:"labels"`
	ID                  int64                 `json:"id"`
	Title               string                `json:"title"`
	TargetProjectID     int64                 `json:"target_project_id"`
	TargetBranch        string                `json:"target_branch"`
	SourceProjectID     int64                 `json:"source_project_id"`
	SourceBranch        string                `json:"source_branch"`
	State               string                `json:"state"`
	MergeStatus         string                `json:"merge_status"`
	Iid                 int64                 `json:"iid"`
	Description         string                `json:"description"`
	CreatedAt           string                `json:"created_at"`
	UpdatedAt           string                `json:"updated_at"`
	ResolvedAt          *time.Time            `json:"resolved_at"`
	MergeType           string                `json:"merge_type"`
	Assignee            *MergeRequestUser     `json:"assignee"`
	Author              *MergeRequestUser     `json:"author"`
	MergeCommitSha      string                `json:"merge_commit_sha"`
	Milestone           *Milestone            `json:"milestone"`
	NecessaryReviewers  []*MergeRequestViewer `json:"necessary_reviewers"`
	SuggestionReviewers []*MergeRequestViewer `json:"suggestion_reviewers"`
	BaseCommit          string                `json:"base_commit"`
	TargetCommit        string                `json:"target_commit"`
	SourceCommit        string                `json:"source_commit"`
	ProjectID           int64                 `json:"project_id"`
	WorkInProgress      bool                  `json:"work_in_progress"`
	Upvotes             int                   `json:"upvotes"`
	Downvotes           int                   `json:"downvotes"`
}

// ListMergeRequests https://code.tencent.com/help/api/mergeRequest#getMergeRequests
func (s *MergeRequestsService) ListMergeRequests(opts *ListMergeRequestsOptions) ([]*MergeRequest, *Response, error) {
	url := fmt.Sprintf("projects/%s/merge_requests", opts.Id)

	req, err := s.client.NewRequest(http.MethodGet, url, opts)
	if err != nil {
		return nil, nil, err
	}

	var m []*MergeRequest
	resp, err := s.client.Do(req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

type ListMergeRequestChangeOptions struct {
	Id             string
	MergeRequestId int
}

type DiffFile struct {
	OldPath     string `json:"old_path"`
	NewPath     string `json:"new_path"`
	AMode       int    `json:"a_mode"`
	BMode       int    `json:"b_mode"`
	Diff        string `json:"diff"`
	NewFile     bool   `json:"new_file"`
	RenamedFile bool   `json:"renamed_file"`
	DeletedFile bool   `json:"deleted_file"`
	IsTooLarge  bool   `json:"is_too_large"`
	IsCollapse  bool   `json:"is_collapse"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
}

type MergeRequestChange struct {
	Labels              []string              `json:"labels"`
	ID                  int64                 `json:"id"`
	Title               string                `json:"title"`
	TargetProjectID     int64                 `json:"target_project_id"`
	TargetBranch        string                `json:"target_branch"`
	SourceProjectID     int64                 `json:"source_project_id"`
	SourceBranch        string                `json:"source_branch"`
	State               string                `json:"state"`
	MergeStatus         string                `json:"merge_status"`
	Iid                 int64                 `json:"iid"`
	Description         string                `json:"description"`
	CreatedAt           string                `json:"created_at"`
	UpdatedAt           string                `json:"updated_at"`
	ResolvedAt          *time.Time            `json:"resolved_at"`
	MergeType           string                `json:"merge_type"`
	Assignee            *MergeRequestUser     `json:"assignee"`
	Author              *MergeRequestUser     `json:"author"`
	MergeCommitSha      string                `json:"merge_commit_sha"`
	Milestone           *Milestone            `json:"milestone"`
	NecessaryReviewers  []*MergeRequestViewer `json:"necessary_reviewers"`
	SuggestionReviewers []*MergeRequestViewer `json:"suggestion_reviewers"`
	BaseCommit          string                `json:"base_commit"`
	TargetCommit        string                `json:"target_commit"`
	SourceCommit        string                `json:"source_commit"`
	ProjectID           inint64t              `json:"project_id"`
	WorkInProgress      bool                  `json:"work_in_progress"`
	Upvotes             int                   `json:"upvotes"`
	Downvotes           int                   `json:"downvotes"`
	Files               []*DiffFile           `json:"files"`
}

// ListMergeRequestChange https://code.tencent.com/help/api/mergeRequest#searchMergeRequest
func (s *MergeRequestsService) ListMergeRequestChange(opts *ListMergeRequestChangeOptions) (*MergeRequestChange, *Response, error) {
	url := fmt.Sprintf("projects/:%s/merge_request/%d/changes", opts.Id, opts.MergeRequestId)

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	var c *MergeRequestChange
	resp, err := s.client.Do(req, &c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}
