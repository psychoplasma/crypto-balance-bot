package concurrency

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"
)

// Retrial represents retrial parameters for recursion
type Retrial struct {
	Limit int           // Max. number of retrial. If it's set to -1, it will run indefinitely
	Delay time.Duration // Delay between each recursion
}

// Retry calls the given function recursively with the given retrial paramters
func Retry(r Retrial, f func() (interface{}, error)) (interface{}, error) {
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()

	for i := r.Limit; ; {
		ret, err := f()
		if err != nil {
			log.Printf("failed(%d) to call function \"%s()\". reason: %s\n", r.Limit-i, funcName, err)

			i--
			if i >= 0 || r.Limit == -1 {
				log.Printf("retrying(%d) to call function \"%s()\"\n", r.Limit-i, funcName)
				time.Sleep(r.Delay)
				continue
			}
			break
		}

		return ret, nil
	}

	return nil, fmt.Errorf("retry has reached to limit: %d", r.Limit)
}
