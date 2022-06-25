/*

We check whether the passed argument is a struct. Then we get the name of the struct from its reflect.Type using the Name() method. Then we use t and start creating the query.

The case statement checks whether the current field is reflect.Int, if that's the case we extract the value of that field as int64 using the Int() method. The if else statement is used to handle edge cases. Please add logs to understand why it is needed. Similar logic is used to extract the string.

We also add checks to prevent the program from crashing when unsupported types are passed to the createQuery function.

How do you know when to use reflection? Rob Pike says: "Clear is better than clever. Reflection is never clear."

Reflection is a very powerful and advanced concept in Go and it should be used with care. It is very difficult to write clear and maintainable code using reflection. It should be avoided wherever possible and should be used only when absolutely necessary.


*/

package main

import (
	"fmt"
	"reflect"
)

type order struct {
	ordId      int
	customerId int
}

type employee struct {
	name    string
	id      int
	address string
	salary  int
	country string
}

func createQuery(q interface{}) {
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		t := reflect.TypeOf(q).Name()
		query := fmt.Sprintf("insert into %s values(", t)
		v := reflect.ValueOf(q)
		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Kind() {
			case reflect.Int:
				if i == 0 {
					query = fmt.Sprintf("%s%d", query, v.Field(i).Int())
				} else {
					query = fmt.Sprintf("%s, %d", query, v.Field(i).Int())
				}
			case reflect.String:
				if i == 0 {
					query = fmt.Sprintf("%s\"%s\"", query, v.Field(i).String())
				} else {
					query = fmt.Sprintf("%s, \"%s\"", query, v.Field(i).String())
				}
			default:
				fmt.Println("Unsupported type")
				return
			}
		}
		query = fmt.Sprintf("%s)", query)
		fmt.Println(query)
		return

	}
	fmt.Println("unsupported type")
}

func main() {
	o := order{
		ordId:      456,
		customerId: 56,
	}
	createQuery(o)

	e := employee{
		name:    "Naveen",
		id:      565,
		address: "Coimbatore",
		salary:  90000,
		country: "India",
	}
	createQuery(e)
	i := 90
	createQuery(i)

}
