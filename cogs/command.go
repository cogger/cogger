package cogs

import (
	"bytes"
	"os/exec"

	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

func Command(ctx context.Context, command string, arg ...string) cogger.Cog {
	return ExecuteCommand(ctx, exec.Command(command, arg...))
}

func CommandWithOutput(ctx context.Context, out *bytes.Buffer, command string, arg ...string) cogger.Cog {
	cmd := exec.Command(command, arg...)
	cmd.Stdout = out

	return ExecuteCommand(ctx, cmd)
}

func ExecuteCommand(ctx context.Context, cmd *exec.Cmd) cogger.Cog {
	return Simple(ctx, func() error {
		return cmd.Run()
	})
}
