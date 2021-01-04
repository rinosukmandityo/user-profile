rm coverage.txt || true 2> /dev/null

go test -p=1 -coverprofile=profile.out ./... | tee -a coverage.txt