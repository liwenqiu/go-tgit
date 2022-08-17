package tgit

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/liwenqiu/go-tgit"
)

func TestProjectsService_ListProjects(t *testing.T) {
	hc := retryablehttp.NewClient()
	xc, _ := tgit.NewClient(hc, "IF9fJGOctqgsDw38SKq4")
	repos, s, err := xc.Projects.ListProjects(&tgit.ListProjectsOptions{
		ListOptions: tgit.ListOptions{Page: 1, PerPage: 30},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(s.StatusCode)
	for _, repo := range repos {
		fmt.Println(repo.ID, repo.Name, repo.NameWithNamespace, repo.Description, repo.DefaultBranch, repo.HTTPSURLToRepo, repo.SSHURLToRepo, repo.WebURL)
	}
}
