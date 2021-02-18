package tools

// HandleErr panics the error if exists
func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}
