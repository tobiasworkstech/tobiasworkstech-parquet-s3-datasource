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
)

// Build builds the plugin (frontend + backend)
func Build() error {
	mg.Deps(BuildFrontend, BuildBackend)
	return nil
}

// BuildFrontend builds the frontend
func BuildFrontend() error {
	fmt.Println("Building frontend...")
	return sh.Run("npm", "run", "build")
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

	ldflags := "-w -s"

	return sh.RunWith(env, "go", "build",
		"-ldflags", ldflags,
		"-o", outputPath,
		"./pkg",
	)
}

// BuildDev builds for current platform only (faster for development)
func BuildDev() error {
	fmt.Println("Building for current platform...")

	binary := binaryName
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}

	outputPath := filepath.Join("dist", binary)

	return sh.Run("go", "build",
		"-o", outputPath,
		"./pkg",
	)
}

// Test runs all tests
func Test() error {
	fmt.Println("Running tests...")
	mg.Deps(TestFrontend, TestBackend)
	return nil
}

// TestFrontend runs frontend tests
func TestFrontend() error {
	return sh.Run("npm", "run", "test:ci")
}

// TestBackend runs Go tests
func TestBackend() error {
	return sh.Run("go", "test", "-v", "./...")
}

// Lint runs linters
func Lint() error {
	fmt.Println("Running linters...")

	if err := sh.Run("npm", "run", "lint"); err != nil {
		return err
	}

	return sh.Run("golangci-lint", "run", "./...")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning...")

	paths := []string{"dist", "node_modules/.cache"}
	for _, p := range paths {
		if err := os.RemoveAll(p); err != nil {
			return err
		}
	}

	return nil
}

// Package creates distributable zip files
func Package() error {
	mg.Deps(Build)

	fmt.Println("Packaging plugin...")

	// Create package directory
	pkgDir := filepath.Join("dist", pluginID)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return err
	}

	// Copy dist contents
	distFiles, err := filepath.Glob("dist/*")
	if err != nil {
		return err
	}

	for _, f := range distFiles {
		base := filepath.Base(f)
		if base == pluginID {
			continue
		}

		dest := filepath.Join(pkgDir, base)
		if err := copyFile(f, dest); err != nil {
			return err
		}
	}

	// Create zip
	zipPath := fmt.Sprintf("dist/%s.zip", pluginID)
	return sh.Run("zip", "-r", zipPath, pluginID, "-j", "dist/"+pluginID)
}

// Sign signs the plugin for distribution
func Sign() error {
	fmt.Println("Signing plugin...")

	rootURL := os.Getenv("GRAFANA_ACCESS_POLICY_TOKEN")
	if rootURL == "" {
		return fmt.Errorf("GRAFANA_ACCESS_POLICY_TOKEN environment variable not set")
	}

	return sh.Run("npx", "@grafana/sign-plugin@latest")
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}

// Dev runs development server with hot reload
func Dev() error {
	mg.Deps(BuildDev)

	// Start frontend watch in background
	cmd := exec.Command("npm", "run", "dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Install installs dependencies
func Install() error {
	fmt.Println("Installing dependencies...")

	if err := sh.Run("npm", "install"); err != nil {
		return err
	}

	return sh.Run("go", "mod", "download")
}
