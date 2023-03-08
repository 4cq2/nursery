package main

import (
   "fmt"
   "time"
)

// Username may only contain alphanumeric characters or single hyphens, and
// cannot begin or end with a hyphen.
const digits = "023456789abcdefghijkmnopqrstuvwxyz"

func format_bits(u uint64) string {
   var a [13]byte // 8qtr74ui5erii
   i := len(a)
   b := uint64(len(digits))
   for u >= b {
      i--
      q := u / b
      a[i] = digits[u-q*b]
      u = q
   }
   i--
   a[i] = digits[u]
   return string(a[i:])
}

func format(f float64) string {
   return format_bits(uint64(f))
}

func main() {
   now := time.Now()
   then := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
   for {
      dur := now.Sub(then)
      id := format(dur.Hours())
      if len(id) >= 4 {
         break
      }
      fmt.Println(id)
      id = format(dur.Minutes())
      fmt.Println(id)
      id = format(dur.Seconds())
      fmt.Print(id, "\n\n")
      then = then.AddDate(-1, 0, 0)
   }
}

