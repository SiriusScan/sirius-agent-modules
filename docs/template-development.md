# Template Development Guide

Sirius templates are YAML files under `templates/`. Use a unique top-level `id`,
an `info` object, and one or more `detection.steps`. Both `.yaml` and `.yml`
are supported in nested directories.

## Template structure

```yaml
id: EXAMPLE-001
info:
  name: Example configuration flag
  severity: info
  author: Sirius contributors
  description: Looks for an example flag in a fixture configuration file.
  remediation: Remove the example flag from the test fixture.

detection:
  logic: all
  steps:
    - type: file_content
      platforms: [linux, darwin]
      weight: 1
      config:
        path: ./example.ini
        regex: '^example_flag=true$'
```

Each `detection.steps[].type` names a built-in detection module. The
agent supports `file_hash`, `file_content`, `file_search`, and `version_cmd`.
See the [agent template architect guide](https://github.com/SiriusScan/app-agent/blob/main/documentation/README.template-architect-guide.md)
for module-specific `config` fields.

- `info.name` and `info.severity` are required. Severities are `critical`,
  `high`, `medium`, `low`, and `info`.
- `detection.logic` defaults to `all`. Use `any` when an applicable step can
  establish the result on its own.
- `detection.steps[].platforms` accepts `darwin`, `linux`, and `windows`. An
  omitted or empty list makes the step applicable to every supported platform.
- `detection.steps[].weight` is between 0 and 1. The parser defaults an omitted
  weight to 1.
- Use `info.description` to explain the result and `info.remediation` to
  describe what an operator can do. Add `info.references`, `info.cve`, or
  `info.cwe` when they apply.

## Regenerate and validate the manifests

1. Add the file under `templates/` and stage the path with Git. Add an entry
   under `templates` in `templates/manifest.json` keyed by its path relative to
   `templates/`. Include relevant descriptive metadata such as `author` and
   `version`. Remove the manifest entry when deleting a template.
2. Check out `SiriusScan/app-agent` at the commit pinned in
   [the validation workflow](../.github/workflows/validate-templates.yml).
   Use the Go version specified by that checkout's `go.mod`.
3. From the root of this module checkout, regenerate and validate with the
   following commands.

```bash
module_checkout="$PWD"
app_agent_checkout="/absolute/path/to/app-agent"

go -C "$app_agent_checkout" run "$module_checkout/tools/rebuild-manifests.go" "$module_checkout"
go -C "$app_agent_checkout" run ./cmd/template-cli validate-repository "$module_checkout"
```

The generator takes inventory from Git's index and reads the exact checkout
bytes. It fills in template IDs, sizes, SHA-256 checksums, detection metadata,
and repository totals. It also checks tracked scripts against their manifest
and updates their sizes and checksums. Stage all changed manifests after
reviewing the diff. LF line endings keep the digests consistent across
platforms.

The template manifest uses `format_version: 2`. Every entry needs a
`sha256:<64 lowercase hex characters>` digest, byte size, and an ID that
matches the YAML file. CI runs the pinned offline validator and checks that
regenerating the manifests leaves them unchanged. Both checks treat templates
and scripts as data.
