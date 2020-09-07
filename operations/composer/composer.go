package composer

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func init() {

}

func loadFile(opts *[]Option) (*ComposeFile, error) {
	// Options
	opt := options{}
	for _, o := range *opts {
		o(&opt)
	}

	// Read yaml file from local
	var bytes []byte
	var err error
	if opt.composeFile != "" {
		bytes, err = ioutil.ReadFile(opt.composeFile)
		if err != nil {
			return nil, err
		}
	} else {
		// try default value
		for _, d := range []string{
			"./containerd-compose.yml",
			"./docker-compose.yml",
		} {
			bytes, err = ioutil.ReadFile(d)
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, err
		}
	}

	// Parse Yaml file
	t := ComposeFile{}
	if err := yaml.Unmarshal(bytes, &t); err != nil {
		return nil, err
	}

	return &t, nil
}
