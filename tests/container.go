package tests

import (
	"fmt"
	"os/exec"
)

// remove the container with the given name
func RemoveContainer(name string) {
	c := exec.Command("docker", []string{"rm", "-f", name}...)
	err := c.Run()
	if err != nil {
		fmt.Println(fmt.Sprintf("got an error %+v when running this command: %s in order to remove a container", err, c.Args))
	}
}
