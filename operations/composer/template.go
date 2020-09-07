package composer

type ComposeFile struct {
	Version  string
	Services map[string]struct {
		Image       string
		Volumes     []string
		Environment []string
	}
}
