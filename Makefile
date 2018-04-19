default: mock

mock:
	#mockgen -package=assert testing TB > _mocktesting.go
	mockgen -package=assert github.com/arteev/tag-assert TB > _mocktesting.go
	rm -f mocktesting.go
	mv _mocktesting.go mocktesting.go