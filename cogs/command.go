package cogs

import (
	"bytes"
	"os/exec"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
)

//Command creates a cog that runs the command with the arguments provided
func Command(ctx context.Context, command string, arg ...string) cogger.Cog {
	return ExecuteCommand(ctx, exec.Command(command, arg...))
}

//CommandWithOutput creates a cog that runs the command with the arguments provided outputting the return from the command to the buffer
func CommandWithOutput(ctx context.Context, out *bytes.Buffer, command string, arg ...string) cogger.Cog {
	cmd := exec.Command(command, arg...)
	cmd.Stdout = out

	return ExecuteCommand(ctx, cmd)
}

//ExecuteCommand creates a cog that runs the cmd provided
func ExecuteCommand(ctx context.Context, cmd *exec.Cmd) cogger.Cog {
	return Simple(ctx, func() error {
		return cmd.Run()
	})
}
