package checker_test

/*! function has more than 5 results, consider to simplify the function */
func bad1() (int, int, int, int, int, int) {
	return 0, 0, 0, 0, 0, 0
}

/*! function has more than 5 results, consider to simplify the function */
func bad2() (_, _, _, _, _, _ int) {
	return 0, 0, 0, 0, 0, 0
}

type badStruct struct{}

/*! function has more than 5 results, consider to simplify the function */
func (badStruct) bad() (_, _, _, _, _, _ int) {
	return 0, 0, 0, 0, 0, 0
}
