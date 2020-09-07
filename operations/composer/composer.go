package composer

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
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
	var buffer []byte
	var err error
	if opt.composeFile != "" {
		buffer, err = ioutil.ReadFile(opt.composeFile)
		if err != nil {
			return nil, err
		}
	} else {
		// try default value
		for _, d := range []string{
			"./containerd-compose.yml",
			"./docker-compose.yml",
		} {
			buffer, err = ioutil.ReadFile(d)
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, err
		}
	}

	// Get Environment Variables
	if opt.envFile != "" {
		_ = godotenv.Load(opt.envFile)
	} else {
		_ = godotenv.Load(".env")
	}
	buffer = []byte(os.ExpandEnv(string(buffer)))

	// Parse Yaml file
	t := ComposeFile{}
	if err := yaml.Unmarshal(buffer, &t); err != nil {
		return nil, err
	}

	return &t, nil
}
