package syncer

import (
	"io"
)

// Clone interface.
type Cloner interface {
	Clone() io.Writer
}
