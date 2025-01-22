# ericool

[![GoDoc](https://godoc.org/github.com/ericxiao417/ericool?status.svg)](https://godoc.org/github.com/ericxiao417/ericool)

Eric's cool utilities for Go development!

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Technologies](#technologies)
- [Installation](#installation)
- [Usage](#usage)
- [Modules](#modules)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Introduction
`ericool` is a collection of utilities for Go development, created by Eric Xiao. It includes various packages that provide functionalities such as assertions for unit testing, consistent hashing, logging, task pooling, and file batch processing.

## Features
- Assertion utilities for unit testing
- Consistent hashing implementation
- Logging utilities with different log levels and output options
- Task pooling for managing concurrent tasks
- File batch processing utilities

## Technologies
- Go (Golang)

## Installation
To install the `ericool` package, you can use `go get`:

```bash
go get github.com/ericxiao417/ericool
```

## Usage
Here are some examples of how to use the different modules in the `ericool` package.

### Assertions
```go
import "github.com/ericxiao417/ericool/assert"

func TestEqual(t *testing.T) {
    assert.Equal(t, expected, actual)
}
```

### Consistent Hashing
```go
import "github.com/ericxiao417/ericool/consistenthash"

ch := consistenthash.New(1000)
ch.Add("node1", "node2", "node3")
node, err := ch.Get("key")
```

### Logging
```go
import "github.com/ericxiao417/ericool/ericlog"

ericlog.Infof("This is an info message")
```

### Task Pool
```go
import "github.com/ericxiao417/ericool/taskpool"

pool, _ := taskpool.NewPool()
pool.Go(func(param ...interface{}) {
    // Task code here
})
```

### File Batch Processing
```go
import "github.com/ericxiao417/ericool/filebatch"

err := filebatch.Walk("/path/to/dir", true, ".txt", func(path string, info os.FileInfo, content []byte, err error) []byte {
    // Process file content here
    return content
})
```

## Modules
### assert
Provides assertion utilities for unit testing.

### bele
Provides functions for big-endian and little-endian binary encoding and decoding.

### consistenthash
Implements consistent hashing.

### ericlog
Provides logging utilities with different log levels and output options.

### fake
Provides utilities for testing, such as fake time and fake OS exit.

### filebatch
Provides utilities for batch processing files.

### taskpool
Provides a task pool for managing concurrent tasks.

## Contributing
Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository
2. Create a new branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -m 'Add some feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Open a pull request

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact
- Name: Eric Xiao
- GitHub: [ericxiao417](https://github.com/ericxiao417)

