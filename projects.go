package tgit

import (
	"net/http"
)

type ProjectsService struct {
	client *Client
}

type VisibilityLevelValue int

const (
	PrivateLevel  VisibilityLevelValue = 0
	InternalLevel VisibilityLevelValue = 10
	PublicLevel   VisibilityLevelValue = 20
)

type ProjectOrderByValue string

const (
	OrderByID             ProjectOrderByValue = "id"
	OrderByName           ProjectOrderByValue = "name"
	OrderByPath           ProjectOrderByValue = "path"
	OrderByCreatedAt      ProjectOrderByValue = "created_at"
	OrderByUpdatedAt      ProjectOrderByValue = "updated_at"
	OrderByLastActivityAt ProjectOrderByValue = "last_activity_at"
)

type SortValue string

const (
	Asc  SortValue = "asc"
	Desc SortValue = "desc"
)

type ListProjectsOptions struct {
	ListOptions
	Search          *string               `url:"search,omitempty" json:"search,omitempty"`
	WithArchived    *bool                 `url:"with_archived,omitempty" json:"with_archived,omitempty"`
	WithPush        *bool                 `url:"with_push,omitempty" json:"with_push,omitempty"`
	Abandoned       *bool                 `url:"abandoned,omitempty" json:"abandoned,omitempty"`
	VisibilityLevel *VisibilityLevelValue `url:"visibility_level,omitempty" json:"visibility_level,omitempty"`
	OrderBy         *ProjectOrderByValue  `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort            *SortValue            `url:"sort,omitempty" json:"sort,omitempty"`
}

type ProjectNamespace struct {
	CreatedAt   *Time  `json:"created_at"`
	Description string `json:"description"`
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	OwnerID     int64  `json:"owner_id"`
	Path        string `json:"path"`
	UpdatedAt   *Time  `json:"updated_at"`
}

type ProjectConfigStorage struct {
	LimitLfsFileSize float64 `json:"limit_lfs_file_size"`
	LimitSize        float64 `json:"limit_size"`
	LimitFileSize    float64 `json:"limit_file_size"`
	LimitLfsSize     float64 `json:"limit_lfs_size"`
}

type ProjectStatistics struct {
	CommitCount    int     `json:"commit_count"`
	RepositorySize float64 `json:"repository_size"`
}

type Permission struct {
	AccessLevel int `json:"access_level"`
}

type ProjectPermission struct {
	ProjectAccess    Permission `json:"project_access"`
	ShareGroupAccess Permission `json:"share_group_access"`
	GroupAccess      Permission `json:"group_access"`
}

type ProjectItem struct {
	ID                        int64                 `json:"id"`
	Description               string                `json:"description"`
	Public                    bool                  `json:"public"`
	Archived                  bool                  `json:"archived"`
	VisibilityLevel           VisibilityLevelValue  `json:"visibility_level"`
	PublicVisibility          int                   `json:"public_visibility"`
	Namespace                 *ProjectNamespace     `json:"namespace"`
	Owner                     *User                 `json:"owner"`
	Name                      string                `json:"name"`
	NameWithNamespace         string                `json:"name_with_namespace"`
	Path                      string                `json:"path"`
	PathWithNamespace         string                `json:"path_with_namespace"`
	DefaultBranch             string                `json:"default_branch"`
	SSHURLToRepo              string                `json:"ssh_url_to_repo"`
	HTTPURLToRepo             string                `json:"http_url_to_repo"`
	HTTPSURLToRepo            string                `json:"https_url_to_repo"`
	WebURL                    string                `json:"web_url"`
	TagList                   []string              `json:"tag_list"`
	IssuesEnabled             bool                  `json:"issues_enabled"`
	MergeRequestsEnabled      bool                  `json:"merge_requests_enabled"`
	WikiEnabled               bool                  `json:"wiki_enabled"`
	SnippetsEnabled           bool                  `json:"snippets_enabled"`
	ReviewEnabled             bool                  `json:"review_enabled"`
	ForkEnabled               bool                  `json:"fork_enabled"`
	TagNameRegex              string                `json:"tag_name_regex"`
	TagCreatePushLevel        int                   `json:"tag_create_push_level"`
	BranchNameRegex           string                `json:"branch_name_regex"`
	CreatedAt                 *Time                 `json:"created_at"`
	LastActivityAt            *Time                 `json:"last_activity_at"`
	CreatorID                 int64                 `json:"creator_id"`
	AvatarURL                 string                `json:"avatar_url"`
	WatchsCount               int                   `json:"watchs_count"`
	StarsCount                int                   `json:"stars_count"`
	ForksCount                int                   `json:"forks_count"`
	ConfigStorage             *ProjectConfigStorage `json:"config_storage"`
	ForkedFromProject         string                `json:"forked_from_project"`
	Statistics                *ProjectStatistics    `json:"statistics"`
	Permissions               *ProjectPermission    `json:"permissions"`
	SuggestionReviewers       []*string             `json:"suggestion_reviewers"`
	NecessaryReviewers        []*string             `json:"necessary_reviewers"`
	PathReviewerRules         string                `json:"path_reviewer_rules"`
	ApproverRule              int                   `json:"approver_rule"`
	NecessaryApproverRule     int                   `json:"necessary_approver_rule"`
	CanApproveByCreator       bool                  `json:"can_approve_by_creator"`
	AutoCreateReviewAfterPush bool                  `json:"auto_create_review_after_push"`
	ForbiddenModifyRule       bool                  `json:"forbidden_modify_rule"`
	PushResetEnabled          bool                  `json:"push_reset_enabled"`
	MergeRequestTemplate      string                `json:"merge_request_template"`
	FileOwnerPathRules        string                `json:"file_owner_path_rules"`
}

// ListProjects https://code.tencent.com/help/api/project#searchProjectByName
func (s *ProjectsService) ListProjects(opts *ListProjectsOptions) ([]*ProjectItem, *Response, error) {
	url := "projects"
	req, err := s.client.NewRequest(http.MethodGet, url, opts)
	if err != nil {
		return nil, nil, err
	}

	var p []*ProjectItem
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}
