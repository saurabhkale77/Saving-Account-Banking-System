run:
	go run main.go

clean:
	cls
	
cleanDB:
	rm ./repository/bank.db

test:
	go test -v ./...

test-cover:
	go test -v ./... -covermode=count -coverprofile=coverage.out