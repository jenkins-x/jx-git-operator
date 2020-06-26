package constants

const (
	// DefaultSelectorKey selector key
	DefaultSelectorKey = "git-operator.jenkins.io/kind"

	// DefaultSelectorValue selector value
	DefaultSelectorValue = "git-operator"

	// DefaultSelector default selector for Secrets
	DefaultSelector = DefaultSelectorKey + "=" + DefaultSelectorValue
)
