package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/tasuku43/agentic-cli-foundry/internal/domain/fault"
)

func TestRootHelpIsDerivedFromCatalog(t *testing.T) {
	var stdout, stderr bytes.Buffer
	command := New(strings.NewReader(""), &stdout, &stderr)
	if code := runCLI(command, []string{"help"}); code != ExitOK {
		t.Fatalf("Run(help) code = %d, stderr = %q", code, stderr.String())
	}
	output := stdout.String()
	for _, spec := range command.catalog.Commands() {
		if !strings.Contains(output, spec.Path) || !strings.Contains(output, spec.Summary) {
			t.Errorf("root help does not contain catalog entry %+v\n%s", spec, output)
		}
	}
}

func TestCommandHelpUsesCatalogMetadataAndDerivedReferences(t *testing.T) {
	var stdout, stderr bytes.Buffer
	command := New(strings.NewReader(""), &stdout, &stderr)
	if code := runCLI(command, []string{"sample", "read", "--help"}); code != ExitOK {
		t.Fatalf("Run(sample read --help) code = %d, stderr = %q", code, stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{
		"Usage:\n  agentic-cli-foundry sample read --id <sample-id> [--format tsv|json]",
		"Read exactly one offline sample by opaque ID.",
		"Effect: read",
		"Role: act",
		"Consumes reference: sample from input --id",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("command help lacks %q\n%s", want, output)
		}
	}
}

func TestAgentAndHumanHelpPublishFixedTarget(t *testing.T) {
	spec := fixedTargetActSpec("auth status")
	help, found := DefaultCatalog().Lookup("help")
	if !found {
		t.Fatal("default catalog lacks help")
	}
	catalog := NewCatalog(help, spec)
	var stdout, stderr bytes.Buffer
	command := newCLI(strings.NewReader(""), &stdout, &stderr, catalog, nil)
	if code := runCLI(command, []string{"help", "auth", "status", "--format=agent"}); code != ExitOK {
		t.Fatalf("agent help code = %d, stderr = %q", code, stderr.String())
	}
	var document agentDocument
	if err := json.Unmarshal(stdout.Bytes(), &document); err != nil {
		t.Fatal(err)
	}
	if len(document.Commands) != 1 || document.Commands[0].Contract.FixedTarget == nil ||
		*document.Commands[0].Contract.FixedTarget != *spec.Agent.FixedTarget ||
		len(document.Commands[0].ProducesRefs) != 0 || len(document.Commands[0].ConsumesRefs) != 0 {
		t.Fatalf("fixed-target agent projection = %+v", document.Commands)
	}

	stdout.Reset()
	stderr.Reset()
	if code := runCLI(command, []string{"auth", "status", "--help"}); code != ExitOK {
		t.Fatalf("human help code = %d, stderr = %q", code, stderr.String())
	}
	for _, want := range []string{"Fixed target:", "auth-config", "selected", "tool_local", "selected authentication configuration"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("human help lacks %q: %s", want, stdout.String())
		}
	}
}

func TestRootAgentHelpIsACompactProjectionOfTheCatalog(t *testing.T) {
	var stdout, stderr bytes.Buffer
	command := New(strings.NewReader(""), &stdout, &stderr)
	if code := runCLI(command, []string{"help", "--format", "agent"}); code != ExitOK {
		t.Fatalf("Run(agent help) code = %d, stderr = %q", code, stderr.String())
	}

	var document agentIndexDocument
	if err := json.Unmarshal(stdout.Bytes(), &document); err != nil {
		t.Fatalf("agent help is not JSON: %v\n%s", err, stdout.String())
	}
	if document.SchemaVersion != 5 || agentHelpSchemaVersion != 5 || document.View != "index" || document.Program != ProgramName {
		t.Fatalf("agent document header = %+v", document)
	}
	if document.ScopeRequest.InvocationTemplate != "agentic-cli-foundry help <command-or-namespace> --format agent" ||
		!reflect.DeepEqual(document.ScopeRequest.SelectorFields, []string{"commands[].path", "commands[].namespace"}) ||
		document.ScopeRequest.UnknownOutcomeMaxInvocations != 2 || document.ScopeRequest.KnownPathMaxInvocations != 1 {
		t.Fatalf("scope request = %+v", document.ScopeRequest)
	}
	specs := command.catalog.Commands()
	if len(document.Commands) != len(specs) {
		t.Fatalf("agent commands = %d, catalog commands = %d", len(document.Commands), len(specs))
	}
	for index, spec := range specs {
		got := document.Commands[index]
		if got.Path != spec.Path || got.Namespace != commandNamespace(spec.Path) || got.Summary != spec.Summary ||
			got.CapabilityID != spec.Agent.CapabilityID || got.Outcome != spec.Agent.Outcome ||
			got.Effect != spec.Effect.String() || got.Role != spec.Role.String() {
			t.Errorf("agent command %d = %+v, want catalog %+v", index, got, spec)
		}
	}
}

func TestScopedAgentHelpIsACompleteProjectionOfEveryCatalogCommand(t *testing.T) {
	command := New(strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	for _, spec := range command.catalog.Commands() {
		t.Run(strings.ReplaceAll(spec.Path, " ", "_"), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			selected := New(strings.NewReader(""), &stdout, &stderr)
			args := append([]string{"help"}, strings.Fields(spec.Path)...)
			args = append(args, "--format=agent")
			if code := runCLI(selected, args); code != ExitOK {
				t.Fatalf("Run(%v) code = %d, stderr = %q", args, code, stderr.String())
			}
			var document agentDocument
			if err := json.Unmarshal(stdout.Bytes(), &document); err != nil {
				t.Fatalf("agent help is not JSON: %v\n%s", err, stdout.String())
			}
			if document.SchemaVersion != agentHelpSchemaVersion || document.View != "scope" || document.Program != ProgramName ||
				document.Scope != (agentScope{Selector: spec.Path, Kind: "command"}) {
				t.Fatalf("agent document header = %+v", document)
			}
			if len(document.GlobalInputs) != 1 || document.GlobalInputs[0].Name != "--error-format" ||
				!reflect.DeepEqual(document.GlobalInputs[0].AllowedValues, []string{"text", "json"}) ||
				document.ErrorContract.CommandErrorsField != "commands[].contract.errors" || len(document.ErrorContract.ExitCodes) != 12 ||
				len(document.ErrorContract.GlobalErrors) != 6 || document.ErrorContract.JSONSchemaVersion != 1 {
				t.Fatalf("global agent contract = %+v / %+v", document.GlobalInputs, document.ErrorContract)
			}
			if document.IOContract.SuccessStream != "stdout" || document.IOContract.ErrorStream != "stderr" ||
				!document.IOContract.SuccessStatusRequiresCompleteWrite || document.IOContract.PartialOutputIsSuccess ||
				document.IOContract.ExternalTextTrust != "untrusted_data" ||
				document.IOContract.ExternalTextProjection != "visible_escape" ||
				document.IOContract.OpaqueReferencePolicy != "validated_exact_bytes" {
				t.Fatalf("I/O contract = %+v", document.IOContract)
			}
			if len(document.Commands) != 1 {
				t.Fatalf("selected commands = %+v", document.Commands)
			}
			got := document.Commands[0]
			if got.Path != spec.Path || got.Summary != spec.Summary || got.Usage != spec.Usage() || got.Args != spec.Args ||
				got.Effect != spec.Effect.String() || got.Role != spec.Role.String() ||
				!reflect.DeepEqual(got.Contract, spec.Agent) ||
				!reflect.DeepEqual(got.ProducesRefs, spec.ProducedRefs()) ||
				!reflect.DeepEqual(got.ConsumesRefs, spec.ConsumedRefs()) {
				t.Errorf("agent command = %+v, want catalog %+v", got, spec)
			}
			if got.Contract.Output.DefaultFormat == OutputFormatUnknown ||
				(containsOutputFormat(got.Contract.Output.Formats, OutputFormatJSON) && got.Contract.Output.JSONSchemaVersion <= 0) {
				t.Errorf("agent command %q has incomplete output metadata: %+v", got.Path, got.Contract.Output)
			}
		})
	}
}

func TestAgentHelpRootAndScopedShapeSnapshots(t *testing.T) {
	root := runAgentHelpForTest(t, []string{"help", "--format=agent"})
	assertJSONKeys(t, root, []string{"commands", "program", "schema_version", "scope_request", "view"})
	var rootCommands []map[string]json.RawMessage
	if err := json.Unmarshal(root["commands"], &rootCommands); err != nil {
		t.Fatal(err)
	}
	for index, command := range rootCommands {
		t.Run(fmt.Sprintf("root_command_%d", index), func(t *testing.T) {
			assertJSONKeys(t, command, []string{"capability_id", "effect", "namespace", "outcome", "path", "role", "summary"})
		})
	}
	var scopeRequest map[string]json.RawMessage
	if err := json.Unmarshal(root["scope_request"], &scopeRequest); err != nil {
		t.Fatal(err)
	}
	assertJSONKeys(t, scopeRequest, []string{"invocation_template", "known_path_max_invocations", "selector_fields", "unknown_outcome_max_invocations"})

	scoped := runAgentHelpForTest(t, []string{"help", "sample", "--format=agent"})
	assertJSONKeys(t, scoped, []string{"commands", "error_contract", "global_inputs", "io_contract", "program", "schema_version", "scope", "view", "workflows"})
	var ioContract map[string]json.RawMessage
	if err := json.Unmarshal(scoped["io_contract"], &ioContract); err != nil {
		t.Fatal(err)
	}
	assertJSONKeys(t, ioContract, []string{"error_stream", "external_text_projection", "external_text_trust", "opaque_reference_policy", "partial_output_is_success", "success_status_requires_complete_write", "success_stream"})
	var scopedCommands []map[string]json.RawMessage
	if err := json.Unmarshal(scoped["commands"], &scopedCommands); err != nil {
		t.Fatal(err)
	}
	for index, command := range scopedCommands {
		t.Run(fmt.Sprintf("scoped_command_%d", index), func(t *testing.T) {
			assertJSONKeys(t, command, []string{"args", "consumes_refs", "contract", "effect", "path", "produces_refs", "role", "summary", "usage"})
			if _, legacy := command["next_actions"]; legacy {
				t.Fatal("scoped agent help retained command-local reference next_actions")
			}
			var contract map[string]json.RawMessage
			if err := json.Unmarshal(command["contract"], &contract); err != nil {
				t.Fatal(err)
			}
			var output map[string]json.RawMessage
			if err := json.Unmarshal(contract["output"], &output); err != nil {
				t.Fatal(err)
			}
			assertJSONKeys(t, output, []string{
				"collection_coverage", "default_format", "delivery", "fields", "formats", "json_envelope", "json_schema_version",
			})
			if _, legacy := output["completeness"]; legacy {
				t.Fatal("scoped agent help retained the ambiguous output completeness field")
			}
		})
	}
	var workflows []map[string]json.RawMessage
	if err := json.Unmarshal(scoped["workflows"], &workflows); err != nil {
		t.Fatal(err)
	}
	if len(workflows) != 1 {
		t.Fatalf("workflows = %+v", workflows)
	}
	assertJSONKeys(t, workflows[0], []string{"consumers", "producers", "reference_kind"})
	var producers []map[string]json.RawMessage
	if err := json.Unmarshal(workflows[0]["producers"], &producers); err != nil {
		t.Fatal(err)
	}
	assertJSONKeys(t, producers[0], []string{"field", "path", "usage"})
	var consumers []map[string]json.RawMessage
	if err := json.Unmarshal(workflows[0]["consumers"], &consumers); err != nil {
		t.Fatal(err)
	}
	assertJSONKeys(t, consumers[0], []string{"input", "path", "usage"})
}

func TestRootAgentHelpSizeGrowthContainsOnlyIndexFields(t *testing.T) {
	command := New(strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	base := utilitySpec("base")
	makeCommands := func(count int) []CommandSpec {
		commands := make([]CommandSpec, 0, count)
		for index := 0; index < count; index++ {
			spec := cloneCommandSpec(base)
			spec.Path = fmt.Sprintf("area inspect%03d", index)
			spec.Summary = "Inspect one bounded synthetic area"
			spec.Agent.CapabilityID = fmt.Sprintf("area.inspect%03d", index)
			spec.Agent.Outcome = "Inspect one bounded synthetic area without external I/O"
			for errorIndex := range spec.Agent.Errors {
				for actionIndex := range spec.Agent.Errors[errorIndex].NextActions {
					spec.Agent.Errors[errorIndex].NextActions[actionIndex].Command = spec.Path
				}
			}
			commands = append(commands, spec)
		}
		return commands
	}
	one, err := command.renderAgentIndex(makeCommands(1))
	if err != nil {
		t.Fatal(err)
	}
	many, err := command.renderAgentIndex(makeCommands(101))
	if err != nil {
		t.Fatal(err)
	}
	perCommandGrowth := (len(many) - len(one)) / 100
	if perCommandGrowth > 320 {
		t.Fatalf("root index grew by %d bytes per command, want <= 320", perCommandGrowth)
	}
	catalog := NewCatalog(makeCommands(100)...)
	if err := catalog.Validate(); err != nil {
		t.Fatalf("100-command catalog failed validation: %v", err)
	}
	if selected, exact := catalog.Select("area"); exact || len(selected) != 100 {
		t.Fatalf("100-command namespace selection exact=%t, commands=%d", exact, len(selected))
	}
	if selected, exact := catalog.Select("area inspect042"); !exact || len(selected) != 1 || selected[0].Path != "area inspect042" {
		t.Fatalf("exact selection exact=%t, commands=%+v", exact, selected)
	}
	if selected, exact := catalog.Select("are"); exact || len(selected) != 0 {
		t.Fatalf("non-boundary selector exact=%t, commands=%+v", exact, selected)
	}
	for _, forbidden := range []string{"global_inputs", "io_contract", "error_contract", "workflows", "contract", "usage", "args", "inputs", "output", "errors", "mutation"} {
		if bytes.Contains(many, []byte(`"`+forbidden+`"`)) {
			t.Errorf("root index leaked detailed field %q", forbidden)
		}
	}

	oversized := cloneCommandSpec(base)
	oversized.Summary = strings.Repeat("s", maxAgentIndexEntryBytes)
	if err := NewCatalog(oversized).Validate(); err == nil || !strings.Contains(err.Error(), "agent index entry") {
		t.Fatalf("oversized root index entry error = %v", err)
	}
}

func TestCatalogSelectReturnsDeepCopiesForScopedProjection(t *testing.T) {
	catalog := DefaultCatalog()
	before := catalog.Commands()

	namespace, exact := catalog.Select("sample")
	if exact || len(namespace) != 2 {
		t.Fatalf("namespace selection exact=%t, commands=%+v", exact, namespace)
	}
	namespace[0].Agent.Inputs[0].AllowedValues[0] = "changed"
	namespace[0].Agent.Output.Fields[0].Name = "changed"
	namespace[0].Agent.Errors[0].NextActions[0].Command = "changed"

	selected, exact := catalog.Select("sample read")
	if !exact || len(selected) != 1 {
		t.Fatalf("exact selection exact=%t, commands=%+v", exact, selected)
	}
	selected[0].Agent.Inputs[0].ReferenceKind = "changed"
	selected[0].Agent.Output.Formats[0] = OutputFormatNone

	after := catalog.Commands()
	for index := range before {
		if before[index].Path != after[index].Path || before[index].Summary != after[index].Summary ||
			before[index].Args != after[index].Args || before[index].Effect != after[index].Effect ||
			before[index].Role != after[index].Role || !reflect.DeepEqual(before[index].Agent, after[index].Agent) {
			t.Fatalf("mutating scoped selections changed catalog command %q", before[index].Path)
		}
	}
}

func runAgentHelpForTest(t *testing.T, args []string) map[string]json.RawMessage {
	t.Helper()
	var stdout, stderr bytes.Buffer
	command := New(strings.NewReader(""), &stdout, &stderr)
	if code := runCLI(command, args); code != ExitOK {
		t.Fatalf("Run(%v) code = %d, stderr = %q", args, code, stderr.String())
	}
	var document map[string]json.RawMessage
	if err := json.Unmarshal(stdout.Bytes(), &document); err != nil {
		t.Fatalf("agent help is not JSON: %v\n%s", err, stdout.String())
	}
	return document
}

func assertJSONKeys(t *testing.T, document map[string]json.RawMessage, want []string) {
	t.Helper()
	got := make([]string, 0, len(document))
	for key := range document {
		got = append(got, key)
	}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("JSON keys = %v, want %v", got, want)
	}
}

func containsOutputFormat(formats []OutputFormat, wanted OutputFormat) bool {
	for _, format := range formats {
		if format == wanted {
			return true
		}
	}
	return false
}

func TestAgentHelpCanSelectNamespaceWithoutLoadingWholeCatalog(t *testing.T) {
	var stdout, stderr bytes.Buffer
	command := New(strings.NewReader(""), &stdout, &stderr)
	if code := runCLI(command, []string{"help", "sample", "--format=agent"}); code != ExitOK {
		t.Fatalf("Run(namespace agent help) code = %d, stderr = %q", code, stderr.String())
	}
	var document agentDocument
	if err := json.Unmarshal(stdout.Bytes(), &document); err != nil {
		t.Fatal(err)
	}
	if len(document.Commands) != 2 || document.Commands[0].Path != "sample list" || document.Commands[1].Path != "sample read" {
		t.Fatalf("namespace commands = %+v", document.Commands)
	}
	if len(document.Workflows) != 1 || document.Workflows[0].ReferenceKind != "sample" ||
		len(document.Workflows[0].Producers) != 1 || len(document.Workflows[0].Consumers) != 1 {
		t.Fatalf("namespace workflows = %+v", document.Workflows)
	}
	for _, entry := range document.Commands {
		if !strings.HasPrefix(entry.Path, "sample ") {
			t.Fatalf("unscoped command leaked into namespace help: %+v", entry)
		}
	}
}

func TestTextHelpCanSelectNamespace(t *testing.T) {
	var stdout, stderr bytes.Buffer
	command := New(strings.NewReader(""), &stdout, &stderr)
	if code := runCLI(command, []string{"help", "sample"}); code != ExitOK {
		t.Fatalf("Run(namespace help) code = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "sample list") || !strings.Contains(stdout.String(), "sample read") || strings.Contains(stdout.String(), "doctor") {
		t.Fatalf("namespace text = %q", stdout.String())
	}
}

func TestAgentHelpPreservesTopLevelCompatibilityFields(t *testing.T) {
	var stdout, stderr bytes.Buffer
	command := New(strings.NewReader(""), &stdout, &stderr)
	if code := runCLI(command, []string{"help", "sample", "list", "--format=agent"}); code != ExitOK {
		t.Fatalf("Run(selected agent help) code = %d, stderr = %q", code, stderr.String())
	}
	var raw struct {
		Commands []map[string]json.RawMessage `json:"commands"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"path", "summary", "usage", "effect", "role", "produces_refs", "consumes_refs"} {
		if _, exists := raw.Commands[0][field]; !exists {
			t.Errorf("scoped agent command lacks compatibility field %q", field)
		}
	}
	if _, exists := raw.Commands[0]["contract"]; !exists {
		t.Error("scoped agent command lacks structured contract")
	}
}

func TestAgentHelpCanSelectOneCatalogCommandWithItsWorkflow(t *testing.T) {
	var stdout, stderr bytes.Buffer
	command := New(strings.NewReader(""), &stdout, &stderr)
	if code := runCLI(command, []string{"help", "sample", "read", "--format=agent"}); code != ExitOK {
		t.Fatalf("Run(selected agent help) code = %d, stderr = %q", code, stderr.String())
	}
	var document agentDocument
	if err := json.Unmarshal(stdout.Bytes(), &document); err != nil {
		t.Fatal(err)
	}
	if len(document.Commands) != 1 || document.Commands[0].Path != "sample read" ||
		document.Commands[0].Effect != "read" || document.Commands[0].Role != "act" {
		t.Fatalf("commands = %+v", document.Commands)
	}
	if len(document.Workflows) != 1 || len(document.Workflows[0].Producers) != 1 ||
		document.Workflows[0].Producers[0].Path != "sample list" || len(document.Workflows[0].Consumers) != 1 ||
		document.Workflows[0].Consumers[0].Path != "sample read" {
		t.Fatalf("selected command workflows = %+v", document.Workflows)
	}
}

func TestAgentHelpPublishesDiscoverToActReferenceFlow(t *testing.T) {
	var stdout, stderr bytes.Buffer
	command := New(strings.NewReader(""), &stdout, &stderr)
	if code := runCLI(command, []string{"help", "sample", "--format", "agent"}); code != ExitOK {
		t.Fatalf("Run(agent help) code = %d, stderr = %q", code, stderr.String())
	}
	var document agentDocument
	if err := json.Unmarshal(stdout.Bytes(), &document); err != nil {
		t.Fatal(err)
	}
	commands := make(map[string]agentCommand, len(document.Commands))
	for _, entry := range document.Commands {
		commands[entry.Path] = entry
	}
	discover := commands["sample list"]
	if discover.Role != "discover" || discover.Effect != "read" ||
		!reflect.DeepEqual(discover.ProducesRefs, []ProducedRef{{Kind: "sample", Field: "id"}}) ||
		len(discover.ConsumesRefs) != 0 {
		t.Fatalf("sample list agent contract = %+v", discover)
	}
	act := commands["sample read"]
	if act.Role != "act" || act.Effect != "read" ||
		!reflect.DeepEqual(act.ConsumesRefs, []ConsumedRef{{Kind: "sample", Argument: "--id"}}) ||
		len(act.ProducesRefs) != 0 {
		t.Fatalf("sample read agent contract = %+v", act)
	}
	if len(document.Workflows) != 1 || document.Workflows[0].ReferenceKind != "sample" ||
		!reflect.DeepEqual(document.Workflows[0].Producers, []agentWorkflowProducer{{
			Path: "sample list", Usage: "agentic-cli-foundry sample list [--format tsv|json]", Field: "id",
		}}) || !reflect.DeepEqual(document.Workflows[0].Consumers, []agentWorkflowConsumer{{
		Path: "sample read", Usage: "agentic-cli-foundry sample read --id <sample-id> [--format tsv|json]", Input: "--id",
	}}) {
		t.Fatalf("derived grouped workflow = %+v", document.Workflows)
	}
}

func TestSelectedProducerDerivesExactNextArgvFromGroupedWorkflow(t *testing.T) {
	var listOut, listErr bytes.Buffer
	listCLI := New(strings.NewReader(""), &listOut, &listErr)
	if code := runCLI(listCLI, []string{"sample", "list", "--format", "json"}); code != ExitOK {
		t.Fatalf("sample list code = %d, stderr = %q", code, listErr.String())
	}
	var listed struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal(listOut.Bytes(), &listed); err != nil || len(listed.Items) == 0 {
		t.Fatalf("sample list output = %q, error = %v", listOut.String(), err)
	}

	var helpOut, helpErr bytes.Buffer
	helpCLI := New(strings.NewReader(""), &helpOut, &helpErr)
	if code := runCLI(helpCLI, []string{"help", "sample", "list", "--format=agent"}); code != ExitOK {
		t.Fatalf("selected producer help code = %d, stderr = %q", code, helpErr.String())
	}
	var document agentDocument
	if err := json.Unmarshal(helpOut.Bytes(), &document); err != nil {
		t.Fatal(err)
	}
	if len(document.Workflows) != 1 || len(document.Workflows[0].Producers) != 1 || len(document.Workflows[0].Consumers) != 1 {
		t.Fatalf("selected producer workflows = %+v", document.Workflows)
	}
	producer := document.Workflows[0].Producers[0]
	consumer := document.Workflows[0].Consumers[0]
	if producer.Path != "sample list" || producer.Field != "id" || consumer.Path != "sample read" ||
		consumer.Input != "--id" || consumer.Usage != "agentic-cli-foundry sample read --id <sample-id> [--format tsv|json]" {
		t.Fatalf("selected producer adjacency = producer %+v consumer %+v", producer, consumer)
	}
	nextArgv := append(strings.Fields(consumer.Path), consumer.Input, listed.Items[0].ID)
	var readOut, readErr bytes.Buffer
	readCLI := New(strings.NewReader(""), &readOut, &readErr)
	if code := runCLI(readCLI, nextArgv); code != ExitOK {
		t.Fatalf("derived next argv %v code = %d, stderr = %q", nextArgv, code, readErr.String())
	}
	if !strings.Contains(readOut.String(), listed.Items[0].ID) {
		t.Fatalf("derived next argv output = %q, want exact ID %q", readOut.String(), listed.Items[0].ID)
	}
}

func TestAgentRoundTripContractCoversDiscoveryActionAndRecovery(t *testing.T) {
	var stdout, stderr bytes.Buffer
	command := New(strings.NewReader(""), &stdout, &stderr)
	if code := runCLI(command, []string{"help", "sample", "--format=agent"}); code != ExitOK {
		t.Fatalf("Run() code = %d, stderr = %q", code, stderr.String())
	}
	var document agentDocument
	if err := json.Unmarshal(stdout.Bytes(), &document); err != nil {
		t.Fatal(err)
	}
	commands := make(map[string]agentCommand, len(document.Commands))
	for _, entry := range document.Commands {
		commands[entry.Path] = entry
	}
	discover := commands["sample list"]
	act := commands["sample read"]
	if discover.Contract.Output.Delivery != OutputDeliveryComplete ||
		discover.Contract.Output.CollectionCoverage != CollectionCoverageExhaustive ||
		len(discover.ProducesRefs) != 1 || discover.ProducesRefs[0] != (ProducedRef{Kind: "sample", Field: "id"}) {
		t.Fatalf("discovery contract = %+v", discover)
	}
	if len(act.Contract.Inputs) < 1 || act.Contract.Inputs[0].Name != "--id" ||
		act.Contract.Inputs[0].Source != InputSourceFlag || act.Contract.Inputs[0].ReferenceKind != "sample" ||
		act.Contract.Inputs[0].Description == "" || act.Contract.Inputs[0].AllowedValues == nil {
		t.Fatalf("action input contract = %+v", act.Contract.Inputs)
	}
	if len(document.Workflows) != 1 || len(document.Workflows[0].Producers) != 1 ||
		document.Workflows[0].Producers[0].Path != discover.Path || len(document.Workflows[0].Consumers) != 1 ||
		document.Workflows[0].Consumers[0].Path != act.Path || document.Workflows[0].Consumers[0].Input != "--id" {
		t.Fatalf("round-trip workflow = %+v", document.Workflows)
	}
	foundRecovery := false
	for _, declared := range act.Contract.Errors {
		if declared.Code == "sample_not_found" && declared.Kind == fault.KindNotFound &&
			len(declared.NextActions) == 1 && declared.NextActions[0].Command == discover.Path {
			foundRecovery = true
		}
	}
	if !foundRecovery {
		t.Fatalf("action errors lack discover recovery: %+v", act.Contract.Errors)
	}
}

type workflowEdge struct {
	ReferenceKind string
	Producer      agentWorkflowProducer
	Consumer      agentWorkflowConsumer
}

type legacyAgentWorkflow struct {
	ReferenceKind string                `json:"reference_kind"`
	Producer      agentWorkflowProducer `json:"producer"`
	Consumer      agentWorkflowConsumer `json:"consumer"`
}

func TestGroupedAgentWorkflowsPreserveEveryReferenceEdge(t *testing.T) {
	alphaList := discoverSpec("alpha list", "alpha")
	alphaSearch := discoverSpec("alpha search", "alpha")
	alphaRead := actSpec("alpha read", "alpha", "--left-id", "--right-id")
	betaList := discoverSpec("beta list", "beta")
	betaRead := actSpec("beta read", "beta", "--id")
	catalog := NewCatalog(alphaList, alphaSearch, alphaRead, betaList, betaRead)
	if err := catalog.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	workflows := catalog.referenceWorkflows()
	if len(workflows) != 2 {
		t.Fatalf("grouped workflows = %+v, want one record per reference kind", workflows)
	}
	got := groupedWorkflowEdges(workflows)
	want := pairExpandedWorkflowEdges(catalog.Commands())
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("grouped edges = %+v, want %+v", got, want)
	}
	if len(got) != 5 {
		t.Fatalf("edge count = %d, want 5", len(got))
	}

	selected, exact := catalog.Select("alpha read")
	if !exact || len(selected) != 1 {
		t.Fatalf("selected exact=%t commands=%+v", exact, selected)
	}
	scoped := workflowsForCommands(workflows, selected)
	if len(scoped) != 1 || scoped[0].ReferenceKind != "alpha" ||
		len(scoped[0].Producers) != 2 || len(scoped[0].Consumers) != 2 {
		t.Fatalf("scoped grouped workflow = %+v", scoped)
	}
}

func TestDerivedScaleScopedAgentHelpFitsWholeResponseBudget(t *testing.T) {
	catalog := derivedScaleHelpCatalog(t)
	selected, exact := catalog.Select("scale")
	if exact || len(selected) != 6 {
		t.Fatalf("scale selection exact=%t commands=%d", exact, len(selected))
	}

	encoded, err := (&CLI{catalog: catalog}).renderAgentHelp("scale", false, selected)
	if err != nil {
		t.Fatal(err)
	}
	const maxScopedHelpBytes = 64 * 1024
	if len(encoded) > maxScopedHelpBytes {
		t.Fatalf("grouped derived-scale scoped help = %d UTF-8 bytes, want <= %d", len(encoded), maxScopedHelpBytes)
	}

	var document agentDocument
	if err := json.Unmarshal(encoded, &document); err != nil {
		t.Fatal(err)
	}
	if document.SchemaVersion != 5 || len(document.Commands) != len(selected) || len(document.Workflows) != 1 ||
		len(document.Workflows[0].Producers) != 18 || len(document.Workflows[0].Consumers) != 18 {
		t.Fatalf("derived-scale grouped document = schema %d commands %d workflows %+v", document.SchemaVersion, len(document.Commands), document.Workflows)
	}
	if got, want := groupedWorkflowEdges(document.Workflows), pairExpandedWorkflowEdges(catalog.Commands()); !reflect.DeepEqual(got, want) {
		t.Fatalf("derived-scale grouped edges = %d, want %d", len(got), len(want))
	}

	legacyWorkflows := pairExpandGroupedWorkflows(document.Workflows)
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &raw); err != nil {
		t.Fatal(err)
	}
	raw["workflows"], err = json.Marshal(legacyWorkflows)
	if err != nil {
		t.Fatal(err)
	}
	legacyEncoded, err := json.Marshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	legacyEncoded = append(legacyEncoded, '\n')
	if len(legacyEncoded) <= maxScopedHelpBytes {
		t.Fatalf("synthetic corpus no longer exposes Cartesian growth: pair-expanded help = %d bytes, budget = %d", len(legacyEncoded), maxScopedHelpBytes)
	}
	t.Logf("derived-scale scoped help: grouped=%d bytes pair-expanded=%d bytes budget=%d bytes edges=%d",
		len(encoded), len(legacyEncoded), maxScopedHelpBytes, len(legacyWorkflows))
}

func derivedScaleHelpCatalog(t *testing.T) Catalog {
	t.Helper()
	const commandsPerRole = 3
	const endpointsPerCommand = 6
	commands := make([]CommandSpec, 0, commandsPerRole*2)
	for commandIndex := 0; commandIndex < commandsPerRole; commandIndex++ {
		spec := discoverSpec(fmt.Sprintf("scale discover%02d", commandIndex), "resource")
		spec.Agent.Output.Fields = make([]OutputField, 0, endpointsPerCommand)
		for endpointIndex := 0; endpointIndex < endpointsPerCommand; endpointIndex++ {
			spec.Agent.Output.Fields = append(spec.Agent.Output.Fields, OutputField{
				Name:          fmt.Sprintf("resource_%02d_%02d", commandIndex, endpointIndex),
				Type:          OutputFieldTypeString,
				Description:   "Opaque synthetic resource reference.",
				ReferenceKind: "resource",
			})
		}
		commands = append(commands, spec)
	}
	for commandIndex := 0; commandIndex < commandsPerRole; commandIndex++ {
		inputs := make([]string, 0, endpointsPerCommand)
		for endpointIndex := 0; endpointIndex < endpointsPerCommand; endpointIndex++ {
			inputs = append(inputs, fmt.Sprintf("--resource-%02d-%02d", commandIndex, endpointIndex))
		}
		commands = append(commands, actSpec(fmt.Sprintf("scale inspect%02d", commandIndex), "resource", inputs...))
	}
	catalog := NewCatalog(commands...)
	if err := catalog.Validate(); err != nil {
		t.Fatalf("derived-scale catalog validation: %v", err)
	}
	return catalog
}

func pairExpandedWorkflowEdges(commands []CommandSpec) map[workflowEdge]struct{} {
	edges := make(map[workflowEdge]struct{})
	for _, producerCommand := range commands {
		for _, produced := range producerCommand.ProducedRefs() {
			for _, consumerCommand := range commands {
				for _, consumed := range consumerCommand.ConsumedRefs() {
					if produced.Kind != consumed.Kind {
						continue
					}
					edges[workflowEdge{
						ReferenceKind: produced.Kind,
						Producer:      agentWorkflowProducer{Path: producerCommand.Path, Usage: producerCommand.Usage(), Field: produced.Field},
						Consumer:      agentWorkflowConsumer{Path: consumerCommand.Path, Usage: consumerCommand.Usage(), Input: consumed.Argument},
					}] = struct{}{}
				}
			}
		}
	}
	return edges
}

func groupedWorkflowEdges(workflows []agentWorkflow) map[workflowEdge]struct{} {
	edges := make(map[workflowEdge]struct{})
	for _, workflow := range workflows {
		for _, producer := range workflow.Producers {
			for _, consumer := range workflow.Consumers {
				edges[workflowEdge{ReferenceKind: workflow.ReferenceKind, Producer: producer, Consumer: consumer}] = struct{}{}
			}
		}
	}
	return edges
}

func pairExpandGroupedWorkflows(workflows []agentWorkflow) []legacyAgentWorkflow {
	expanded := make([]legacyAgentWorkflow, 0)
	for _, workflow := range workflows {
		for _, producer := range workflow.Producers {
			for _, consumer := range workflow.Consumers {
				expanded = append(expanded, legacyAgentWorkflow{
					ReferenceKind: workflow.ReferenceKind,
					Producer:      producer,
					Consumer:      consumer,
				})
			}
		}
	}
	return expanded
}

func TestHelpRejectsUnknownSelectorsAndFormats(t *testing.T) {
	for _, args := range [][]string{
		{"help", "missing"},
		{"help", "--format", "yaml"},
		{"help", "--unknown"},
	} {
		var stdout, stderr bytes.Buffer
		command := New(strings.NewReader(""), &stdout, &stderr)
		if code := runCLI(command, args); code != ExitUsage {
			t.Errorf("Run(%v) code = %d, want %d", args, code, ExitUsage)
		}
		if stdout.Len() != 0 || !strings.Contains(stderr.String(), "error:") {
			t.Errorf("Run(%v) stdout = %q, stderr = %q", args, stdout.String(), stderr.String())
		}
	}
}
