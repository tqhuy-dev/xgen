package main

import (
	"fmt"

	"github.com/tqhuy-dev/xgen/dsa"
)

func main() {
	type Data struct {
		Name string
		Age  int
	}
	data := dsa.HeapList[Data]{
		{Data: Data{Name: "tqhuy", Age: 20}, Point: 20},
		{Data: Data{Name: "tqhuy", Age: 50}, Point: 50},
		{Data: Data{Name: "tqhuy", Age: 60}, Point: 60},
		{Data: Data{Name: "tqhuy", Age: 10}, Point: 10},
		{Data: Data{Name: "tqhuy", Age: 30}, Point: 30},
	}
	fmt.Println(dsa.TopMaxPoint(data, 3))
}
