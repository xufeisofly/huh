package huh

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func multiArgsToMap(args ...interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for i := 0; i < len(args); i += 2 {
		result[args[i].(string)] = args[i+1]
	}
	return result
}

func uint8ToString(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}
