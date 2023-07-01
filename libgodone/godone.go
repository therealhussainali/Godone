package libgodone

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Execute() {
	args := ParseArgs()

	switch args.operation {
	case NEW, INIT:
		InitFolder(args)

	case RUN:
		UpdateMakefileSources(args.target)
		RunTarget(args.target)

	case BUILD:
		UpdateMakefileSources(args.target)
		BuildTarget(args.target)

	case BUILDRUN:
		UpdateMakefileSources(args.target)
		BuildRunTarget(args.target)

	case CLEAN:
		UpdateMakefileSources(args.target)
		CleanTarget(args.target)

	default:
		fmt.Println("Godone Command Line Utility By MrMadPie.")
		fmt.Println("\nOperations:")
		fmt.Println("help (or anything that isn't a known command) - brings up this message")
		fmt.Println("new <project name> <flags>                    - creates a project in a new folder with name <project name> and flags <flags>")
		fmt.Println("init <flags>                                  - creates a project in the current folder with flags <flags>")
		fmt.Println("build <target>                                - builds for the specified target (available: Linux and Windows if you've enabled cross-compilation)")
		fmt.Println("run <target>                                  - runs the binary of target if it exists (you must have built at least once beforehand, check build/buildrun)")
		fmt.Println("buildrun <target>                             - does the two above one after another so that you don't have to do it yourself")
		fmt.Println("clean <target>                                - removes all object files as well as the produced binary of target")
		fmt.Println("\nFlags:")
		fmt.Println("-c, --comp, --compiler <compiler>             - select your compiler (available: GCC, CLANG, default: GCC)")
		fmt.Println("-l, --lang, --language <language>             - select your language (available: C, CPP, default: CPP)")
		fmt.Println("-x, --cross, --cross-compile                  - enable cross-compilation to windows")
	}
}

func GenerateReadMeString(projectName string) string {
	projectVersion := "0.0.1"
	projectDate := CurrentTimeAndDate()
	projectDescription := "<Insert Project Description Here>"

	return fmt.Sprintf("Name: %s\nVersion: %s\nCreation Date: %s\nDescription: %s\n",
		projectName, projectVersion, projectDate, projectDescription)
}

func GenerateHelloWorldString(language Language) string {
	var ioInclude, mainFunc, printStatement string

	switch language {
	case C:
		ioInclude = "#include <stdio.h>"
		mainFunc = "int main(int argc, char *argv[])"
		printStatement = "puts(\"Hello, World!\");"

	default:
		ioInclude = "#include <iostream>"
		mainFunc = "int main(int argc, char* argv[])"
		printStatement = "std::cout << \"Hello, World!\" << std::endl;"
	}

	return fmt.Sprintf("%s\n\n%s {\n\t%s\n\treturn 0;\n}\n", ioInclude, mainFunc, printStatement)
}

func GenerateMakefile(specifications CmdArgs, forWindows bool) string {
	compiler := GetCompilerVariant(specifications.language, specifications.compiler)

	if forWindows {
		specifications.name += ".exe" // windows executable suffix
		compiler = GetCrossCompilerVariant(specifications.language)
	}

	// variables
	makefile := fmt.Sprintf("CC = %s\nTARGET = %s\nCC_FLAGS = -c -Wall -Wextra --pedantic-errors\nLD_FLAGS = -o $(TARGET)\n\n",
		compiler, specifications.name)

	// all unit which just calls the build unit
	makefile += "all: build\n\n"

	// build unit which calls the compilation of all files and then links them
	makefile += "build: main\n\t$(CC) out/*.o $(LD_FLAGS)\n\n"

	// run unit which runs the executable
	if forWindows {
		makefile += "run:\n\twine ./$(TARGET)\n\n"
	} else {
		makefile += "run:\n\t./$(TARGET)\n\n"
	}

	// unnamed unit that compiles every file provided in the build unit
	makefile += fmt.Sprintf("%%:\n\t$(CC) $(CC_FLAGS) ../../src/$@.%s -o out/$@.o\n\n",
		LanguageAsString(specifications.language))

	// clean unit
	makefile += "clean:\n\trm -f out/*.o $(TARGET)\n\techo \"CLEANED TARGET!\"\n"

	return makefile
}

// Adds newFile in the build unit of the target makefile.
func InsertNewFileToMakefile(target Target, newFile string) {
	makefilePath := FullPathToBuildTarget(target) + "/Makefile"

	data, err := os.ReadFile(makefilePath)
	if err != nil {
		Die(RUNTIME, "Unable to insert new source file into Makefile!")
	}

	newFileWithSpace := " " + newFile
	substring := "build: main"
	dataString := string(data)

	index := strings.Index(dataString, substring) + len(substring)
	newMakefile := dataString[:index] + newFileWithSpace + dataString[index:]

	os.WriteFile(makefilePath, []byte(newMakefile), os.ModePerm)
}

// Removes deadFile from the build unit of the target makefile as well as from the out directory.
func RemoveDeadFileFromMakefile(target Target, deadFile string) {
	makefilePath := FullPathToBuildTarget(target) + "/Makefile"

	file, err := os.Open(makefilePath)
	if err != nil {
		Die(RUNTIME, "Unable to update target Makefile!")
	}
	defer file.Close()

	newMakefile := ""
	substring := "build: main"
	scanner := bufio.NewScanner(file)
	foundSubstring := false

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), substring) {
			foundSubstring = true
			newMakefile += strings.ReplaceAll(scanner.Text(), deadFile, "") + "\n"
		} else {
			newMakefile += scanner.Text() + "\n"
		}
	}

	if !foundSubstring {
		Die(RUNTIME, "Unable to update target Makefile!")
	}

	os.WriteFile(makefilePath, []byte(newMakefile), os.ModePerm)

	// removes the object file from the out directory
	os.Remove(fmt.Sprintf("%s/out/%s.o", FullPathToBuildTarget(target), deadFile))
}

// Checks all files in src/ and updates the makefiles of all targets accordingly.
func UpdateMakefileSources(target Target) {
	sourceFiles := CollectAllFilesFromDir("src/")
	makefileFiles := CollectAllFilesFromMakefile(target)

	// remove .c/.cpp/.cc etc suffixes from the source files
	for i := 0; i < len(sourceFiles); i++ {
		sourceFiles[i] = strings.Split(sourceFiles[i], ".")[0]
	}

	// find source files which need to be added in the makefile and insert them
	for _, file := range sourceFiles {
		if !Contains(makefileFiles, file) {
			InsertNewFileToMakefile(target, file)
		}
	}

	// find source files which no longer exist and need to be removed from the makefile and purge them
	for _, file := range makefileFiles {
		if !Contains(sourceFiles, file) {
			RemoveDeadFileFromMakefile(target, file)
		}
	}
}

// Setup the contents inside of project/build/.
func InitBuildFolder(path string, specifications *CmdArgs, forWindows bool) {
	os.Mkdir(path, os.ModePerm)
	os.Mkdir(path+"/out/", os.ModePerm)

	makefile, err := os.Create(path + "/Makefile")
	if err != nil {
		Die(RUNTIME, "Unable to create Makefile")
	}
	makefile.WriteString(GenerateMakefile(*specifications, forWindows))
}

// Setup the project.
func InitFolder(specifications *CmdArgs) {
	var path string // empty by default

	if specifications.operation == NEW {
		os.Mkdir(specifications.name, os.ModePerm)
		path = specifications.name + "/"
	}

	os.Mkdir(path+"build/", os.ModePerm)
	os.Mkdir(path+"include/", os.ModePerm)
	os.Mkdir(path+"src/", os.ModePerm)

	readme, err := os.Create(path + "README.txt")
	if err != nil {
		Die(RUNTIME, "Unable to create README.txt")
	}
	readme.WriteString(GenerateReadMeString(specifications.name))

	// won't react if the hello world sample fails to be created
	helloworld, err := os.Create(path + "src/main." + LanguageAsString(specifications.language))
	if err != nil {
		Die(RUNTIME, fmt.Sprintf("Unable to create main.%s",
			LanguageAsString(specifications.language)))
	}
	helloworld.WriteString(GenerateHelloWorldString(specifications.language))

	InitBuildFolder(path+"build/linux/", specifications, false)
	if specifications.crossCompile {
		InitBuildFolder(path+"build/windows/", specifications, true)
	}
}
