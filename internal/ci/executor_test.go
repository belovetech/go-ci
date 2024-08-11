package ci

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWorkspace struct {
	mock.Mock
}

func (ws *mockWorkspace) Branch() string {
	args := ws.Called()
	return args.String(0)
}

func (ws *mockWorkspace) Commit() string {
	args := ws.Called()
	return args.String(0)
}

func (ws *mockWorkspace) Dir() string {
	args := ws.Called()
	return args.String(0)
}

func (ws *mockWorkspace) Env() []string {
	args := ws.Called()
	return args.Get(0).([]string)
}

func (ws *mockWorkspace) LoadPipeline() (*Pipeline, error) {
	args := ws.Called()
	return args.Get(0).(*Pipeline), args.Error(1)
}

func (ws *mockWorkspace) ExecuteCommand(ctx context.Context, cmd string, arguements []string) ([]byte, error) {
	args := ws.Called(ctx, cmd, arguements)
	return args.Get(0).([]byte), args.Error(1)
}

func TestRunDefaultHappyPath(t *testing.T) {
	wsMock := mockWorkspace{}

	wsMock.On("LoadPipeline").Return(
		&Pipeline{
			Name: "Test pipeline",
			Steps: []Step{
				{Name: "Step 1", Commands: []string{"cmd1 arg1 arg2"}},
			},
		}, nil,
	)

	wsMock.On("ExecuteCommand", mock.Anything, "cmd1", []string{"arg1", "arg2"}).Return([]byte("output"), nil)

	executor := NewExecutor(&wsMock)
	str, err := executor.RunDefault(context.Background())
	assert.Nil(t, err)
	expectedOutput := "Executing pipeline: Test pipeline\nExecuting step: Step 1\noutput\n"
	assert.Equal(t, expectedOutput, str)
}
