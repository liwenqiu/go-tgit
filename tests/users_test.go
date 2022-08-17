package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/liwenqiu/go-tgit"
)

func TestUsersService_List(t *testing.T) {
	hc := retryablehttp.NewClient()
	c, _ := tgit.NewClient(hc, "4n-Bvu7L9OI8dVuTTT92")

	r, s, err := c.Users.Get("")
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println(r)
	fmt.Println(s.StatusCode)
}
