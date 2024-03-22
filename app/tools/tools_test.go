package tools

import (
	"fmt"
	"testing"
)

func TestGetUID(t *testing.T) {
	id := GetUID()
	fmt.Printf("id:%d\n", id)
	id1 := GetUID()
	fmt.Printf("id:%d\n", id1)
	id2 := GetUID()
	fmt.Printf("id:%d\n", id2)
	id3 := GetUID()
	fmt.Printf("id:%d\n", id3)
	id4 := GetUID()
	fmt.Printf("id:%d\n", id4)
}
