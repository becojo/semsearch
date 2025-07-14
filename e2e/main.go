package main

import (
	"context"
	"dagger/e-2-e/internal/dagger"
	"fmt"
)

type E2E struct {
}

type WithSemsearch struct {
	Semsearch *dagger.File
}

const (
	GoVersion             = "1.24.2"
	SemgrepVersion        = "1.124.1"
	OpengrepVersion       = "1.6.0"
	OpengrepInstallScript = "https://raw.githubusercontent.com/opengrep/opengrep/0b89f9df946a0320b0b765e4f7f88b1711104b9d/install.sh"
)

func (m *E2E) FromSource(
	//+ignore=["*", "!pkg/**/*.go", "!cmd/**/*.go", "**_test.go", "!go.mod", "!go.sum"]
	src *dagger.Directory,
) *WithSemsearch {
	ctr := m.Semsearch(src, GoVersion)
	return &WithSemsearch{
		ctr.File("/usr/bin/semsearch"),
	}
}

func (m *E2E) FromBinary(
	file *dagger.File,
) *WithSemsearch {
	return &WithSemsearch{file}
}

func (m *E2E) Semsearch(
	//+ignore=["*", "!pkg/**/*.go", "!cmd/**/*.go", "**_test.go", "!go.mod", "!go.sum"]
	src *dagger.Directory,
	//+default=""
	goVersion string,
) *dagger.Container {
	if goVersion == "" {
		goVersion = GoVersion
	}

	goCache := dag.CacheVolume("go_cache" + goVersion)
	buildCache := dag.CacheVolume("go_build" + goVersion)
	ctr := dag.Container().
		From(fmt.Sprintf("golang:%s-alpine", goVersion)).
		WithEnvVariable("CGO_ENABLED", "0").
		WithMountedCache("/go/pkg/mod", goCache).
		WithMountedCache("/root/.cache/go-build", buildCache).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "/usr/bin/semsearch", "./cmd/semsearch"})

	return ctr
}

func (m *WithSemsearch) Semgrep(
	ctx context.Context,
	//+default=""
	version string,
) *dagger.Container {
	if version == "" {
		version = SemgrepVersion
	}
	return dag.Container().
		From(fmt.Sprintf("semgrep/semgrep:%s", version)).
		WithEnvVariable("SEMGREP_SEND_METRICS", "off").
		WithFile("/usr/bin/semsearch", m.Semsearch).
		WithEnvVariable("SEMSEARCH_COMMAND", "semgrep")
}

func (m *WithSemsearch) Opengrep(
	ctx context.Context,
	//+default=""
	version string,
) *dagger.Container {
	if version == "" {
		version = OpengrepVersion
	}

	return dag.Container().
		From("alpine:latest").
		WithExec(append([]string{"apk", "add", "--no-cache", "curl", "cosign"})).
		WithEnvVariable("OPENGREP_INSTALL_SCRIPT_URL", OpengrepInstallScript).
		WithExec([]string{"sh", "-ec", `
			curl -sSL "$OPENGREP_INSTALL_SCRIPT_URL" > "/usr/bin/opengrep-install.sh"
		`}).
		WithExec([]string{"sh", "/usr/bin/opengrep-install.sh", "-v", "v" + version}).
		WithSymlink("/root/.opengrep/cli/latest/opengrep", "/usr/local/bin/opengrep").
		WithFile("/usr/bin/semsearch", m.Semsearch)
}
