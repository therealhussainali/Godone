package libgodone

import (
	"fmt"
	"os"
)

// Struct used for command line argument parsing.
type CmdArgs struct {
	operation    Operation
	name         string
	compiler     Compiler
	language     Language
	target       Target
	crossCompile bool
}

// Parses all command line arguments into a CmdArgs struct.
func ParseArgs() *CmdArgs {
	args := CmdArgs{
		operation:    HELP,
		name:         GetwdBaseName(),
		compiler:     GCC,
		language:     CPP,
		target:       LINUX,
		crossCompile: false,
	}

	for i := 0; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "new":
			args.operation = NEW
			if i+1 >= len(os.Args) {
				Die(USAGE, fmt.Sprintf("'%s' option expects an argument which wasn't provided", os.Args[i]))
			}
			args.name = os.Args[i+1]
			if len(args.name) == 0 {
				Die(USAGE, "The name of your new project cannot be empty!")
			}
			i++

		case "run":
			args.operation = RUN
			if i+1 >= len(os.Args) {
				Die(USAGE, fmt.Sprintf("'%s' option expects an argument which wasn't provided", os.Args[i]))
			}
			target, err := StringAsTarget(os.Args[i+1])
			if err != nil {
				Die(USAGE, err.Error())
			}
			args.target = target
			i++

		case "build":
			args.operation = BUILD
			if i+1 >= len(os.Args) {
				Die(USAGE, fmt.Sprintf("'%s' option expects an argument which wasn't provided", os.Args[i]))
			}
			target, err := StringAsTarget(os.Args[i+1])
			if err != nil {
				Die(USAGE, err.Error())
			}
			args.target = target
			i++

		case "buildrun":
			args.operation = BUILDRUN
			if i+1 >= len(os.Args) {
				Die(USAGE, fmt.Sprintf("'%s' option expects an argument which wasn't provided", os.Args[i]))
			}
			target, err := StringAsTarget(os.Args[i+1])
			if err != nil {
				Die(USAGE, err.Error())
			}
			args.target = target
			i++

		case "clean":
			args.operation = CLEAN
			if i+1 >= len(os.Args) {
				Die(USAGE, fmt.Sprintf("'%s' option expects an argument which wasn't provided", os.Args[i]))
			}
			target, err := StringAsTarget(os.Args[i+1])
			if err != nil {
				Die(USAGE, err.Error())
			}
			args.target = target
			i++

		case "init":
			args.operation = INIT

		case "-c", "--comp", "--compiler":
			if i+1 >= len(os.Args) {
				Die(USAGE, fmt.Sprintf("'%s' option expects an argument which wasn't provided", os.Args[i]))
			}
			compiler, err := StringAsCompiler(os.Args[i+1])
			if err != nil {
				Die(USAGE, err.Error())
			}
			args.compiler = compiler
			i++

		case "-l", "--lang", "--language":
			if i+1 >= len(os.Args) {
				Die(USAGE, fmt.Sprintf("'%s' option expects an argument which wasn't provided", os.Args[i]))
			}
			language, err := StringAsLanguage(os.Args[i+1])
			if err != nil {
				Die(USAGE, err.Error())
			}
			args.language = language
			i++

		case "-x", "--cross", "--cross-compile":
			args.crossCompile = true
		}
	}

	return &args
}
