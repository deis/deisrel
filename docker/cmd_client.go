package docker

import (
	"os/exec"
)

type cmdClient struct{}

// NewCmdClient creates a new Client that does all of its operations by shelling out to the docker CLI
func NewCmdClient() Client {
	return &cmdClient{}
}

func (c *cmdClient) Push(img *Image) error {
	cmd := exec.Command("docker", "push", img.String())
	return cmd.Run()
}

func (c *cmdClient) Pull(img *Image) error {
	cmd := exec.Command("docker", "pull", img.String())
	return cmd.Run()
}

func (c *cmdClient) Retag(src *Image, tar *Image) error {
	cmd := exec.Command("docker", "tag", src.String(), tar.String())
	return cmd.Run()
}
