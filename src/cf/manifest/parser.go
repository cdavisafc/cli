package manifest

import (
	"generic"
	"io"
	"fmt"
	"github.com/kylelemons/go-gypsy/yaml"
)

func Parse(reader io.Reader) (manifest *Manifest, errs ManifestErrors) {
	yamlMap := generic.NewMap()

	// gypsy
	node, err := yaml.Parse(reader)
	if err != nil {
		errs = append(errs, err)
		return
	}
	fmt.Printf("\n\n%#v\n", node)

	manifest, errs = NewManifest(yamlMap)
	return
}
