package stevedore

// IgnoredRelease represents necessary information for
// ignoring a release
type IgnoredRelease struct {
	Name   string `validate:"required"`
	Reason string
}

// IgnoredReleases is a list of release which needs to ignored
type IgnoredReleases []IgnoredRelease

// Names returns list of release names
func (releases IgnoredReleases) Names() []string {
	var result []string

	for _, release := range releases {
		result = append(result, release.Name)
	}
	return result
}

// Find returns IgnoredComponent matched by name and the condition representing match
func (releases IgnoredReleases) Find(name string) (IgnoredRelease, bool) {
	for _, ignoreComponent := range releases {
		if ignoreComponent.Name == name {
			return ignoreComponent, true
		}
	}
	return IgnoredRelease{}, false
}
