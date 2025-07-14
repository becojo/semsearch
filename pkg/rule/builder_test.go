package rule

import (
	"testing"
)

func TestState(t *testing.T) {

	runner := &Runner{
		state: Builder().
			Rule().
			Language("go").
			Eval("package main\nvar x = 123").
			Format("json").
			Pattern("package $PKG\n...\nvar $VAR = $VAL").
			Export(),
	}

	runner.Prepare()
	runner.Run()
	runner.Cleanup()
}
