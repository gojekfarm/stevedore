package helm

import "github.com/cucumber/godog/gherkin"

// Repo represents helm repo
type Repo struct {
	Name string
	URL  string
}

// Repos represents list of HelmRepo
type Repos []Repo

// Contains returns if the item is contained
func (repos Repos) Contains(item Repo) bool {
	for _, repo := range repos {
		if repo.Name == item.Name && repo.URL == item.URL {
			return true
		}
	}
	return false
}

// Diff returns the HelmRepos which are not contained in the source list
func (repos Repos) Diff(with Repos) Repos {
	var result Repos

	for _, repo := range repos {
		if !with.Contains(repo) {
			result = append(result, repo)
		}
	}
	return result
}

// NewRepos create list of HelmRepo from gherkin.DataTable
func NewRepos(helmRepos *gherkin.DataTable) Repos {
	result := make(Repos, 0, len(helmRepos.Rows)-1)
	for _, row := range helmRepos.Rows[1:] {
		if len(row.Cells) != 2 {
			continue
		}

		result = append(result, Repo{
			Name: row.Cells[0].Value,
			URL:  row.Cells[1].Value,
		})
	}

	return result
}
