package modulemath

import (
	"fmt"
	"strconv"
	"strings"
)

type ModMath struct {
	N int32 `json:"N"`
}

func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func (mm *ModMath) SetN(n int32) {
	if n < 1 {
		fmt.Printf("ERROR by " + string(n))
	}
	mm.N = n
}

func (mm *ModMath) IntToMod(input int32) string {
	if input == 0 {
		return "0"
	}
	var sb strings.Builder
	for ; input > 0; input /= mm.N {
		sb.WriteString(strconv.Itoa(int(input % mm.N)))
	}

	return reverse(sb.String())
}

func (mm *ModMath) NextOne(input string) string {
	array := []rune(input)
	overflow := true
	outArray := make([]rune, 0)
	i := len(array) - 1
	for ; i >= 0 && overflow; i-- {
		digit := int(array[i] - '0')
		digit++
		if digit == int(mm.N) {
			digit = 0
		} else {
			overflow = false
		}
		outArray = append([]rune{rune('0' + digit)}, outArray...)
	}
	if overflow {
		if outArray[0] == '0' {
			outArray = append([]rune{'0'}, outArray...)
			outArray[len(outArray)-1] = '1'
		}
	} else if i >= 0 {
		outArray = append(array[:i+1], outArray...)
	}
	if outArray[len(outArray)-1] == '0' {
		outArray[len(outArray)-1] = '1'
	}
	return string(outArray)
}

func (mm *ModMath) ModToInt(input string) int {

	output := 0
	for _, ar := range input {
		output = output*int(mm.N) + int(ar-'0')
	}

	return output
}

func (mm *ModMath) CompareTwoNumbs(num1, num2 string) int {
	if len(num1) > len(num2) {
		return 1
	}
	if len(num1) < len(num2) {
		return -1
	}

	return mm.ModToInt(num1) - mm.ModToInt(num2)
}

func EditDistance(str1, str2 string) int {
	arr1 := []rune(str1)
	arr2 := []rune(str2)
	lenDiff := len(arr1) - len(arr2)
	// fmt.Printf("%v vs. %v := %d\n", str1, str2, lenDiff)
	if lenDiff != 0 {
		for i := 0; i < (-lenDiff); i++ {
			arr1 = append(arr1, '0')
		}
		for i := 0; i < (lenDiff); i++ {
			arr2 = append(arr2, '0')
		}
	}

	// fmt.Printf("%s vs. %s := %d\n", string(arr1), string(arr2), lenDiff)

	dist := 0
	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			dist++
		}
	}

	return dist
}
