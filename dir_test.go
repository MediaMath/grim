package grim
import (
	"fmt"
	"testing"
	"io/ioutil"
	"os"
	"strings"
)

var pathsNames = &pathNames{} // use to store file  workspace pth and result path

func TestOnWorkSpaceAndResultNameConsistency(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "results-dir-error")
	defer os.RemoveAll(tempDir)
	builtForHook(tempDir, "MediaMath", "grim", 0)
	ismatched, err := pathsNames.isConsistent()
	if err != nil{
		t.Fatal(err.Error())
	}
	if !ismatched{
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
	if len(a)==0 || len(b) == 0 {
		return false, fmt.Errorf("empty dir name ")
	}
	if len(a)!=len(b) {
		return false, fmt.Errorf("inconsistent dir names")
	}
return strings.EqualFold(a, b), nil
}


func builtForHook(tempDir, owner, repo string, exitCode int) error {
	return onHook("not-used", &effectiveConfig{resultRoot: tempDir, workspaceRoot:tempDir}, hookEvent{Owner: owner, Repo: repo}, StubBuild)
}

func StubBuild(configRoot string, resultPath string, config *effectiveConfig, hook hookEvent, basename string) (*executeResult, string, error) {
	//fmt.Print("result path:"+resultPath+"\n")
	pathsNames.resultPath=resultPath
	return built(config.gitHubToken, configRoot, config.workspaceRoot, resultPath, config.pathToCloneIn, hook.Owner, hook.Repo, hook.Ref, hook.env(), basename)
}

func built(token, configRoot, workspaceRoot, resultPath, clonePath, owner, repo, ref string, extraEnv []string, basename string) (*executeResult, string, error) {
	ws := &testworkSpaceBuilder{workspaceBuilder{workspaceRoot, clonePath, token, configRoot, owner, repo, ref, extraEnv}}
	return grimBuild(ws, resultPath, basename)
}

type testworkSpaceBuilder struct {
	workspaceBuilder
}

func (tb *testworkSpaceBuilder) PrepareWorkspace(basename string) (string, error) {
	workSpacePath, err := createWorkspaceDirectory(tb.workspaceRoot, tb.owner, tb.repo, basename)
	//fmt.Print("workspace:"+workSpacePath)
	pathsNames.workspacePath=workSpacePath
	return workSpacePath, err
}