package compose

import (
	"path/filepath"

	"github.com/openshift/origin/third_party/github.com/docker/libcompose/project"
)

func Load(paths ...string) (*project.Project, error) {
	for i := range paths {
		path, err := filepath.Abs(paths[i])
		if err != nil {
			return nil, err
		}
		paths[i] = path
	}
	var bases []string
	for _, s := range paths {
		bases = append(bases, filepath.Dir(s))
	}

	context := &project.Context{
		ComposeFiles: paths,
	}
	p := project.NewProject(context)
	if err := p.Parse(); err != nil {
		return nil, err
	}

	return p, err
}
