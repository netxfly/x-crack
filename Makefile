# linux
$env:CGO_ENABLED="0"
$env:GOOS="linux"
$env:GOARCH="amd64"

go build -ldflags "-w -s" -o xcrack_linux_amd64/xcrack

# windows
$env:CGO_ENABLED="0"
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -ldflags "-w -s" -o xcrack_windows_amd64.exe