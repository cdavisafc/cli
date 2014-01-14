package manifest

import (
	"generic"
	"github.com/cloudfoundry/gamble"
	"io"
	"io/ioutil"
)

func Parse(reader io.Reader) (manifest *Manifest, errs ManifestErrors) {
	yamlBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		errs = append(errs, err)
		return
	}

	document, err := gamble.Parse(string(yamlBytes))
	if err != nil {
		errs = append(errs, err)
		return
	}

	yamlMap := generic.NewMap(document)
	manifest, errs = NewManifest(yamlMap)
	return
}
