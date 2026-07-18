// Package authconfigfile persists non-secret authentication setup metadata.
// It is deliberately separate from every credential store.
package authconfigfile

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tasuku43/agentic-cli-foundry/internal/domain/authn"
)

var (
	ErrUnsafePath  = errors.New("authentication configuration path is unsafe")
	ErrInvalidData = errors.New("authentication configuration data is invalid")
)

// Store owns one injected non-secret configuration path.
type Store struct {
	path string
}

// New returns a store without creating files or directories.
func New(path string) *Store { return &Store{path: path} }

// Decode strictly decodes one bounded schema-versioned document.
func Decode(reader io.Reader) (authn.UserConfiguration, error) {
	if reader == nil {
		return authn.UserConfiguration{}, ErrInvalidData
	}
	limited := io.LimitReader(reader, authn.MaxUserConfigurationBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil || len(data) > authn.MaxUserConfigurationBytes {
		return authn.UserConfiguration{}, ErrInvalidData
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var configuration authn.UserConfiguration
	if err := decoder.Decode(&configuration); err != nil {
		return authn.UserConfiguration{}, ErrInvalidData
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return authn.UserConfiguration{}, ErrInvalidData
	}
	if err := configuration.Validate(); err != nil {
		return authn.UserConfiguration{}, ErrInvalidData
	}
	return configuration.Clone(), nil
}

// Encode validates and returns the canonical bounded document.
func Encode(configuration authn.UserConfiguration) ([]byte, error) {
	if err := configuration.Validate(); err != nil {
		return nil, ErrInvalidData
	}
	data, err := json.Marshal(configuration)
	if err != nil || len(data)+1 > authn.MaxUserConfigurationBytes {
		return nil, ErrInvalidData
	}
	return append(data, '\n'), nil
}

// Load reads a safe regular owner-only file. Missing is not an error; every
// corrupt or unsafe present state fails closed.
func (s *Store) Load(ctx context.Context) (authn.UserConfiguration, bool, error) {
	if err := validateStoreContext(ctx, s); err != nil {
		return authn.UserConfiguration{}, false, err
	}
	info, err := os.Lstat(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return authn.UserConfiguration{}, false, nil
	}
	if err != nil || !safeFileInfo(info) {
		return authn.UserConfiguration{}, true, ErrUnsafePath
	}
	file, err := os.Open(s.path)
	if err != nil {
		return authn.UserConfiguration{}, true, ErrUnsafePath
	}
	defer file.Close()
	opened, err := file.Stat()
	current, currentErr := os.Lstat(s.path)
	if err != nil || currentErr != nil || !safeFileInfo(current) || !os.SameFile(info, opened) || !os.SameFile(opened, current) {
		return authn.UserConfiguration{}, true, ErrUnsafePath
	}
	configuration, err := Decode(file)
	if err != nil {
		return authn.UserConfiguration{}, true, err
	}
	return configuration, true, nil
}

// Save atomically replaces the target through an owner-only same-directory
// temporary file. It never creates the parent directory.
func (s *Store) Save(ctx context.Context, configuration authn.UserConfiguration) (err error) {
	if err := validateStoreContext(ctx, s); err != nil {
		return err
	}
	data, err := Encode(configuration)
	if err != nil {
		return err
	}
	parent := filepath.Dir(s.path)
	parentInfo, err := os.Lstat(parent)
	if err != nil || parentInfo.Mode()&os.ModeSymlink != 0 || !parentInfo.IsDir() || parentInfo.Mode().Perm()&0o077 != 0 {
		return ErrUnsafePath
	}
	if targetInfo, statErr := os.Lstat(s.path); statErr == nil {
		if !safeFileInfo(targetInfo) {
			return ErrUnsafePath
		}
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return ErrUnsafePath
	}
	temporary, err := os.CreateTemp(parent, ".auth-config-*")
	if err != nil {
		return fmt.Errorf("create authentication configuration temporary file: %w", err)
	}
	temporaryPath := temporary.Name()
	defer func() {
		_ = temporary.Close()
		_ = os.Remove(temporaryPath)
	}()
	if err := temporary.Chmod(0o600); err != nil {
		return fmt.Errorf("set authentication configuration permissions: %w", err)
	}
	if _, err := temporary.Write(data); err != nil {
		return fmt.Errorf("write authentication configuration: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		return fmt.Errorf("sync authentication configuration: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close authentication configuration: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, s.path); err != nil {
		return fmt.Errorf("replace authentication configuration: %w", err)
	}
	return nil
}

// Status reconciles persistent state without writing or repairing it.
func (s *Store) Status(ctx context.Context) authn.ConfigurationStatus {
	configuration, present, err := s.Load(ctx)
	if err != nil {
		problem := "invalid_data"
		if errors.Is(err, ErrUnsafePath) {
			problem = "unsafe_file"
		}
		return authn.ConfigurationStatus{State: authn.ConfigurationStateInvalid, Problem: problem}
	}
	if !present {
		return authn.ConfigurationStatus{State: authn.ConfigurationStateMissing}
	}
	return authn.ConfigurationStatus{State: authn.ConfigurationStateValid, SchemaVersion: configuration.SchemaVersion, Method: configuration.Method}
}

func validateStoreContext(ctx context.Context, store *Store) error {
	if ctx == nil || store == nil || store.path == "" || !filepath.IsAbs(store.path) {
		return ErrUnsafePath
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}

func safeFileInfo(info os.FileInfo) bool {
	return info != nil && info.Mode()&os.ModeSymlink == 0 && info.Mode().IsRegular() && info.Mode().Perm() == 0o600
}
