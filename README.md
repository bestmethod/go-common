# Common Golang functions

### Some common tools I keep using all the time

### import and example
```go
package main

import "github.com/bestmethod/go-common"
import "fmt"

func main() {
    array := []string{"robert","bestmethod","testing","some","content"}
    fmt.Println(gocommon.inArray(array,"bestmethod"))
    fmt.Println(gocommon.inArray(array,"hah"))
}
```
OUTPUT:
```
1
-1
```

##### From exec package, if spawned binary returns an error, extract the error code:
```go
func check_exec_retcode(err error) int
```

##### cut (more like awk -F'split' '{print pos}')

```go
func cut(line string, pos int, split string) string
```

##### Is element in array? If so, return first matching position index. If not, return -1

```go
func inArray(array interface{}, element interface{}) (index int)
```

##### makeError wrapper - because I shouldn't have to errors.New(fmt.Sprintf())

```go
func makeError(format string, args ...interface{}) error
```

##### like makeError, but append err (previous error) as next line

```go
func appendErrorln(err error, format string, args ...interface{}) error
```

Example:

```go
func myCode() {
	some_param := "robert"
    out, err := some_function(some_param)
    if err != nil {
        return appendErrorln(err,"Error in myCode executing some_function(%s):", some_param)
    }
}
```

Will print nicely:
```go
Error in myCode executing some_function(robert):
contents of original err here
```

##### like appendErrorln, but appends in same line

```go
func appendError(err error, format string, args ...interface{}) error
```

Example:

```go
func myCode() {
	some_param := "robert"
    out, err := some_function(some_param)
    if err != nil {
        return appendError(err,"Error in myCode executing some_function(%s): ", some_param)
    }
}
```

Will print nicely:
```go
Error in myCode executing some_function(robert): contents of original err here
```

#### ssh

##### ssh, run cmd and attach stdin and stdout to current stdin and stdout. If cmd is /bin/bash for example, you will be presented with an interactive shell
```go
func RemoteAttachAndRun(user string, addr string, privateKey string, cmd string, stdin *os.File, stdout *os.File, stderr *os.File) error
```

stdin, stdout, stderr are optional. If set, these will be used. If 'nil', standard os.Std* will be used.

##### ssh, run the specified command, and return output as []byte
```go
func RemoteRun(user string, addr string, privateKey string, cmd string) ([]byte, error)
```

##### scp files over to the remote machine
```go
func Scp(user string, addr string, privateKey string, files []fileList) error

type FileList struct {
	SourceFilePath      *string
	SourceFileReader    *bytes.Reader
	DestinationFilePath string
}
```

You can either present a `SourceFilePath` or `SourceFileReader` (if you have contents to be written to a file in a variable, as opposed to a source file). If `SourceFilePath` is set, that will be used (file opened, read and scp over to the remote machine's `DestinationFilePath`). If `SourceFilePath` is nil, `SourceFileReader` will be used instead.