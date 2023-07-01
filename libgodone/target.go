package libgodone

import (
	"fmt"
	"os"
	"os/exec"
)

type Target int

const (
	LINUX   Target = iota
	WINDOWS Target = iota
)

func TargetAsString(target Target) string {
	switch target {
	case LINUX:
		return "linux"

	case WINDOWS:
		return "windows"

	default:
		return "linux"
	}
}

func StringAsTarget(target string) (Target, error) {
	switch target {
	case "linux", "Linux", "lin", "Lin":
		return LINUX, nil

	case "windows", "Windows", "win", "Win":
		return WINDOWS, nil

	default:
		return LINUX, fmt.Errorf("unknown target")
	}
}

func RunTarget(target Target) {
	cmd := exec.Command("make", "run")
	cmd.Dir = FullPathToBuildTarget(target)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		Die(RUNTIME, fmt.Sprintf("Unable to run target \"%s\"\n", TargetAsString(target)))
	}
}

func BuildTarget(target Target) {
	cmd := exec.Command("make", "build")
	cmd.Dir = FullPathToBuildTarget(target)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		Die(RUNTIME, fmt.Sprintf("Unable to build target \"%s\"\n", TargetAsString(target)))
	}
}

func BuildRunTarget(target Target) {
	cmd := exec.Command("make", "build", "run")
	cmd.Dir = FullPathToBuildTarget(target)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		Die(RUNTIME, fmt.Sprintf("Unable to build and run target \"%s\"", TargetAsString(target)))
	}
}

func CleanTarget(target Target) {
	cmd := exec.Command("make", "clean")
	cmd.Dir = FullPathToBuildTarget(target)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		Die(RUNTIME, fmt.Sprintf("Unable to clean target \"%s\"", TargetAsString(target)))
	}
}
