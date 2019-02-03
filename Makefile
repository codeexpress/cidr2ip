all:
	GOOS=windows GOARCH=386 go build -o binaries/cidr2ip-win32.exe cidr2ip.go
	GOOS=windows GOARCH=amd64 go build -o binaries/cidr2ip-win64.exe cidr2ip.go
	GOOS=linux GOARCH=386 go build  -o binaries/cidr2ip-linux32 cidr2ip.go
	GOOS=linux GOARCH=amd64 go build -o binaries/cidr2ip-linux64 cidr2ip.go
	GOOS=darwin GOARCH=386 go build -o binaries/cidr2ip-osx32 cidr2ip.go
	GOOS=darwin GOARCH=amd64 go build -o binaries/cidr2ip-osx64 cidr2ip.go
