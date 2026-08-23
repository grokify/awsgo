package inspector2util

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
)

func writeTempFile(t *testing.T, contents string) string {
	t.Helper()
	name := filepath.Join(t.TempDir(), "findings.json")
	if err := os.WriteFile(name, []byte(contents), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	return name
}

func TestReadFileListFindingsOutputStream_Success(t *testing.T) {
	filename := writeTempFile(t, `{"findings":[
		{"awsAccountId":"111111111111","description":"d1","findingArn":"arn:1"},
		{"awsAccountId":"222222222222","description":"d2","findingArn":"arn:2"}
	],"nextToken":"abc"}`)

	var got []types.Finding
	err := ReadFileListFindingsOutputStream(filename, func(f types.Finding) error {
		got = append(got, f)
		return nil
	})
	if err != nil {
		t.Fatalf("ReadFileListFindingsOutputStream() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if arn := *got[0].FindingArn; arn != "arn:1" {
		t.Errorf("got[0].FindingArn = %q, want %q", arn, "arn:1")
	}
	if arn := *got[1].FindingArn; arn != "arn:2" {
		t.Errorf("got[1].FindingArn = %q, want %q", arn, "arn:2")
	}
}

func TestReadFileListFindingsOutputStream_KeyOrderIndependent(t *testing.T) {
	filename := writeTempFile(t, `{"nextToken":"abc","findings":[
		{"awsAccountId":"111111111111","description":"d1","findingArn":"arn:1"}
	]}`)

	var got []types.Finding
	err := ReadFileListFindingsOutputStream(filename, func(f types.Finding) error {
		got = append(got, f)
		return nil
	})
	if err != nil {
		t.Fatalf("ReadFileListFindingsOutputStream() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
}

func TestReadFileListFindingsOutputStream_NoFindingsKey(t *testing.T) {
	filename := writeTempFile(t, `{"nextToken":"abc"}`)

	called := false
	err := ReadFileListFindingsOutputStream(filename, func(f types.Finding) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("ReadFileListFindingsOutputStream() error = %v", err)
	}
	if called {
		t.Error("callback was invoked with no findings present")
	}
}

func TestReadFileListFindingsOutputStream_NilFunc(t *testing.T) {
	filename := writeTempFile(t, `{"findings":[]}`)

	err := ReadFileListFindingsOutputStream(filename, nil)
	if err == nil {
		t.Fatal("ReadFileListFindingsOutputStream() error = nil, want error for nil fn")
	}
}

func TestReadFileListFindingsOutputStream_FileNotFound(t *testing.T) {
	err := ReadFileListFindingsOutputStream(filepath.Join(t.TempDir(), "missing.json"), func(f types.Finding) error {
		return nil
	})
	if err == nil {
		t.Fatal("ReadFileListFindingsOutputStream() error = nil, want error for missing file")
	}
}

func TestReadFileListFindingsOutputStream_NotAnObject(t *testing.T) {
	filename := writeTempFile(t, `[1,2,3]`)

	err := ReadFileListFindingsOutputStream(filename, func(f types.Finding) error {
		return nil
	})
	if err == nil {
		t.Fatal("ReadFileListFindingsOutputStream() error = nil, want error for non-object root")
	}
}

func TestReadFileListFindingsOutputStream_CallbackErrorPropagates(t *testing.T) {
	filename := writeTempFile(t, `{"findings":[
		{"awsAccountId":"111111111111","description":"d1","findingArn":"arn:1"},
		{"awsAccountId":"222222222222","description":"d2","findingArn":"arn:2"}
	]}`)

	wantErr := errors.New("boom")
	callCount := 0
	err := ReadFileListFindingsOutputStream(filename, func(f types.Finding) error {
		callCount++
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("ReadFileListFindingsOutputStream() error = %v, want %v", err, wantErr)
	}
	if callCount != 1 {
		t.Errorf("callCount = %d, want 1 (should stop at first error)", callCount)
	}
}
