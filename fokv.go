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
    f.FilePath = path
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
                if v != "" {
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
            k = s[0:si]
            v = s[si+1:]
        }
        f.rows = append(f.rows, KeyValue{key: k, value: v})
        k = ""
    }
    return f
}

func (f *Fokv) Put(k string, v string) {
    f.rows = append(f.rows, KeyValue{key: k, value: v})
}

func (f *Fokv) Save() {
   file, err := os.Create(f.FilePath)
   if err != nil {
        panic(fmt.Sprintf("error write %s: %v", f.FilePath, err))
   }
   defer file.Close()
   w := bufio.NewWriter(file)
    row := ""   
    for _, kv := range f.rows {
        if kv.key == "#" {
            row = kv.key + kv.value + "\n"
        } else if strings.Index(kv.value, "\n") < 0 {
            row = kv.key + " " + kv.value + "\n"
        } else {
            row = kv.key + "\n" + kv.value + "\n#\n"
        }
      n4, err := w.WriteString(row)
      if err != nil {
        panic(fmt.Sprintf("error WriteString %s: %v", row, err))
      }
      fmt.Printf("wrote %d bytes\n", n4)
   }
   w.Flush()
}