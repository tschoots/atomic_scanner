
rm -rf atomic_scanner

GO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o atomic_scanner
