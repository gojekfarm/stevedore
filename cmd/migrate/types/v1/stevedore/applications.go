package stevedore

import "github.com/gojek/stevedore/pkg/stevedore"

// Applications is a collection of helm release specification
type Applications []Application

// ReleaseSpecifications convert applications to stevedore.ReleaseSpecifications
func (applications Applications) ReleaseSpecifications() stevedore.ReleaseSpecifications {
	result := stevedore.ReleaseSpecifications{}
	for _, application := range applications {
		result = append(result, application.ReleaseSpecification())
	}
	return result
}
