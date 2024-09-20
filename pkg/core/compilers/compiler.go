package compilers

type Compiler interface {
	Compile(string) (Program, error)
}
