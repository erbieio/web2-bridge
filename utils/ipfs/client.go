package ipfs

import (
	shell "github.com/ipfs/go-ipfs-api"
)

func NewClient(api string) *shell.Shell {
	return shell.NewShell(api)
}
