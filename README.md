# Fixedwidth
[![Build Status](https://travis-ci.org/huydang284/fixedwidth.svg?branch=master)](https://travis-ci.org/huydang284/fixedwidth)
[![Report](https://goreportcard.com/badge/github.com/huydang284/fixedwidth)](https://goreportcard.com/badge/github.com/huydang284/fixedwidth)

Fixedwidth is a Go package that provides a simple way to define fixed-width data, fast encoding and decoding also is the project's target.

## Character encoding supported
UTF-8

## Getting Started
### Installation
To start using Fixedwidth, run `go get`:
``` 
$ go get github.com/huydang284/fixedwidth
```

### How we limit a struct field
To limit a struct field, we use `fixed` tag.

Example:
```go
type people struct {
    Name string `fixed:"10"`
    Age  int    `fixed:"3"`
}
```

If the value of struct field is longer than the limit that we defined, redundant characters will be truncated.

Otherwise, if the value of struct field is less than the limit, additional spaces will be appended.

### Encoding
We can use `Marshal` function directly to encode fixed-width data.

```go
package main

import (
    "fmt"
    "github.com/huydang284/fixedwidth"
)

type people struct {
    Name string `fixed:"10"`
    Age  int    `fixed:"3"`
}

func main() {
    me := people {
        Name: "Huy",
        Age: 25,
    }
    data, _ := fixedwidth.Marshal(me)
    fmt.Println(string(data))
}
```

The result will be:
```
Huy       25 
```

### Decoding
For decoding, we use `Unmarshal`.

```go
package main

import (
    "fmt"
    "github.com/huydang284/fixedwidth"
)

type people struct {
    Name string `fixed:"10"`
    Age  int    `fixed:"3"`
}

func main() {
    var me people
    data := []byte("Huy       25 ")
    fixedwidth.Unmarshal(data, &me)
    fmt.Printf("%+v", me)
}
```

The result will be:
```
{Name:Huy Age:25}
```

## Author
Huy Dang ([huydangg28@gmail.com](mailto:huydangg28@gmail.com))

## License
Fixedwidth source code is available under the [MIT License](https://github.com/huydang284/fixedwidth/blob/master/LICENSE).