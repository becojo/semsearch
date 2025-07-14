package rule

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Runner struct {
	// current rule state
	state *State
	// temporary directory with rules, evals
	tmpDir string
	// additional paths to scan
	paths []string
}

func NewRunner(state *State) *Runner {
	return &Runner{
		state: state,
	}
}

func (r *Runner) Prepare() error {
	if err := r.createTempDir(); err != nil {
		return err
	}

	if _, err := r.writeTempFile("rules.yaml", r.state.MarshalRules()); err != nil {
		return err
	}

	for i, eval := range r.state.evals {
		p, err := r.writeTempFile(fmt.Sprintf("eval-%d", i), []byte(eval))
		if err != nil {
			return err
		}
		r.paths = append(r.paths, p)
	}

	return nil
}

func (r *Runner) createTempDir() error {
	tmpDir, err := os.MkdirTemp("", "semsearch-")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	r.tmpDir = tmpDir
	return nil
}

func (r *Runner) writeTempFile(name string, content []byte) (string, error) {
	filePath := path.Join(r.tmpDir, name)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return "", fmt.Errorf("failed to write file in temporary directory: %w", err)
	}
	return filePath, nil
}

func (r *Runner) Cleanup() error {
	return os.RemoveAll(r.tmpDir)
}

func (r *Runner) Args() []string {
	args := []string{
		"scan",
		"--no-rewrite-rule-ids",
		"--disable-version-check",
		fmt.Sprintf("--%s", r.state.format),
	}

	for _, config := range r.state.configs {
		args = append(args, "--config", config)
	}

	args = append(args, "--config", path.Join(r.tmpDir, "rules.yaml"))

	if len(r.state.evals) > 0 {
		args = append(args, "--scan-unknown-extensions")
	}

	if r.state.autofix {
		args = append(args, "--autofix")
	}

	if r.state.verbose {
		args = append(args, "--verbose")
	} else {
		args = append(args, "--quiet")
	}

	args = append(args, r.paths...)
	args = append(args, r.state.paths...)

	return args
}

func (r *Runner) Run() error {
	if len(r.state.warnings) > 0 {
		for _, warning := range r.state.warnings {
			fmt.Fprintln(os.Stderr, "Warning:", warning)
		}
	}

	if r.state.debug {
		fmt.Fprintln(os.Stderr, string(r.state.MarshalRules()))
		fmt.Fprintf(os.Stderr, "command: %s %s\n", r.state.command, strings.Join(r.Args(), " "))
	}

	if r.state.export {
		fmt.Print(string(r.state.MarshalRules()))
		return nil
	}

	cmd := exec.Command(r.state.command, r.Args()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
