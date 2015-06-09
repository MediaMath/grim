package grim

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var pathsNames = &pathNames{} // a global variable to store file paths of workspace and result

/**
Test on the consistency of timestamp of workspace and result.
*/
func TestOnWorkSpaceAndResultNameConsistency(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "results-dir-consistencyCheck")
	defer os.RemoveAll(tempDir)
	//trigger a build to have file paths of result and workspace
	builtForHook(tempDir, "MediaMath", "grim", 0)

	isMatched, err := pathsNames.isConsistent()
	if err != nil {
		t.Fatal(err.Error())
	}
	if !isMatched {
		t.Fatalf("inconsistent dir name")
	}
}

type pathNames struct {
	workspacePath string
	resultPath    string
}

func (pn *pathNames) isConsistent() (bool, error) {
	workspacePaths := strings.Split(pn.workspacePath, "/")
	resultPaths := strings.Split(pn.resultPath, "/")
	a := workspacePaths[len(workspacePaths)-1]
	b := resultPaths[len(resultPaths)-1]

	if len(a) == 0 {
		return false, fmt.Errorf("empty workspacePaths name ")
	}

	if len(b) == 0 {
		return false, fmt.Errorf("empty resultPaths name ")
	}

	if len(a) != len(b) || !strings.EqualFold(a, b) {
		return false, fmt.Errorf("inconsistent timestamp workspace:" + a + " and resultpath:" + b)
	}

	return true, nil
}

func builtForHook(tempDir, owner, repo string, exitCode int) error {
	return onHook("not-used", &effectiveConfig{resultRoot: tempDir, workspaceRoot: tempDir}, hookEvent{Owner: owner, Repo: repo}, StubBuild)
}

func StubBuild(configRoot string, resultPath string, config *effectiveConfig, hook hookEvent, basename string) (*executeResult, string, error) {
	pathsNames.resultPath = resultPath
	return built(config.gitHubToken, configRoot, config.workspaceRoot, resultPath, config.pathToCloneIn, hook.Owner, hook.Repo, hook.Ref, hook.env(), basename)
}

func built(token, configRoot, workspaceRoot, resultPath, clonePath, owner, repo, ref string, extraEnv []string, basename string) (*executeResult, string, error) {
	ws := &testWorkSpaceBuilder{workspaceBuilder{workspaceRoot, clonePath, token, configRoot, owner, repo, ref, extraEnv}}
	return grimBuild(ws, resultPath, basename)
}

type testWorkSpaceBuilder struct {
	workspaceBuilder
}

func (tb *testWorkSpaceBuilder) PrepareWorkspace(basename string) (string, error) {
	workSpacePath, err := createWorkspaceDirectory(tb.workspaceRoot, tb.owner, tb.repo, basename)
	pathsNames.workspacePath = workSpacePath
	return workSpacePath, err
}
