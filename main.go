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
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"
)

type Config struct {
	Device          string            `json:"device"`
	Properties      map[string]string `json:"properties"`
	IntervalSeconds int               `json:"interval"`
	EnableLog       bool              `json:"log"`
}

type LoopFlag struct {
	mu    sync.Mutex
	value bool
}

// This is probably a slow way of doing this but for will work for now.
func stringSelectUntilSpace(value string) string {
	var sb strings.Builder

	for _, r := range value {
		if unicode.IsSpace(r) {
			break
		}
		sb.WriteRune(r)
	}

	return sb.String()
}

func main() {
	cfg := Config{
		Device:          "",
		Properties:      make(map[string]string),
		IntervalSeconds: 0,
		EnableLog:       false,
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
		if cfg.EnableLog {
			log.Printf("starting with the following configuration:\n\tDevice name:\t%v\n\tInterval:\t%v seconds.\n\n", cfg.Device, cfg.IntervalSeconds)
		}
		for {
			loopFlag.mu.Lock()
			if !loopFlag.value {
				loopFlag.mu.Unlock()
				break
			}
			loopFlag.mu.Unlock()

			devListData, err := exec.Command("xinput", "--list").Output()
			if err != nil {
				panic(err)
			}

			if len(devListData) == 0 {
				panic(errors.New("empty command output"))
			}

			for _, dev := range strings.Split(string(devListData), "\n") {
				if strings.Contains(dev, "id=") &&
					strings.Contains(dev, "pointer") &&
					strings.Contains(dev, cfg.Device) {
					dataAlpha := strings.SplitN(dev, "id=", 2)
					if len(dataAlpha) != 2 {
						continue
					}
					devIdStr := stringSelectUntilSpace(dataAlpha[1])
					devIdInt, err := strconv.Atoi(devIdStr)
					if err != nil {
						continue
					}
					if devIdInt < 0 {
						continue
					}
					if cfg.EnableLog {
						log.Printf("match device with id %v\n", devIdStr)
					}
					for propName, propValue := range cfg.Properties {
						if len(propName) == 0 || len(propValue) == 0 {
							continue
						}
						out, err := exec.Command("xinput", "--set-prop", devIdStr, "libinput "+propName, propValue).CombinedOutput()
						if cfg.EnableLog {
							if err != nil {
								log.Printf("xinput error: %v", err)
							}
							if len(out) > 0 {
								log.Printf("xinput output: %v\n", string(out))
							}
						}
					}
				}
			}

			time.Sleep(time.Duration(cfg.IntervalSeconds) * time.Second)
		}
	}()
	sig := <-cancelChan
	if cfg.EnableLog {
		log.Printf("caught signal %v\n", sig)
	}
	loopFlag.mu.Lock()
	loopFlag.value = false
	loopFlag.mu.Unlock()
}
