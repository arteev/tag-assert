default: test lint

test:
	go test -v ./... 

mock:
	mockgen -package=assert github.com/arteev/tag-assert TB > _mocktesting.go
	rm -f mocktesting.go
	mv _mocktesting.go mocktesting.go

lint:
	gometalinter.v2 --enable-all --exclude mock* ./...

example:
	go test -v ./_example/
