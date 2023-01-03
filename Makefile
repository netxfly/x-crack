
$env:OPTS=-ldflags "-w -s"
$env:WINDOWS="build/xcrack.exe"
$env:LINUX="build/xcrack_linux_amd64"
$env:MAC="build/xcrack_mac"

build:

	# linux
	$env:CGO_ENABLED="0"
	$env:GOOS="linux"
	$env:GOARCH="amd64"
	go build $env:OPTS -o $env:LINUX
	# windows
	$env:CGO_ENABLED="0"
	$env:GOOS="windows"
	$env:GOARCH="amd64"
	go build $env:OPTS  -o $env:WINDOWS

	# mac
	$env:CGO_ENABLED="0"
	$env:GOOS="darwin"
	$env:GOARCH="arm64"
	go build -o $env:MAC $env:OPTS

clean:
	rm -f *.o test