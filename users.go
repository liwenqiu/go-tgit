package tgit

import (
	"fmt"
	"net/http"
	"strings"
)

type UsersService struct {
	client *Client
}

type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	WebURL    string `json:"web_url"`
	IsAdmin   bool   `json:"is_admin"`
	Bio       string `json:"bio"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
}

// Get fetches a user. Passing the empty string will fetch the authenticated user.
// tgit doc: https://code.tencent.com/help/api/user
func (s *UsersService) Get(user string) (*User, *Response, error) {
	url := "user"
	if strings.TrimSpace(user) != "" {
		url = fmt.Sprintf("users/%s", strings.TrimSpace(user))
	}
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	var p *User
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}
