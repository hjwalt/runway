test:
	go test ./... -cover -coverprofile cover.out 
	
testv:
	go test ./... -cover -coverprofile cover.out -v 

cov: test
	go tool cover -func cover.out

htmlcov: test
	go tool cover -html cover.out -o cover.html

# --------------------

tidy:
	go mod tidy
	go fmt ./...

update:
	go get -u ./...
	go mod tidy
	go fmt ./...

