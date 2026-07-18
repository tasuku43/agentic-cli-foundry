package authconfigfile

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/tasuku43/agentic-cli-foundry/internal/domain/authn"
)

func validConfiguration(method authn.Method) authn.UserConfiguration {
	return authn.UserConfiguration{
		SchemaVersion: authn.UserConfigurationSchemaVersion,
		Method:        method,
		Parameters:    []authn.PublicParameter{{Name: "public_client_id", Value: "example-public-client"}},
	}
}

func TestCodecStrictlyDecodesVersionedBoundedFixture(t *testing.T) {
	data, err := os.ReadFile("testdata/user-configuration-v1.json")
	if err != nil {
		t.Fatal(err)
	}
	configuration, err := Decode(bytes.NewReader(data))
	if err != nil || configuration.Method != authn.MethodOAuth2 || len(configuration.Parameters) != 2 {
		t.Fatalf("Decode() = %+v, %v", configuration, err)
	}
	for name, document := range map[string]string{
		"unknown field":  `{"schema_version":1,"method":"pat","parameters":[],"token":"forbidden"}`,
		"unknown schema": `{"schema_version":2,"method":"pat","parameters":[]}`,
		"trailing value": `{"schema_version":1,"method":"pat","parameters":[]} {}`,
		"missing list":   `{"schema_version":1,"method":"pat"}`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := Decode(strings.NewReader(document)); !errors.Is(err, ErrInvalidData) {
				t.Fatalf("Decode() error = %v", err)
			}
		})
	}
	oversized := bytes.Repeat([]byte("x"), authn.MaxUserConfigurationBytes+1)
	if _, err := Decode(bytes.NewReader(oversized)); !errors.Is(err, ErrInvalidData) {
		t.Fatalf("oversized Decode() error = %v", err)
	}
}

func TestStoreAtomicallyReplacesAndReportsConfiguration(t *testing.T) {
	directory := t.TempDir()
	if err := os.Chmod(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(directory, "auth.json")
	store := New(path)
	if status := store.Status(context.Background()); status.State != authn.ConfigurationStateMissing {
		t.Fatalf("missing status = %+v", status)
	}
	first := validConfiguration(authn.MethodPAT)
	if err := store.Save(context.Background(), first); err != nil {
		t.Fatal(err)
	}
	info, err := os.Lstat(path)
	if err != nil || info.Mode().Perm() != 0o600 || !info.Mode().IsRegular() {
		t.Fatalf("saved file = %+v, %v", info, err)
	}
	second := validConfiguration(authn.MethodOAuth2)
	if err := store.Save(context.Background(), second); err != nil {
		t.Fatal(err)
	}
	loaded, present, err := store.Load(context.Background())
	if err != nil || !present || !reflect.DeepEqual(loaded, second) {
		t.Fatalf("Load() = %+v, %t, %v", loaded, present, err)
	}
	status := store.Status(context.Background())
	if status.State != authn.ConfigurationStateValid || status.Method != authn.MethodOAuth2 || status.SchemaVersion != 1 {
		t.Fatalf("valid status = %+v", status)
	}
	entries, err := os.ReadDir(directory)
	if err != nil || len(entries) != 1 {
		t.Fatalf("temporary residue = %+v, %v", entries, err)
	}
}

func TestStoreRejectsUnsafeAndCorruptFilesWithoutRepair(t *testing.T) {
	directory := t.TempDir()
	if err := os.Chmod(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(directory, "auth.json")
	if err := os.WriteFile(target, []byte(`{"schema_version":2}`), 0o600); err != nil {
		t.Fatal(err)
	}
	store := New(target)
	before, _ := os.ReadFile(target)
	if _, present, err := store.Load(context.Background()); !present || !errors.Is(err, ErrInvalidData) {
		t.Fatalf("corrupt Load() present=%t error=%v", present, err)
	}
	if status := store.Status(context.Background()); status.State != authn.ConfigurationStateInvalid || status.Problem != "invalid_data" {
		t.Fatalf("corrupt status = %+v", status)
	}
	after, _ := os.ReadFile(target)
	if !bytes.Equal(before, after) {
		t.Fatal("read-only status repaired corrupt state")
	}

	if err := os.Chmod(target, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.Load(context.Background()); !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("permissive file error = %v", err)
	}
	if err := os.Remove(target); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(directory, "missing"), target); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.Load(context.Background()); !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("symlink error = %v", err)
	}
	if err := store.Save(context.Background(), validConfiguration(authn.MethodPAT)); !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("save over symlink error = %v", err)
	}
}
