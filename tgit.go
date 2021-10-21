package tgit

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	defaultBaseURL = "https://git.code.tencent.com/"
	apiVersionPath = "api/v3/"
	userAgent      = "go-tgit"
)

var (
	chinaLoc = time.FixedZone("Asia/Shanghai", int((8 * time.Hour).Seconds()))
)

type authType int

const (
	basicAuth authType = iota
	oAuthToken
	privateToken
)

type Client struct {
	// HTTP client used to communicate with the API.
	client *retryablehttp.Client

	baseURL *url.URL

	// Token types used to make authenticated API calls.
	authType authType

	// Username and password used for basic authentication.
	username, password string

	// Token used to make authenticated API calls.
	token string

	// User agent used when communicating with the TGit API.
	UserAgent string

	// Services used for talking to different parts of the TGit API.
	Branches        *BranchesService
	Commits         *CommitsService
	Repositories    *RepositoriesService
	RepositoryFiles *RepositoryFilesService
	Tags            *TagsService
	Projects        *ProjectsService
	MergeRequests   *MergeRequestsService
}

type ListOptions struct {
	Page    int `url:"page,omitempty" json:"page,omitempty"`
	PerPage int `url:"per_page,omitempty" json:"per_page,omitempty"`
}

func NewClient(hc *retryablehttp.Client, token string) (*Client, error) {
	client, err := newClient(hc)
	if err != nil {
		return nil, err
	}

	client.authType = privateToken
	client.token = token
	return client, nil
}

func NewBasicAuthClient(hc *retryablehttp.Client, username, password string) (*Client, error) {
	client, err := newClient(hc)
	if err != nil {
		return nil, err
	}

	client.authType = basicAuth
	client.username = username
	client.password = password

	return client, nil
}

func NewOAuthClient(hc *retryablehttp.Client, token string) (*Client, error) {
	client, err := newClient(hc)
	if err != nil {
		return nil, err
	}

	client.authType = oAuthToken
	client.token = token
	return client, nil
}

func newClient(hc *retryablehttp.Client) (*Client, error) {
	c := &Client{UserAgent: userAgent}

	c.client = hc

	c.setBaseURL(defaultBaseURL)

	c.Branches = &BranchesService{client: c}
	c.Commits = &CommitsService{client: c}
	c.Repositories = &RepositoriesService{client: c}
	c.RepositoryFiles = &RepositoryFilesService{client: c}
	c.Tags = &TagsService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.MergeRequests = &MergeRequestsService{client: c}

	return c, nil
}

// BaseURL return a copy of the baseURL.
func (c *Client) BaseURL() *url.URL {
	u := *c.baseURL
	return &u
}

// setBaseURL sets the base URL for API requests to a custom encoding.
func (c *Client) setBaseURL(urlStr string) error {
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(baseURL.Path, apiVersionPath) {
		baseURL.Path += apiVersionPath
	}

	// Update the base URL of the client.
	c.baseURL = baseURL

	return nil
}

func (c *Client) NewRequest(method, path string, opt interface{}) (*retryablehttp.Request, error) {
	u := *c.baseURL
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	// Set the encoded path data
	u.RawPath = c.baseURL.Path + path
	u.Path = c.baseURL.Path + unescaped

	// Create a request specific headers map.
	reqHeaders := make(http.Header)
	reqHeaders.Set("Accept", "application/json")

	if c.UserAgent != "" {
		reqHeaders.Set("User-Agent", c.UserAgent)
	}

	var body interface{}
	switch {
	case method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete:
		reqHeaders.Set("Content-Type", "application/json")

		if opt != nil {
			body, err = json.Marshal(opt)
			if err != nil {
				return nil, err
			}
		}
	case opt != nil:
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req, err := retryablehttp.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	return req, nil
}

type Response struct {
	*http.Response

	TotalItems   int
	TotalPages   int
	ItemsPerPage int
	CurrentPage  int
	NextPage     int
	PreviousPage int
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	response.populatePageValues()
	return response
}

const (
	xTotal      = "X-Total"
	xTotalPages = "X-Total-Pages"
	xPerPage    = "X-Per-Page"
	xPage       = "X-Page"
	xNextPage   = "X-Next-Page"
	xPrevPage   = "X-Prev-Page"
)

// populatePageValues parses the HTTP Link response headers and populates the
// various pagination link values in the Response.
func (r *Response) populatePageValues() {
	if totalItems := r.Response.Header.Get(xTotal); totalItems != "" {
		r.TotalItems, _ = strconv.Atoi(totalItems)
	}
	if totalPages := r.Response.Header.Get(xTotalPages); totalPages != "" {
		r.TotalPages, _ = strconv.Atoi(totalPages)
	}
	if itemsPerPage := r.Response.Header.Get(xPerPage); itemsPerPage != "" {
		r.ItemsPerPage, _ = strconv.Atoi(itemsPerPage)
	}
	if currentPage := r.Response.Header.Get(xPage); currentPage != "" {
		r.CurrentPage, _ = strconv.Atoi(currentPage)
	}
	if nextPage := r.Response.Header.Get(xNextPage); nextPage != "" {
		r.NextPage, _ = strconv.Atoi(nextPage)
	}
	if previousPage := r.Response.Header.Get(xPrevPage); previousPage != "" {
		r.PreviousPage, _ = strconv.Atoi(previousPage)
	}
}

func (c *Client) Do(req *retryablehttp.Request, v interface{}) (*Response, error) {

	switch c.authType {
	case oAuthToken:
		req.Header.Set("OAUTH-TOKEN", c.token)
	case privateToken:
		req.Header.Set("PRIVATE-TOKEN", c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}

// Helper function to accept and format both the project ID or name as project
// identifier for all API calls.
func parseID(id interface{}) (string, error) {
	switch v := id.(type) {
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("invalid ID type %#v, the ID must be an int or a string", id)
	}
}

// Helper function to escape a project identifier.
func pathEscape(s string) string {
	return strings.Replace(url.PathEscape(s), ".", "%2E", -1)
}

type ErrorResponse struct {
	Body     []byte
	Response *http.Response
	Message  string
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}

func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Body = data

		var raw interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			errorResponse.Message = "failed to parse unknown error format"
		} else {
			errorResponse.Message = parseError(raw)
		}
	}

	return errorResponse
}

func parseError(raw interface{}) string {
	switch raw := raw.(type) {
	case string:
		return raw

	case []interface{}:
		var errs []string
		for _, v := range raw {
			errs = append(errs, parseError(v))
		}
		return fmt.Sprintf("[%s]", strings.Join(errs, ", "))

	case map[string]interface{}:
		var errs []string
		for k, v := range raw {
			errs = append(errs, fmt.Sprintf("{%s: %s}", k, parseError(v)))
		}
		sort.Strings(errs)
		return strings.Join(errs, ", ")

	default:
		return fmt.Sprintf("failed to parse unexpected error type: %T", raw)
	}
}

type Time struct {
	time.Time
}

func (t *Time) AsTime() *time.Time {
	if t == nil {
		return nil
	}
	date := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	return &date
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	s := string(data)
	if s == "null" {
		return nil
	}

	s = strings.ReplaceAll(s, "+0000", "+08:00")
	t.Time, err = time.Parse(`"`+time.RFC3339+`"`, s)
	return err
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return t.Time.MarshalJSON()
}
