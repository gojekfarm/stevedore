package helm

// Resource represents name and kind of a k8s resource
type Resource struct {
	// Name of the k8s resource
	Name string
	// Kind of the k8s resource (Eg: Deployment,Service,Job,etc)
	Kind string
}

// Resources represents a collection of k8s resource
type Resources []Resource

// GroupByKind groups the resources based on their Kind
func (resources Resources) GroupByKind() map[string]Resources {
	group := make(map[string]Resources)
	for _, resource := range resources {
		kind := resource.Kind
		group[kind] = append(group[kind], resource)
	}
	return group
}

// Summary represents all changes to a helm release
type Summary struct {
	Added     Resources
	Modified  Resources
	Destroyed Resources
}
