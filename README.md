# no-acceld
A program that constantly disables mouse acceleration  (For libinput)

# Prerequisites
 - Go programming language - [Get it here](https://golang.org/doc/install)
 - A distro/operative system with libinput (xinput command should be available)

# How to install
**Optimization note**: Since Go 1.18, if you are using an x86_64 CPU you can setup the environment variable `GOAMD64`

You can found a detailed description [here](https://github.com/golang/go/wiki/MinimumRequirements#amd64)

***TL;DR:** Use v1 for any x64 cpu, v2 for circa 2009 "Nehalem and Jaguar", v3 for circa 2015 "Haswell and Excavator" and v4 for AVX-512.*

**Example:** `export GOAMD64=v2`

	git clone https://git.lucdev.net/luc/no-acceld.git
	cd no-acceld
	make
	sudo make install
	cp config.json.example ~/.config/no-acceld.json

Find the name of your mouse using this command:

    xinput --list

Edit the config file to match the device name:

    $EDITOR ~/.config/no-acceld.json

You can override the location of the config file by using this environment variable:

    CONFIG_FILE=/home/user/example/no-acceld.json no-acceld

# Disclaimer

This program is a work in progress and it's not done.

I made this program for myself, I don't have the intention (for now) to document it and/or give support, so pretty much you are on your own.

This piece of software is barely tested, don't use it on production, or use it at your own risk.