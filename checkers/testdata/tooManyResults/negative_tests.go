package checker_test

func good1() (int, int) {
	return 0, 0
}

func good2() (string, int) {
	return "", 0
}

func good3() int {
	return 0
}

type goodStruct struct{}

func (goodStruct) good() (_, _, _, _, _ int) {
	return 0, 0, 0, 0, 0
}
