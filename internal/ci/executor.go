package ci

import (
	"context"
	"os/exec"
	"strings"
)

type Step struct {
	Name     string   `yaml:"name"`
	Commands []string `yaml:"commands"`
}

type Pipeline struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}
type Workspace interface {
	Branch() string
	Commit() string
	Dir() string
	Env() []string
	LoadPipeline() (*Pipeline, error)
	ExecuteCommand(ctx context.Context, cmd string, args []string) ([]byte, error)
}

type Executor struct {
	ws Workspace
}

func NewExecutor(ws Workspace) *Executor {
	return &Executor{ws: ws}
}

func (e *Executor) Run(ctx context.Context, pipeline *Pipeline) (string, error) {
	output := strings.Builder{}
	output.WriteString("Executing pipeline: " + pipeline.Name + "\n")

	for _, step := range pipeline.Steps {
		output.WriteString("Executing step: " + step.Name + "\n")
		for _, cmd := range step.Commands {
			withArgs := strings.Fields(cmd)
			cmd := withArgs[:1][0]
			args := withArgs[1:]
			out, err := e.ws.ExecuteCommand(ctx, cmd, args)
			output.Write(out)
			output.WriteRune('\n')
			if err != nil {
				return output.String(), err
			}
		}
	}
	return output.String(), nil
}
func (e *Executor) RunDefault(ctx context.Context) (string, error) {
	pipeline, err := e.ws.LoadPipeline()
	if err != nil {
		return "", err
	}
	return e.Run(ctx, pipeline)
}

func (ws *workspaceImpl) ExecuteCommand(ctx context.Context, cmd string, args []string) ([]byte, error) {
	command := exec.Command(cmd, args...)
	command.Dir = ws.Dir()
	command.Env = append(command.Environ(), ws.Env()...)
	return command.CombinedOutput()
}
