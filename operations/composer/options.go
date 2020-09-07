package composer

type Option func(*options)

type options struct {
	composeFile string
	projectDir  string
	envFile     string
}

func WithComposeFile(composeFile string) Option {
	return func(opts *options) {
		opts.composeFile = composeFile
	}
}

func WithProjectDir(projectDir string) Option {
	return func(opts *options) {
		opts.projectDir = projectDir
	}
}

func WithEnvFile(envFile string) Option {
	return func(opts *options) {
		opts.envFile = envFile
	}
}
