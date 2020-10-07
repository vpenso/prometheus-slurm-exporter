/* Copyright 2020 Victor Penso

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
        "io/ioutil"
        "os/exec"
        "log"
        "strings"
)

func AccountsData() []byte {
        cmd := exec.Command("squeue", "-h", "-o '%A|%a|%T'")
        stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	out, _ := ioutil.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return out
}

type AccountMetrics struct {
        resv float64
}

func ParseAccountsMetrics(input []byte) map[string]map[string]int {
        accounts := make(map[string]map[string]int)
        lines := strings.Split(string(input), "\n")
        for _, line := range lines {
                if strings.Contains(line,"|") {
                        log.Print(line)

                        account := strings.Split(line,"|")[1]
                        _,key := accounts[account]
                        if !key {
                                accounts[account] = make(map[string]int)
                        }

                        state := strings.Split(line,"|")[2]
                        state = strings.ToLower(state)
                        _,key = accounts[account][state]
                        if !key {
                                accounts[account][state] = 1
                        } else {
                                accounts[account][state] += 1
                        }
                }
        }
        return accounts
}

