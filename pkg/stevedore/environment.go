package stevedore

// Environment to which Stevedore Manifest is to be applied
type Environment string

// Environments is a collection of environment
type Environments []Environment

// NewEnvironments returns Environments created from given string array
func NewEnvironments(environmentsList []string) Environments {
	var environments Environments
	for _, env := range environmentsList {
		environments = append(environments, Environment(env))
	}
	return environments
}

// Contains returns true if environment is present
func (environments Environments) Contains(environment Environment) bool {
	for _, item := range environments {
		if item == environment {
			return true
		}
	}
	return false
}
