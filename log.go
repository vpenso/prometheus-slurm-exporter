/* Copyright 2021 Rovanion Luckey

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>. */

package main

import (
	"fmt"
	"log"
	"runtime"
)

// Log messages to stdout preceeded by the current timestamp and the
// module.Function where Log was called from.
func Log(messages ...interface{}) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		log.Print(details.Name(), ": ", fmt.Sprint(messages...))
	} else {
		log.Print(fmt.Sprint(messages...))
	}
}

// Log messages to stdout preceeded by the current timestamp and the
// module.Function where Log was called from and then terminate the
// program.
func Fatal(messages ...interface{}) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		log.Fatal(details.Name(), ": ", fmt.Sprint(messages...))
	} else {
		log.Fatal(fmt.Sprint(messages...))
	}
}

// Log messages to stdout preceeded by the current timestamp and the
// module.Function of the parent function from where Log was called
// from. This Log function is meant to be called from utility
// functions.
func Log2(messages ...interface{}) {
	pc, _, _, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		log.Print(details.Name(), ": ", fmt.Sprint(messages...))
	} else {
		log.Print(fmt.Sprint(messages...))
	}
}

// Log messages to stdout preceeded by the current timestamp and the
// module.function of the parent function from where Log was called
// from and then terminate the program. This Log function is meant to
// be called from utility functions.
func Fatal2(messages ...interface{}) {
	pc, _, _, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		log.Fatal(details.Name(), ": ", fmt.Sprint(messages...))
	} else {
		log.Fatal(fmt.Sprint(messages...))
	}
}
