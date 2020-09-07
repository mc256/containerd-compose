package composer

import (
	"testing"
)

func Test_ComposerLoaderProvideComposeFile(t *testing.T) {
	var opts []Option

	//opts = append(opts, WithComposeFile("./test.yml"))
	opts = append(opts, WithComposeFile("./docker-compose.yml"))
	var compose *ComposeFile
	var err error
	if compose, err = loadFile(&opts); err != nil {
		t.Error(err)
	}
	if compose.Version != "2" {
		t.Error("version error")
	}
}

func Test_ComposerLoaderNonProvideComposeFile(t *testing.T) {
	var opts []Option

	opts = append(opts)
	opts = append(opts, WithComposeFile("./docker-compose.yml"))
	var compose *ComposeFile
	var err error
	if compose, err = loadFile(&opts); err != nil {
		t.Error(err)
	}
	if compose.Version != "2" {
		t.Error("version error")
	}
}
