package model

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetPackageConfig_Empty(t *testing.T) {
	cm := ConfigMap{}
	cfg, err := cm.GetPackageConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.PackageSources) != 0 {
		t.Fatalf("expected 0 sources, got %d", len(cfg.PackageSources))
	}
}

func TestGetPackageConfig_ValidSources(t *testing.T) {
	cm := ConfigMap{
		"packageSources": []any{
			map[string]any{"name": "pacman", "command": "sudo pacman -S {package}"},
			map[string]any{"name": "npm", "command": "npm install -g {package}"},
		},
	}
	cfg, err := cm.GetPackageConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.PackageSources) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(cfg.PackageSources))
	}
	if cfg.PackageSources[0].Name != "pacman" || cfg.PackageSources[0].Command != "sudo pacman -S {package}" {
		t.Errorf("unexpected first source: %+v", cfg.PackageSources[0])
	}
	if cfg.PackageSources[1].Name != "npm" || cfg.PackageSources[1].Command != "npm install -g {package}" {
		t.Errorf("unexpected second source: %+v", cfg.PackageSources[1])
	}
}

func TestGetPackageConfig_MissingFields(t *testing.T) {
	cm := ConfigMap{
		"packageSources": []any{
			map[string]any{"name": "pacman"},
			map[string]any{"command": "npm install -g {package}"},
		},
	}
	cfg, err := cm.GetPackageConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.PackageSources) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(cfg.PackageSources))
	}
	if cfg.PackageSources[0].Name != "pacman" || cfg.PackageSources[0].Command != "" {
		t.Errorf("unexpected source: %+v", cfg.PackageSources[0])
	}
	if cfg.PackageSources[1].Name != "" || cfg.PackageSources[1].Command != "npm install -g {package}" {
		t.Errorf("unexpected source: %+v", cfg.PackageSources[1])
	}
}

func TestGetPackageConfig_WrongType(t *testing.T) {
	cm := ConfigMap{
		"packageSources": "not a slice",
	}
	cfg, err := cm.GetPackageConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.PackageSources) != 0 {
		t.Fatalf("expected 0 sources, got %d", len(cfg.PackageSources))
	}
}

func TestGetPackageConfig_NonMapEntries(t *testing.T) {
	cm := ConfigMap{
		"packageSources": []any{"not a map", 42, nil},
	}
	cfg, err := cm.GetPackageConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.PackageSources) != 0 {
		t.Fatalf("expected 0 sources (non-map entries skipped), got %d", len(cfg.PackageSources))
	}
}

func TestMergeConfigs_DstOnly(t *testing.T) {
	dst := ConfigMap{"a": "1", "b": "2"}
	src := ConfigMap{}
	result := mergeConfigs(dst, src)
	if result["a"] != "1" || result["b"] != "2" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestMergeConfigs_SrcOnly(t *testing.T) {
	dst := ConfigMap{}
	src := ConfigMap{"a": "1", "b": "2"}
	result := mergeConfigs(dst, src)
	if result["a"] != "1" || result["b"] != "2" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestMergeConfigs_ScalarOverride(t *testing.T) {
	dst := ConfigMap{"a": "old", "b": "keep"}
	src := ConfigMap{"a": "new"}
	result := mergeConfigs(dst, src)
	if result["a"] != "new" {
		t.Errorf("expected src to win for scalar, got %v", result["a"])
	}
	if result["b"] != "keep" {
		t.Errorf("expected dst key to remain, got %v", result["b"])
	}
}

func TestMergeConfigs_NestedMaps(t *testing.T) {
	dst := ConfigMap{"nested": ConfigMap{"x": "1", "y": "2"}}
	src := ConfigMap{"nested": ConfigMap{"y": "overridden", "z": "3"}}
	result := mergeConfigs(dst, src)
	nested, ok := result["nested"].(ConfigMap)
	if !ok {
		t.Fatalf("expected nested to be ConfigMap, got %T", result["nested"])
	}
	if nested["x"] != "1" {
		t.Errorf("expected nested.x=1, got %v", nested["x"])
	}
	if nested["y"] != "overridden" {
		t.Errorf("expected nested.y=overridden, got %v", nested["y"])
	}
	if nested["z"] != "3" {
		t.Errorf("expected nested.z=3, got %v", nested["z"])
	}
}

func TestMergeConfigs_DoesNotMutateDst(t *testing.T) {
	dst := ConfigMap{"a": "1"}
	src := ConfigMap{"a": "2"}
	_ = mergeConfigs(dst, src)
	if dst["a"] != "1" {
		t.Errorf("mergeConfigs mutated dst")
	}
}

func TestMergeValues_BothMaps(t *testing.T) {
	dst := ConfigMap{"a": "1"}
	src := ConfigMap{"b": "2"}
	result := mergeValues(dst, src)
	m, ok := result.(ConfigMap)
	if !ok {
		t.Fatalf("expected ConfigMap result, got %T", result)
	}
	if m["a"] != "1" || m["b"] != "2" {
		t.Errorf("unexpected merge result: %v", m)
	}
}

func TestMergeValues_BothSlices(t *testing.T) {
	dst := []any{"a", "b"}
	src := []any{"b", "c"}
	result := mergeValues(dst, src)
	s, ok := result.([]any)
	if !ok {
		t.Fatalf("expected []any result, got %T", result)
	}
	if len(s) != 3 {
		t.Errorf("expected 3 elements (a,b,c), got %d: %v", len(s), s)
	}
}

func TestMergeValues_ScalarWins(t *testing.T) {
	result := mergeValues("old", "new")
	if result != "new" {
		t.Errorf("expected src scalar to win, got %v", result)
	}
}

func TestMergeValues_MapVsScalar(t *testing.T) {
	dst := ConfigMap{"a": "1"}
	src := "scalar"
	result := mergeValues(dst, src)
	if result != "scalar" {
		t.Errorf("expected scalar src to win over map dst, got %v", result)
	}
}

func TestMergeValues_SliceVsMap(t *testing.T) {
	dst := []any{"a"}
	src := ConfigMap{"b": "2"}
	result := mergeValues(dst, src)
	if _, ok := result.(ConfigMap); !ok {
		t.Errorf("expected ConfigMap src to win when types differ, got %T", result)
	}
}

func TestMergeSlices_NoDuplicates(t *testing.T) {
	a := []any{"x", "y"}
	b := []any{"z"}
	result := mergeSlices(a, b)
	if len(result) != 3 {
		t.Errorf("expected 3 elements, got %d: %v", len(result), result)
	}
}

func TestMergeSlices_Deduplication(t *testing.T) {
	a := []any{"x", "y"}
	b := []any{"y", "z"}
	result := mergeSlices(a, b)
	if len(result) != 3 {
		t.Errorf("expected 3 elements after dedup, got %d: %v", len(result), result)
	}
	seen := map[any]int{}
	for _, v := range result {
		seen[v]++
	}
	if seen["y"] != 1 {
		t.Errorf("expected y to appear once, appeared %d times", seen["y"])
	}
}

func TestMergeSlices_NonComparableAlwaysAppended(t *testing.T) {
	a := []any{ConfigMap{"k": "v"}}
	b := []any{ConfigMap{"k": "v"}}
	result := mergeSlices(a, b)
	if len(result) != 2 {
		t.Errorf("expected 2 non-comparable elements (no dedup), got %d", len(result))
	}
}

func TestMergeSlices_NestedSlicesNotDeduped(t *testing.T) {
	a := []any{[]any{"a", "b"}}
	b := []any{[]any{"a", "b"}}
	result := mergeSlices(a, b)
	if len(result) != 2 {
		t.Errorf("expected 2 (nested slices not comparable), got %d", len(result))
	}
}

func TestMergeSlices_Order(t *testing.T) {
	a := []any{"a", "b"}
	b := []any{"c", "d"}
	result := mergeSlices(a, b)
	expected := []any{"a", "b", "c", "d"}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("order mismatch at index %d: expected %v, got %v", i, v, result[i])
		}
	}
}

func TestIsComparable_ConfigMap(t *testing.T) {
	if isComparable(ConfigMap{}) {
		t.Error("ConfigMap should not be comparable")
	}
}

func TestIsComparable_Slice(t *testing.T) {
	if isComparable([]any{}) {
		t.Error("[]any should not be comparable")
	}
}

func TestIsComparable_Primitives(t *testing.T) {
	for _, v := range []any{"str", 42, true, 3.14, nil} {
		if !isComparable(v) {
			t.Errorf("expected %T(%v) to be comparable", v, v)
		}
	}
}

func TestLoadConfigFiles_NoPaths(t *testing.T) {
	_, err := LoadConfigFiles(nil)
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestLoadConfigFiles_SingleYAML(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "config.yaml")
	os.WriteFile(f, []byte("key: value\n"), 0600)

	result, err := LoadConfigFiles([]string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("expected key=value, got %v", result["key"])
	}
}

func TestLoadConfigFiles_SingleJSON(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "config.json")
	os.WriteFile(f, []byte(`{"key":"value"}`), 0600)

	result, err := LoadConfigFiles([]string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("expected key=value, got %v", result["key"])
	}
}

func TestLoadConfigFiles_MergesMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "base.yaml")
	f2 := filepath.Join(dir, "override.yaml")
	os.WriteFile(f1, []byte("a: base\nb: keep\n"), 0600)
	os.WriteFile(f2, []byte("a: override\nc: new\n"), 0600)

	result, err := LoadConfigFiles([]string{f1, f2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["a"] != "override" {
		t.Errorf("expected a=override, got %v", result["a"])
	}
	if result["b"] != "keep" {
		t.Errorf("expected b=keep, got %v", result["b"])
	}
	if result["c"] != "new" {
		t.Errorf("expected c=new, got %v", result["c"])
	}
}

func TestLoadConfigFiles_NonexistentFile(t *testing.T) {
	_, err := LoadConfigFiles([]string{"/nonexistent/path/config.yaml"})
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLoadConfigFiles_UnsupportedExtension(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "config.toml")
	os.WriteFile(f, []byte("key = \"value\"\n"), 0600)

	_, err := LoadConfigFiles([]string{f})
	if err == nil {
		t.Fatal("expected error for unsupported extension")
	}
}

func TestLoadConfigFiles_SecondFileError(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "base.yaml")
	os.WriteFile(f1, []byte("a: 1\n"), 0600)

	_, err := LoadConfigFiles([]string{f1, "/nonexistent/config.yaml"})
	if err == nil {
		t.Fatal("expected error when second file is missing")
	}
}

func TestLoadConfigFiles_YMLExtension(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "config.yml")
	os.WriteFile(f, []byte("foo: bar\n"), 0600)

	result, err := LoadConfigFiles([]string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %v", result["foo"])
	}
}
