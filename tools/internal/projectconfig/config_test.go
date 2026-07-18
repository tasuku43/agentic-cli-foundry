package projectconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsUnknownFields(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".harness"), 0o755); err != nil {
		t.Fatal(err)
	}
	raw := `{"schema_version":1,"profile":"template","project":{},"public_guard":{},"unknown":true}`
	if err := os.WriteFile(filepath.Join(root, ".harness", "project.json"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(root); err == nil || !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("error = %v", err)
	}
}

func TestReadyProblemsRequireProjectSpecificIdentity(t *testing.T) {
	wantTemplateProblems := []string{
		"name still uses the runnable template default",
		"binary_name still uses the runnable template default",
		"go_module still uses the runnable template default",
		"github_repository still uses the runnable template default",
		"description still uses the runnable template default",
		"formula_class still uses the runnable template default",
		"security_contact still uses the runnable template default",
	}
	if got := ReadyProblems(Defaults); strings.Join(got, "\n") != strings.Join(wantTemplateProblems, "\n") {
		t.Fatalf("problems = %v, want %v", got, wantTemplateProblems)
	}
	project := Defaults
	project.Name = "Example Tool"
	project.BinaryName = "example-tool"
	project.GoModule = "github.com/" + Defaults.GitHubOwner + "/example-tool"
	project.GitHubRepository = "example-tool"
	project.Description = "An example command-line tool."
	project.FormulaClass = "ExampleTool"
	project.SecurityContact = "security@example.com"
	if got := ReadyProblems(project); len(got) != 0 {
		t.Fatalf("problems = %v", got)
	}
}

func TestReadyProblemsRequireEachMeaningfulDerivedField(t *testing.T) {
	project := Defaults
	project.Name = "Example Tool"
	project.BinaryName = "example-tool"
	project.GoModule = "github.com/" + Defaults.GitHubOwner + "/example-tool"
	project.GitHubRepository = "example-tool"
	project.Description = "An example command-line tool."
	project.FormulaClass = "ExampleTool"
	project.SecurityContact = "security@example.com"

	tests := []struct {
		name    string
		restore func(*Project)
	}{
		{name: "name", restore: func(p *Project) { p.Name = Defaults.Name }},
		{name: "binary_name", restore: func(p *Project) { p.BinaryName = Defaults.BinaryName }},
		{name: "go_module", restore: func(p *Project) { p.GoModule = Defaults.GoModule }},
		{name: "github_repository", restore: func(p *Project) { p.GitHubRepository = Defaults.GitHubRepository }},
		{name: "description", restore: func(p *Project) { p.Description = Defaults.Description }},
		{name: "formula_class", restore: func(p *Project) { p.FormulaClass = Defaults.FormulaClass }},
		{name: "security_contact", restore: func(p *Project) { p.SecurityContact = Defaults.SecurityContact }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := project
			test.restore(&candidate)
			got := ReadyProblems(candidate)
			want := test.name + " still uses the runnable template default"
			if len(got) != 1 || got[0] != want {
				t.Fatalf("problems = %v, want [%q]", got, want)
			}
		})
	}
}

func TestConfigRejectsWindowsReservedBinaryNames(t *testing.T) {
	config := Config{
		SchemaVersion: SchemaVersion,
		Profile:       "template",
		Project:       Defaults,
		PublicGuard:   PublicGuard{DenylistFile: ".harness/denylist.txt"},
	}
	for _, name := range []string{"con", "aux", "prn", "nul", "com1", "com9", "lpt1", "lpt9", "Con", "cOm1", "LpT9"} {
		config.Project.BinaryName = name
		if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "reserved Windows device") {
			t.Fatalf("Validate() accepted binary_name %q: %v", name, err)
		}
	}
	for _, name := range []string{"console", "auxiliary", "null", "com0", "com10", "lpt0", "lpt10"} {
		config.Project.BinaryName = name
		if err := config.Validate(); err != nil {
			t.Fatalf("Validate() rejected binary_name %q: %v", name, err)
		}
	}
	for _, name := range []string{"Con", "cOm1", "LpT9"} {
		if !isWindowsReservedBaseName(name) {
			t.Fatalf("isWindowsReservedBaseName(%q) = false", name)
		}
	}
}
