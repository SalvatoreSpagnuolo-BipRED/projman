package project

import (
	"os"
	"path/filepath"
)

type Project struct {
	Name string
	Path string
}

func Discover(root string) ([]Project, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var list []Project
	for _, e := range entries {
		if e.IsDir() {
			p := filepath.Join(root, e.Name())
			if _, err := os.Stat(filepath.Join(p, "pom.xml")); err == nil {
				list = append(list, Project{Name: e.Name(), Path: p})
			}
		}
	}
	return list, nil
}

func Names(projs []Project) []string {
	names := make([]string, len(projs))
	for i, p := range projs {
		names[i] = p.Name
	}
	return names
}

func Filter(all []Project, selected []string) []Project {
	if len(selected) == 0 {
		return all
	}
	m := map[string]bool{}
	for _, s := range selected {
		m[s] = true
	}
	var out []Project
	for _, p := range all {
		if m[p.Name] {
			out = append(out, p)
		}
	}
	return out
}
