// Package driver specifies the interface for storage drivers of a Merkle
// weave.
package driver

import "time"

// Interface is the interface a Merkle weave storage driver must satisfy.
type Interface interface {
	Get
	WriteNext(data []byte) (node []byte, n int64, timestamp time.Time, err error)
}
