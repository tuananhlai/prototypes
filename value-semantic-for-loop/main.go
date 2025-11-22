package main

import "fmt"

// Sample output:
// == demoSlice ==
// friend[1] = Bob
// friend = Dave
// == demoArray ==
// friend[1] = Bob
// friend = Bob
func main() {
	fmt.Println("== demoSlice ==")
	demoSlice()
	fmt.Println("== demoArray ==")
	demoArray()
}

// Iterating through a slice of strings, updating the value of
// a particular element and checking its value afterward.
func demoSlice() {
	friends := []string{"Alice", "Bob", "Charlie"}
	fmt.Printf("friend[1] = %s\n", friends[1])

	for i, friend := range friends {
		friends[1] = "Dave"

		if i == 1 {
			fmt.Printf("friend = %s\n", friend)
		}
	}
}

// Iterating through an array of strings, updating the value of
// a particular element and checking its value afterward.
func demoArray() {
	friends := [3]string{"Alice", "Bob", "Charlie"}
	fmt.Printf("friend[1] = %s\n", friends[1])

	// The for loop iterates over a copy of friends instead
	// of the actual `friends` variable.
	for i, friend := range friends {
		friends[1] = "Dave"

		if i == 1 {
			fmt.Printf("friend = %s\n", friend)
		}
	}
}
