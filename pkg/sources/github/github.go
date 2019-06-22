package github

type GithubReceiver struct {
	Secret string
}

func (g *GithubReceiver) HandleEvent() error {
	return nil
}

func (g *GithubReceiver) ValidateType() bool {
	// validate based on existance of GH header.
	return true
}
