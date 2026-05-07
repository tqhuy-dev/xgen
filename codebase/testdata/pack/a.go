package pack

func helper() {}

// Exported calls helper in the same file.
func Exported() {
	helper()
}
