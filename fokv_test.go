// Package fokv read and write array of pair to file
package fokv

import (
    "fmt"
    "testing"
)

func TestOpen(t *testing.T) {
    f := Open("fokv_test.txt")
    fmt.Printf("%+v\n", f)
    f.Put("key5","value5")
    f.Put("key6","value\nmulti\nline")
    f.Save()
}