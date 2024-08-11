package ci

type Workspace interface {
	Branch() string
	Commit() string
	Dir() string
	Env() []string
}
