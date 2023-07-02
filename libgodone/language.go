package libgodone

import "fmt"

type Language int

const (
	C   Language = iota
	CPP Language = iota
)

// gay as fuckkkkkkkkkkkkkkk

func LanguageAsString(language Language) string {
	switch language {
	case C:
		return "c"

	case CPP:
		return "cpp"

	default:
		return "cpp"
	}
}

func StringAsLanguage(language string) (Language, error) {
	switch language {
	case "c", "C":
		return C, nil

	case "cpp", "c++", "C++", "Cpp", "CPP", "cc", "CC", "cxx", "Cxx", "CXX":
		return CPP, nil

	default:
		return CPP, fmt.Errorf("unknown language")
	}
}
