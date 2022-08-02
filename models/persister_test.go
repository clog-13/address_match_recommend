package models

import (
	"fmt"
	"testing"
)

func TestRootRegion(t *testing.T) {
	root := NewAddressPersister().RootRegion()
	fmt.Println(root.Types)
	fmt.Println(root.Children[0].Types)
}
