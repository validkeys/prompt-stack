package greeter

// Hello returns a short greeting for the given name.
func Hello(name string) string {
	if name == "" {
		name = "world"
	}
	return "hello, " + name
}
