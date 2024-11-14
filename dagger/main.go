package main

import (
	"context"

	"dagger/myboth/internal/dagger"
)

type Myboth struct{}

// Lint runs golangci-lint on the project with all linters enabled
func (m *Myboth) Lint(ctx context.Context, source *dagger.Directory) *dagger.Container {
	return dag.GolangciLint().
		WithModuleCache(dag.CacheVolume("gomod")).
		WithBuildCache(dag.CacheVolume("gobuild")).
		WithLinterCache(dag.CacheVolume("golangci")).
		Run(source, dagger.GolangciLintRunOpts{
			RawArgs: []string{
				"--issues-exit-code=0",
				"--enable-all",
			},
		})
}

// Returns a container that echoes whatever string argument is provided
func (m *Myboth) ContainerEcho(stringArg string) *dagger.Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
}

// Returns lines that match a pattern in the files of the provided Directory
func (m *Myboth) GrepDir(ctx context.Context, directoryArg *dagger.Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}
