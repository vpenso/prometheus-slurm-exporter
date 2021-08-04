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
	"io/ioutil"
	"os/exec"
	"strings"
)

func Subprocess(executable string, arguments ...string) []byte {
	subprocess := exec.Command(executable, arguments...)
	stdout, err := subprocess.StdoutPipe()
	if err != nil {
		Fatal2("Unable to open stdout for ", executable, ": ", err)
	}
	stderr, err := subprocess.StderrPipe()
	if err != nil {
		Fatal2("Unable to open stderr for ", executable, ": ", err)
	}
	if err := subprocess.Start(); err != nil {
		Fatal2("Failed to start ", executable, ": ", err)
	}
	out, err := ioutil.ReadAll(stdout)
	errOut, _ := ioutil.ReadAll(stderr)
	if err := subprocess.Wait(); err != nil {
		Log2("The command '", executable, " ", strings.Join(arguments, " "), "' failed with the following error:")
		fmt.Print(string(errOut))
		Fatal2("The subprocess ", executable, " terminated with: ", err)
	}
	return out
}
