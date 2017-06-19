SHELL=bash
export SQL_DRIVER=pgsql

.PHONY: .prepare test db init
# db init uuid
all: .prepare test db init uuid

# ignore vendor and vault directories
deps:=$(shell find . \( ! -wholename '*/vendor/*' -a ! -wholename '*/vault/*' -a -iname '*.go' \))

.prepare: $(deps) 
	govendor init
	govendor sync -v

db: db/db

db/db: $(deps) Makefile
	govendor build -o $@ db/main.go

init: init/init

init/init: $(deps) Makefile
	govendor build -o $@ init/main.go


test: db init
	. db/environment; init/init
	. db/environment; db/db

clean:
	rm -f init/init db/db
