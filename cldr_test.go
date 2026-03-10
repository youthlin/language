package language

import (
	"testing"
)

func TestGetParent(t *testing.T) {
	i := _zh_Hans_CN
	p := i.Parent()
	t.Logf("%v -> %v", i, p)
}
