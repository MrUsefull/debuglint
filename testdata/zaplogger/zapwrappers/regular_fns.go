package zapwrappers

import "fmt"

func JustAnEveryDayNormalFunc() {
	fmt.Println("I'm just a regular every day normal function")
}

func JustAnEveryDayNormalFuncWithLogging() {
	fmt.Println("I'm just a regular every day normal function")
	SimpleSugaredWrapperFunc("I log things for breakfast")
	JustAnEveryDayNormalFunc()
}

func DoAddition(a, b int) int {
	myVal := 1 + 2
	return myVal
}
