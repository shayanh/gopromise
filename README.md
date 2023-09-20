# Go Promise

Go Promise is a simple library implementing promises in Go. 
Inspired by https://github.com/chebyrash/promise/blob/master/promise.go

```go
package main

import (
	"fmt"

	"github.com/shayanh/gopromise"
)

func main() {
	p := gopromise.New(func(resolve gopromise.ResolveFn[int], reject gopromise.RejectFn) {
		resolve(2)
	}).Then(func(d int, err error) (int, error) {
		return d * 2, err
	}).Then(func(d int, err error) (int, error) {
		return d * 3, err
	})

	d, err := p.Await()
	fmt.Println(d, err)
}
```