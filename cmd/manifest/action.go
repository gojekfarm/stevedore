package manifest

// Action type represents manifest related actions
type Action interface {
	Do() (Info, error)
}

// NewAction to create action based on command name
func NewAction(cmd Command, info Info) Action {
	switch cmd.name {
	case applyCommand:
		return NewHelmAction(info, cmd.kubeconfig, false, false, false, cmd.helmRepoName,
			cmd.helmTimeout, cmd.helmAtomic)
	case planCommand:
		return NewHelmAction(info, cmd.kubeconfig, true, true, true, cmd.helmRepoName,
			cmd.helmTimeout, false)
	default:
		return RenderAction{info: info}
	}
}
