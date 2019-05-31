# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gwtc android ios gwtc-cross swarm evm all test clean
.PHONY: gwtc-linux gwtc-linux-386 gwtc-linux-amd64 gwtc-linux-mips64 gwtc-linux-mips64le
.PHONY: gwtc-linux-arm gwtc-linux-arm-5 gwtc-linux-arm-6 gwtc-linux-arm-7 gwtc-linux-arm64
.PHONY: gwtc-darwin gwtc-darwin-386 gwtc-darwin-amd64
.PHONY: gwtc-windows gwtc-windows-386 gwtc-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

gwtc:
	build/env.sh go run build/ci.go install ./cmd/gwtc
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gwtc\" to launch gwtc."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gwtc.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Gwtc.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

clean:
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/jteeuwen/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go install ./cmd/abigen

# Cross Compilation Targets (xgo)

gwtc-cross: gwtc-linux gwtc-darwin gwtc-windows gwtc-android gwtc-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-*

gwtc-linux: gwtc-linux-386 gwtc-linux-amd64 gwtc-linux-arm gwtc-linux-mips64 gwtc-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-linux-*

gwtc-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gwtc
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-linux-* | grep 386

gwtc-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gwtc
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-linux-* | grep amd64

gwtc-linux-arm: gwtc-linux-arm-5 gwtc-linux-arm-6 gwtc-linux-arm-7 gwtc-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-linux-* | grep arm

gwtc-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gwtc
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-linux-* | grep arm-5

gwtc-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gwtc
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-linux-* | grep arm-6

gwtc-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gwtc
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-linux-* | grep arm-7

gwtc-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gwtc
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-linux-* | grep arm64

gwtc-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gwtc
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-linux-* | grep mips

gwtc-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gwtc
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-linux-* | grep mipsle

gwtc-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gwtc
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-linux-* | grep mips64

gwtc-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gwtc
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-linux-* | grep mips64le

gwtc-darwin: gwtc-darwin-386 gwtc-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-darwin-*

gwtc-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gwtc
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-darwin-* | grep 386

gwtc-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gwtc
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-darwin-* | grep amd64

gwtc-windows: gwtc-windows-386 gwtc-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-windows-*

gwtc-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gwtc
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-windows-* | grep 386

gwtc-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gwtc
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gwtc-windows-* | grep amd64
