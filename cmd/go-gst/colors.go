package main

import (
	"fmt"
	"io"
	"strings"
)

type color string

var (
	colorReset       color = "\033[0m"
	colorBlack       color = "\033[0;30m"
	colorRed         color = "\033[0;31m"
	colorGreen       color = "\033[0;32m"
	colorOrange      color = "\033[0;33m"
	colorBlue        color = "\033[0;34m"
	colorPurple      color = "\033[0;35m"
	colorCyan        color = "\033[0;36m"
	colorLightGray   color = "\033[0;37m"
	colorDarkGray    color = "\033[1;30m"
	colorLightRed    color = "\033[1;31m"
	colorLightGreen  color = "\033[1;32m"
	colorYellow      color = "\033[1;33m"
	colorLightBlue   color = "\033[1;34m"
	colorLightPurple color = "\033[1;35m"
	colorLightCyan   color = "\033[1;36m"
	colorWhite       color = "\033[1;37m"
)

func disableColor() {
	colorReset = ""
	colorBlack = ""
	colorRed = ""
	colorGreen = ""
	colorOrange = ""
	colorBlue = ""
	colorPurple = ""
	colorCyan = ""
	colorLightGray = ""
	colorDarkGray = ""
	colorLightRed = ""
	colorLightGreen = ""
	colorYellow = ""
	colorLightBlue = ""
	colorLightPurple = ""
	colorLightCyan = ""
	colorWhite = ""
}

func (c color) print(s string)                       { fmt.Printf("%s%s%s", c, s, colorReset) }
func (c color) printIndent(i int, s string)          { c.print(fmt.Sprintf("%s%s", strings.Repeat(" ", i), s)) }
func (c color) printf(f string, args ...interface{}) { c.print(fmt.Sprintf(f, args...)) }
func (c color) printfIndent(i int, f string, args ...interface{}) {
	c.printf(fmt.Sprintf("%s%s", strings.Repeat(" ", i), fmt.Sprintf(f, args...)))
}
func (c color) fprint(w io.Writer, s string) { fmt.Fprintf(w, fmt.Sprintf("%s%s%s", c, s, colorReset)) }
func (c color) fprintf(w io.Writer, f string, args ...interface{}) {
	c.fprint(w, fmt.Sprintf(f, args...))
}
func (c color) fprintIndent(w io.Writer, i int, s string) {
	c.fprint(w, fmt.Sprintf("%s%s", strings.Repeat(" ", i), s))
}
func (c color) fprintfIndent(w io.Writer, i int, f string, args ...interface{}) {
	c.fprint(w, fmt.Sprintf("%s%s", strings.Repeat(" ", i), fmt.Sprintf(f, args...)))
}
