package huh

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
