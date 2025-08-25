package worker

import (
	"fmt"
	"os/exec"
)

type Builder struct {}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) BuildServer() error {
	
	cmd := exec.Command("echo", "hello")

	// Run and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	fmt.Println("Output:", string(output))
	return nil
}
