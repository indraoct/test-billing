package main

import (
	"fmt"
	"runtime"
	"test-billing/config"
)

func main() {
	goVersion := runtime.Version()
	osArch := fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)

	// Read config and check if error occurred
	conf := config.Load()

	fmt.Println("App Name	:", conf.Name)
	fmt.Println("Go Version	:", goVersion)
	fmt.Println("OS Arch	:", osArch)

}
