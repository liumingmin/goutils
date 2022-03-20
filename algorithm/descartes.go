package algorithm

func DescartesCombine(dimValue [][]string) [][]string {
	if len(dimValue) == 0 {
		return [][]string{}
	}

	result := make([][]string, 0)
	curList := make([]string, 0, len(dimValue))
	backtrace(dimValue, 0, &result, &curList)
	return result
}

func backtrace(dimValue [][]string, index int, result *[][]string, curList *[]string) {
	if len(*curList) == len(dimValue) {
		newList := make([]string, len(dimValue))
		copy(newList, *curList)
		*result = append(*result, newList)
		return
	}

	subDimValue := dimValue[index]
	for i := 0; i < len(subDimValue); i++ {
		*curList = append(*curList, subDimValue[i])

		backtrace(dimValue, index+1, result, curList)

		*curList = (*curList)[:len(*curList)-1]
	}
}
