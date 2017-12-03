// Package fokv read and write array of pair to file
package fokv

import (
    "os"
    "fmt"
    "bufio"
    "strings"
)

// Options define a set of properties that dictate Diskv behavior.
// All values are optional.
type Options struct {
  FilePath     string
}

type KeyValue struct {
  key   string
  value string
}

type Fokv struct {
  Options
  rows []KeyValue
}

// Open file and fill Fokv
func Open(path string) *Fokv {
// Open file and create scanner on top of it
    file, err := os.Open(path)
    if err != nil {
        panic(fmt.Sprintf("error opening %s: %v", path, err))
    }
    f := &Fokv{}
    //f.FilePath := path
    scanner := bufio.NewScanner(file)
    k := ""
    v := ""

    // Scan for next token. 
    for i := 0; scanner.Scan(); i++ {
         // False on error or EOF. Check error
        if err := scanner.Err(); err != nil {
            fmt.Fprintln(os.Stderr, "error reading from file:", err)
	    os.Exit(3)
        }
        s := scanner.Text()
        if k != "" {
            if  s != "#" {
                if v!="" {
                    v += "\n"
                }
              v += s
              continue
            }
        } else if s[0:1] == "#" {
            k = "#"
            v = s[1:]
        } else if len(s) == 0 {
            continue
        } else {
            si := strings.Index(s," ")
            if si < 0 {
              k = s
              v = ""
              continue  
            }
            k = s[0:si-1]
            v = s[si+1:]
        }
        f.rows = append(f.rows, KeyValue{key: k, value: v})
        k = ""
    }
    return f
}