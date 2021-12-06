#	======================================================================	#
#	HELPERS
# 	======================================================================	#

## help: print this help message
.PHONY: help
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# Create the new confirm target.
.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]


#	======================================================================	#
#	DEVELOPMENT
# 	======================================================================	#

## run/api: run the cmd/api application	
.PHONY: run/api
run/api:
	@go run ./cmd/api -db-dsn=${DB_DSN}

## bd/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DB_DSN}


## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

# db/migrations/up:
# 	@echo 'Running up migrations...'
# 	migrate -path ./migrations -database ${DB_DSN} up

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	ls

#	======================================================================	#
#	QUALITY CONTROL
# 	======================================================================	#

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

#	======================================================================	#
#	BUILD
# 	======================================================================	#

current_time = $(shell date +"%Y-%m-%dT%H:%M:%S")
git_description = $(shell git describe --always --dirty)
linker_flags = '-s -w -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api
