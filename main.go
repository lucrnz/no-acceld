/*
* no-acceld: a program to keep mouse acceleration disabled
* Copyright (C) 2022 Lucie <lucdev.net>
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Config struct {
	Device          string            `json:"device"`
	Properties      map[string]string `json:"properties"`
	IntervalSeconds int               `json:"interval"`
}

type LoopFlag struct {
	mu    sync.Mutex
	value bool
}

func main() {
	cfg := Config{
		Device:          "",
		Properties:      make(map[string]string),
		IntervalSeconds: 0,
	}

	cfgFilePath := os.Getenv("CONFIG_FILE")

	if len(cfgFilePath) == 0 {
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")

		if len(xdgConfigHome) > 0 {
			cfgFilePath = xdgConfigHome + "/no-acceld.json"
		} else {
			cfgFilePath = os.Getenv("HOME") + "/no-acceld.json"
		}
	}

	cfgFileData, err := os.ReadFile(cfgFilePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(cfgFileData, &cfg)
	if err != nil {
		panic(err)
	}

	if len(cfg.Device) == 0 || len(cfg.Properties) == 0 || cfg.IntervalSeconds <= 0 {
		panic(errors.New("invalid config file, empty properties"))
	}

	cancelChan := make(chan os.Signal, 1)
	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	loopFlag := LoopFlag{
		value: true,
	}

	go func() {
		fmt.Printf("no-acceld: Starting with the following configuration:\n\tDevice name: %v\n\tInterval: %v seconds.\n\n", cfg.Device, cfg.IntervalSeconds)
		for {
			loopFlag.mu.Lock()
			if !loopFlag.value {
				loopFlag.mu.Unlock()
				break
			}
			loopFlag.mu.Unlock()

			listDevices, err := exec.Command("xinput", "--list").Output()
			if err != nil {
				panic(err)
			}

			fmt.Printf("%v\n", string(listDevices))
			time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
		}
		log.Printf("Goodbye!\n")
	}()
	sig := <-cancelChan
	log.Printf("Caught signal %v\n", sig)
	loopFlag.mu.Lock()
	loopFlag.value = false
	loopFlag.mu.Unlock()
}
