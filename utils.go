package huh

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func multiArgsToMap(args ...interface{}) map[string]interface{} {
	var result map[string]interface{}
	for i := 0; i < len(args); i += 2 {
		result[args[i].(string)] = args[i+1]
	}
	return result
}
