all: no-acceld 
.PHONY: clean install

no-acceld:
	go build

install:
	cp no-acceld /usr/bin/no-acceld

clean:
	test -f no-acceld && rm no-acceld 
