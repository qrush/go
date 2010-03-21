package ls

type (
	// Scanner function, must return next token or false.
	// Advance past next token if argument is true.
	Scanner func(bool) (string, bool)

	// Parser function
	Parser func(Scanner) Eval

	// Interpreter function; returns value or os.Error.
	Eval func() interface{}
)
