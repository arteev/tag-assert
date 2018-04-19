default: test

test:
	go test -v ./... 

mock:
	mockgen -package=assert github.com/arteev/tag-assert TB > _mocktesting.go
	rm -f mocktesting.go
	mv _mocktesting.go mocktesting.go

example:
	go test -v ./_example/
