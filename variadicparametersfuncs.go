//variadic functions, last arg presumed to be last name

// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	var line string

	line = fullname("Robert", "Michael", "Alexander", "Loblaw")
	fmt.Println(line)

	line = fullname("Robert", "Loblaw")
	fmt.Println(line)

	line = fullname("Robert")
	fmt.Println(line)
}

func fullname(values ...string) string {
	var line string
	var lastname string
	for i, v := range values {

		if i == len(values)-1 && len(values) > 1 {
			lastname = v + ", "
		} else {
			line = line + v + " "
		}
	}
	return lastname + line
}

/*
//required first name

// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	var line string

	line = fullname("Robert", "Michael", "Alexander", "Loblaw")
	fmt.Println(line)

	line = fullname("Robert", "Loblaw")
	fmt.Println(line)

	line = fullname("Robert")
	fmt.Println(line)

}

func fullname(fname string, values ...string) string {
	var line string
	var lastname string
	for i, v := range values {

		if i == len(values)-1 && len(values) > 1 {
			lastname = v + ", "
		} else {
			line = line + v + " "
		}
	}
	return lastname + fname + " " + line
}

*/

//can only use ... with final parameter in list

/*
//if middle part has one, assume it's last name and assign NMN

package main

import "fmt"

func main() {
	var line string

	line = fullname("Robert", "Michael", "Alexander", "Loblaw")
	fmt.Println(line)

	line = fullname("Robert", "Loblaw")
	fmt.Println(line)

	line = fullname("Robert")
	fmt.Println(line)

}

func fullname(fname string, values ...string) string {
	var line string
	//var lastname string
	for _, v := range values {

		if len(values) == 1 {
			line = "NMN " + v
		} else {
			line = line + v + " "
		}

		//if i == len(values)-1 && len(values) > 1 {
		//		lastname = v + ", "
		//	} else {
		//line = line + v + " "
		//}

	}
	return fname + " " + line
}

*/

/*
//accepting delimiter

package main

import "fmt"

func main() {
	var line string

	line = fullname(" ", "Robert", "Michael", "Alexander", "Loblaw")
	fmt.Println(line)

	line = fullname(" ", "Robert", "Loblaw")
	fmt.Println(line)

	line = fullname(" ", "Robert")
	fmt.Println(line)

}

func fullname(delimiter string, fname string, values ...string) string {
	var line string
	//var lastname string
	for _, v := range values {

		if len(values) == 1 {
			line = "NMN" + delimiter + v
		} else {
			line = line + v + delimiter
		}

		//if i == len(values)-1 && len(values) > 1 {
		//		lastname = v + "," + delimiter
		//	} else {
		//line = line + v + delimiter
		//}

	}
	return fname + " " + line
}

*/
/*
//exploded slice

//required first name

// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	var name string
	names := []string{"Alexander", "Loblaw"}

	name = fullname(" ", "Robert", names...)
	fmt.Println(name)

}

func fullname(delimiter string, fname string, values ...string) string {
	var name string
	var lastname string
	for i, v := range values {

		if i == len(values)-1 && len(values) > 1 {
			lastname = v + "," + delimiter
		} else {
			name = name + v + delimiter
		}
	}
	return lastname + fname + delimiter + name
}

*/
