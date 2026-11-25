// Rebuild the published module manifests from staged paths and checkout bytes.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	repositoryManifest = "repository-manifest.json"
	templateManifest   = "templates/manifest.json"
	scriptManifest     = "scripts/manifest.json"
)

type templateHeader struct {
	ID   string `yaml:"id"`
	Info struct {
		Severity string `yaml:"severity"`
	} `yaml:"info"`
	Detection struct {
		Steps []struct {
			Type      string   `yaml:"type"`
			Platforms []string `yaml:"platforms"`
		} `yaml:"steps"`
	} `yaml:"detection"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run tools/rebuild-manifests.go <module-checkout>")
		os.Exit(1)
	}
	if err := rebuild(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rebuild(directory string) error {
	root, err := os.OpenRoot(directory)
	if err != nil {
		return err
	}
	defer root.Close()

	manifests := make(map[string]map[string]any)
	names := []string{repositoryManifest, templateManifest, scriptManifest}
	for _, name := range names {
		data, err := readFile(root, name)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		var manifest map[string]any
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber()
		if !json.Valid(data) || decoder.Decode(&manifest) != nil || manifest == nil {
			return fmt.Errorf("%s: expected a JSON object", name)
		}
		manifests[name] = manifest
	}

	cmd := exec.Command("git", "ls-files", "-z", "--", "templates", "scripts")
	cmd.Dir = directory
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("list tracked module files: %w", err)
	}

	templateStats, err := updateInventory(root, strings.Split(string(output), "\x00"), manifests[templateManifest], "templates")
	if err != nil {
		return err
	}
	scriptStats, err := updateInventory(root, strings.Split(string(output), "\x00"), manifests[scriptManifest], "scripts")
	if err != nil {
		return err
	}

	manifests[templateManifest]["format_version"] = 2
	manifests[templateManifest]["statistics"] = map[string]any{
		"total_templates": templateStats.total,
		"by_type":         templateStats.byKind,
		"by_severity":     templateStats.bySeverity,
		"by_platform":     templateStats.byPlatform,
	}
	manifests[scriptManifest]["statistics"] = map[string]any{
		"total_scripts": scriptStats.total,
		"by_language":   scriptStats.byKind,
		"by_severity":   scriptStats.bySeverity,
		"by_platform":   scriptStats.byPlatform,
	}
	platforms := make(map[string]int)
	for platform, count := range templateStats.byPlatform {
		if platform == "darwin" {
			platform = "macos"
		}
		platforms[platform] += count
	}
	for platform, count := range scriptStats.byPlatform {
		if platform != "cross-platform" {
			platforms[platform] += count
		}
	}
	manifests[repositoryManifest]["statistics"] = map[string]any{
		"total_templates": templateStats.total,
		"total_scripts":   scriptStats.total,
		"by_type":         templateStats.byKind,
		"by_platform":     platforms,
	}

	// Keep the root component references aligned with the declared component releases.
	components, ok := manifests[repositoryManifest]["components"].(map[string]any)
	if !ok {
		return fmt.Errorf("%s: components must be an object", repositoryManifest)
	}
	for _, kind := range []string{"templates", "scripts"} {
		component, ok := components[kind].(map[string]any)
		if !ok {
			return fmt.Errorf("%s: missing components.%s object", repositoryManifest, kind)
		}
		manifest := manifests[kind+"/manifest.json"]
		component["version"] = manifest["version"]
		component["updated"] = manifest["updated"]
	}

	for _, name := range names {
		var data bytes.Buffer
		encoder := json.NewEncoder(&data)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(manifests[name]); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		if err := os.WriteFile(filepath.Join(directory, filepath.FromSlash(name)), data.Bytes(), 0o644); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

type inventoryStats struct {
	total      int
	byKind     map[string]int
	bySeverity map[string]int
	byPlatform map[string]int
}

func updateInventory(root *os.Root, files []string, manifest map[string]any, kind string) (inventoryStats, error) {
	stats := inventoryStats{
		byKind:     make(map[string]int),
		bySeverity: make(map[string]int),
		byPlatform: make(map[string]int),
	}
	manifestPath := kind + "/manifest.json"
	entries, ok := manifest[kind].(map[string]any)
	if !ok {
		return stats, fmt.Errorf("%s: %s must be an object", manifestPath, kind)
	}

	seenFiles := make(map[string]bool)
	seenIDs := make(map[string]string)
	var issues []error
	for _, name := range files {
		relative, underComponent := strings.CutPrefix(name, kind+"/")
		if !underComponent || name == manifestPath {
			continue
		}
		extension := strings.ToLower(path.Ext(relative))
		if kind == "templates" && extension != ".yaml" && extension != ".yml" {
			continue
		}
		if kind == "scripts" && extension != ".sh" && extension != ".ps1" && extension != ".py" {
			continue
		}
		seenFiles[relative] = true

		entry, ok := entries[relative].(map[string]any)
		if !ok {
			issues = append(issues, fmt.Errorf("%s: add a metadata object in %s", name, manifestPath))
			continue
		}
		content, err := readFile(root, name)
		if err != nil {
			issues = append(issues, fmt.Errorf("%s: %w", name, err))
			continue
		}
		if bytes.Contains(content, []byte("\r\n")) {
			issues = append(issues, fmt.Errorf("%s: use LF line endings before regenerating", name))
			continue
		}
		entry["size"] = len(content)
		entry["checksum"] = fmt.Sprintf("sha256:%x", sha256.Sum256(content))

		var platformNames []string
		category, _ := entry["language"].(string)
		if kind == "templates" {
			header, platforms, err := readTemplateHeader(content)
			if err != nil {
				issues = append(issues, fmt.Errorf("%s: %w", name, err))
				continue
			}
			platformNames = platforms
			entry["id"] = header.ID
			if header.Info.Severity != "" {
				entry["severity"] = header.Info.Severity
			}
			category = strings.ReplaceAll(header.Detection.Steps[0].Type, "_", "-")
			entry["detection_type"] = category
			entry["platforms"] = platforms
		} else {
			platforms, ok := entry["platforms"].([]any)
			if !ok || len(platforms) == 0 {
				issues = append(issues, fmt.Errorf("%s: script metadata requires a platforms array", name))
				continue
			}
			for _, value := range platforms {
				platform, ok := value.(string)
				if !ok || strings.TrimSpace(platform) == "" {
					issues = append(issues, fmt.Errorf("%s: platform names must be nonempty strings", name))
					continue
				}
				platformNames = append(platformNames, platform)
			}
			if strings.HasPrefix(relative, "cross-platform/") {
				platformNames = append(platformNames, "cross-platform")
			}
		}

		id, _ := entry["id"].(string)
		severity, _ := entry["severity"].(string)
		if strings.TrimSpace(id) == "" || strings.TrimSpace(severity) == "" || strings.TrimSpace(category) == "" {
			issues = append(issues, fmt.Errorf("%s: metadata requires an ID, severity, and detection type or language", name))
			continue
		}
		if previous, ok := seenIDs[id]; ok {
			issues = append(issues, fmt.Errorf("%s: duplicate ID %q also used by %s", name, id, previous))
			continue
		}
		seenIDs[id] = name

		stats.total++
		stats.byKind[category]++
		stats.bySeverity[severity]++
		slices.Sort(platformNames)
		for _, platform := range slices.Compact(platformNames) {
			stats.byPlatform[platform]++
		}
	}
	for _, name := range slices.Sorted(maps.Keys(entries)) {
		if !seenFiles[name] {
			issues = append(issues, fmt.Errorf("%s/%s: listed in %s but absent from tracked module files", kind, name, manifestPath))
		}
	}
	return stats, errors.Join(issues...)
}

func readTemplateHeader(content []byte) (templateHeader, []string, error) {
	var header templateHeader
	if err := yaml.Unmarshal(content, &header); err != nil {
		return header, nil, err
	}
	if strings.TrimSpace(header.ID) == "" || len(header.Detection.Steps) == 0 {
		return header, nil, fmt.Errorf("template metadata requires a YAML ID and detection.steps")
	}
	platforms := make(map[string]bool)
	for _, step := range header.Detection.Steps {
		if len(step.Platforms) == 0 {
			return header, []string{"darwin", "linux", "windows"}, nil
		}
		for _, platform := range step.Platforms {
			platforms[platform] = true
		}
	}
	return header, slices.Sorted(maps.Keys(platforms)), nil
}

func readFile(root *os.Root, name string) ([]byte, error) {
	localName, err := filepath.Localize(name)
	if err != nil {
		return nil, err
	}
	info, err := root.Lstat(localName)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("must be a regular file")
	}
	file, err := root.Open(localName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}
