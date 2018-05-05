# errs

[![GoDoc](https://godoc.org/github.com/zeebo/errs?status.svg)](https://godoc.org/github.com/zeebo/errs)
[![Sourcegraph](https://sourcegraph.com/github.com/zeebo/errs/-/badge.svg)](https://sourcegraph.com/github.com/zeebo/errs?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeebo/errs)](https://goreportcard.com/report/github.com/zeebo/errs)

errs is a package for making errors friendly and easy.

### Creating Errors

The easiest way to use it, is to use the package level [New][New] function.
It's much like `fmt.Errorf`, but better. For example:

```go
func checkThing() error {
	return errs.New("what's up with %q?", "zeebo")
}
```

Why is it better? Errors come with a stack trace that is only printed
when a `"+"` character is used in the format string. This should retain the
benefits of being able to diagnose where and why errors happen, without all of
the noise of printing a stack trace in every situation. For example:

```go
func doSomeRealWork() {
	err := checkThing()
	if err != nil {
		fmt.Printf("%+v\n", err) // contains stack trace if it's a errs error.
		fmt.Printf("%v\n", err)  // does not contain a stack trace
		return
	}
}
```

### Error Classes

You can create a [Class][Class] of errors and check if any error was created by
that [Class][Class]. The [Class][Class] name is prefixed to all of the errors
it creates. For example:

```go
var Unauthorized = errs.Class("unauthorized")

func checkUser(username, password string) error {
	if username != "zeebo" {
		return Unauthorized.New("who is %q?", username)
	}
	if password != "hunter2" {
		return Unauthorized.New("that's not a good password, jerkmo!")
	}
	return nil
}

func handleRequest() {
	if err := checkUser("zeebo", "hunter3"); Unauthorized.Has(err) {
		fmt.Println(err)
	}

	// output:
	// unauthorized: that's not a good password, jerkmo!
}
```

[Class][Class]es can also [Wrap][Wrap] other errors, and errors may be
[Wrap][Wrap]ped multiple times. For example:

```go
var (
	Error        = errs.Class("mypackage")
	Unauthorized = errs.Class("unauthorized")
)

func deep3() error {
	return fmt.Errorf("ouch")
}

func deep2() error {
	return Unauthorized.Wrap(deep3())
}

func deep1() error {
	return Error.Wrap(deep2())
}

func deep() {
	fmt.Println(deep1())

	// output:
	// mypackage: unauthorized: ouch
}
```

In the above example, both `Error.Has(deep1())` and `Unauthorized.Has(deep1())`
would return `true`, and the stack trace would only be recorded once at the
`deep2` call.

In addition, when an error has been [Wrap][Wrap]ped, [Wrap][Wrap]ping it again
with the same [Class][Class] will not do anything. For example:

```go
func doubleWrap() {
	fmt.Println(Error.Wrap(error.New("foo")))

	// output:
	// mypackage: foo
}
```

This is to make it an easier decision if you should [Wrap][Wrap] or not.

### Utilities

[Classes][Classes] is a helper function to get a slice of [Class][Class]es
that an error has. The earliest wrap is first in the slice. For example:

```go
func getClasses() {
	classes := errs.Classes(deep1())
	fmt.Println(classes[0] == &Unauthorized)
	fmt.Println(classes[1] == &Error)

	// output:
	// true
	// true
}
```

Finally, a helper function, [Unwrap][Unwrap] is provided to get the
[Wrap][Wrap]ped error in cases where you might want to inspect details. For
example:

```go
var Error = Class("mypackage")

func getHandle() (*os.File, error) {
	fh, err := os.Open("neat_things")
	if err != nil {
		return nil, Error.Wrap(err)
	}
	return fh, nil
}

func checkForNeatThings() {
	fh, err := getHandle()
	if os.IsNotExist(errs.Unwrap(err)) {
		panic("no neat things?!")
	}
	if err != nil {
		panic("phew, at least there are neat things, even if i can't see them")
	}
	fh.Close()
}
```

### Contributing

errs is released under an MIT License. If you want to contribute, be sure to
add yourself to the list in AUTHORS.


[New]: https://godoc.org/github.com/zeebo/errs#New
[Class]: https://godoc.org/github.com/zeebo/errs#Class
[Wrap]: https://godoc.org/github.com/zeebo/errs#Class.Wrap
[Unwrap]: https://godoc.org/github.com/zeebo/errs#Unwrap
[Classes]: https://godoc.org/github.com/zeebo/errs#Classes
