# errdata

[![GoDoc](https://godoc.org/github.com/zeebo/errs/errdata?status.svg)](https://godoc.org/github.com/zeebo/errs/errdata)
[![Sourcegraph](https://sourcegraph.com/github.com/zeebo/errs/-/badge.svg)](https://sourcegraph.com/github.com/zeebo/errs?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeebo/errs/errdata)](https://goreportcard.com/report/github.com/zeebo/errs/errdata)

errdata helps with associating some data to error classes.

### Adding data

The [Set][Set] function associates some data with some error class and key. For
example:

```go
var (
	Unauthorized = errs.Class("unauthorized")
	NotFound     = errs.Class("not found")
)

type httpErrorCodeKey struct{}

func init() {
	errdata.Set(Unauthorized, httpErrorCodeKey{}, http.StatusUnauthorized)
	errdata.Set(NotFound, httpErrorCodeKey{}, http.StatusNotFound)
}
```

Why do that? [Get][Get] can read the associated data for an error if it was
wrapped by any of the Classes you have set data on. For example:

```go
func getStatusCode(err error) int {
	code, _ := errdata.Get(err, httpErrorCodeKey{}).(int)
	if code == 0 {
		code = http.StatusInternalServerError
	}
	return code
}
```

If the error has been wrapped by multiple classes for that key, the value for
the most recently wrapped class is returned. For example:

```go
func whatStatusCodeCode() {
	err := NotFound.Wrap(Unauthorized.New("test"))
	fmt.Println(getStatusCode(err))

	// output:
	// 404
}
```

### Contributing

errdata is released under an MIT License. If you want to contribute, be sure to
add yourself to the list in AUTHORS.

[Set]: https://godoc.org/github.com/zeebo/errs/errdata#Set
[Get]: https://godoc.org/github.com/zeebo/errs/errdata#Get
