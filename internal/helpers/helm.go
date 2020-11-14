package helpers

import (
	"regexp"

	"github.com/gojek/stevedore/internal/cli/helm"
)

var helmRepoRegex = regexp.MustCompile(`(?m)^([\w-]*)\s*([\w:./-]*)`)

// AddHelmRepos if not already added
func AddHelmRepos(helmRepos helm.Repos) error {
	existingRepos, err := existingHelmRepos()
	if err != nil {
		return err
	}

	diff := helmRepos.Diff(existingRepos)
	for _, repo := range diff {
		command, err := helm.RepoAdd(repo)
		if err != nil {
			return err
		}

		_, err = execute(command)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateHelmRepo updates local helm repo cache
func UpdateHelmRepo() error {
	command, err := helm.RepoUpdate()
	if err != nil {
		return err
	}

	_, err = execute(command)
	return err
}

func existingHelmRepos() (helm.Repos, error) {
	var result helm.Repos

	command, err := helm.RepoList()
	if err != nil {
		return result, err
	}

	out, err := execute(command)
	if err != nil {
		return result, err
	}

	matches := helmRepoRegex.FindAllStringSubmatch(out, -1)
	for _, match := range matches[1:] {
		if len(match) != 3 {
			continue
		}

		result = append(result, helm.Repo{
			Name: match[1],
			URL:  match[2],
		})
	}
	return result, nil
}
