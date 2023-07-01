package libgodone

import "fmt"

type Compiler int

const (
	GCC   Compiler = iota
	CLANG Compiler = iota
)

func CompilerAsString(compiler Compiler) string {
	switch compiler {
	case GCC:
		return "gcc"

	case CLANG:
		return "clang"

	default:
		return "gcc"
	}
}

func StringAsCompiler(compiler string) (Compiler, error) {
	switch compiler {
	case "gcc", "cc", "CC", "GCC":
		return GCC, nil

	case "clang", "Clang", "CLang", "CLANG":
		return CLANG, nil

	default:
		return GCC, fmt.Errorf("unknown compiler")
	}
}

func GetCrossCompilerVariant(language Language) string {
	switch language {
	case C:
		return "x86_64-w64-mingw32-gcc"

	default:
		return "x86_64-w64-mingw32-g++"
	}
}

func GetCompilerVariant(language Language, compiler Compiler) string {
	switch language {
	case C:
		if compiler == GCC {
			return "gcc"
		} else {
			return "clang"
		}

	default:
		if compiler == GCC {
			return "g++"
		} else {
			return "clang++"
		}
	}
}
