//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
var Default = Build

const (
	pluginID   = "tobiasworkstech-parquet-s3-datasource"
	binaryName = "gpx_parquet_s3_datasource"
	pluginDir  = "plugin"
)

// Build builds the plugin (frontend + backend)
func Build() error {
	mg.Deps(BuildFrontend, BuildBackend)
	return nil
}

// BuildFrontend builds the frontend
func BuildFrontend() error {
	fmt.Println("Building frontend...")
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// BuildBackend builds the Go backend for multiple platforms
func BuildBackend() error {
	fmt.Println("Building backend...")

	targets := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
	}

	for _, t := range targets {
		if err := buildBinary(t.os, t.arch); err != nil {
			return err
		}
	}

	return nil
}

func buildBinary(targetOS, targetArch string) error {
	binary := binaryName
	if targetOS == "windows" {
		binary += ".exe"
	}

	suffix := fmt.Sprintf("_%s_%s", targetOS, targetArch)
	if targetOS == "windows" {
		suffix = fmt.Sprintf("_%s_%s.exe", targetOS, targetArch)
	}

	outputPath := filepath.Join("dist", binary+suffix)
	fmt.Printf("Building %s...\n", outputPath)

	env := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        targetOS,
		"GOARCH":      targetArch,
	}

	cmd := exec.Command("go", "build",
		"-ldflags", "-w -s",
		"-o", outputPath,
		"./pkg",
	)
	cmd.Dir = pluginDir
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// BuildDev builds for current platform only (faster for development)
func BuildDev() error {
	fmt.Println("Building for current platform...")

	binary := binaryName
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}

	outputPath := filepath.Join("dist", binary)

	cmd := exec.Command("go", "build",
		"-o", outputPath,
		"./pkg",
	)
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Test runs all tests
func Test() error {
	fmt.Println("Running tests...")
	mg.Deps(TestFrontend, TestBackend)
	return nil
}

// TestFrontend runs frontend tests
func TestFrontend() error {
	cmd := exec.Command("npm", "run", "test:ci")
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// TestBackend runs Go tests
func TestBackend() error {
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Lint runs linters
func Lint() error {
	fmt.Println("Running linters...")

	cmd := exec.Command("npm", "run", "lint")
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("golangci-lint", "run", "./...")
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning...")
	return sh.Rm(filepath.Join(pluginDir, "dist"))
}

// Sign signs the plugin for distribution
func Sign() error {
	fmt.Println("Signing plugin...")

	token := os.Getenv("GRAFANA_ACCESS_POLICY_TOKEN")
	if token == "" {
		return fmt.Errorf("GRAFANA_ACCESS_POLICY_TOKEN environment variable not set")
	}

	cmd := exec.Command("npx", "@grafana/sign-plugin@latest")
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Dev runs development server with hot reload
func Dev() error {
	mg.Deps(BuildDev)

	cmd := exec.Command("npm", "run", "dev")
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Install installs dependencies
func Install() error {
	fmt.Println("Installing dependencies...")

	cmd := exec.Command("npm", "install")
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("go", "mod", "download")
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
