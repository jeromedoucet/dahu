package model

type HttpAuthConfig struct {
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (a HttpAuthConfig) IsValid() bool {
	if a.Url == "" {
		return false
	}
	return true
}

func (a *HttpAuthConfig) ToPublicModel() {
	a.User = ""
	a.Password = ""
}

type SshAuthConfig struct {
	Url         string `json:"url"`
	Key         string `json:"key"`
	KeyPassword string `json:"keyPassword"`
}

func (a SshAuthConfig) IsValid() bool {
	if a.Url == "" {
		return false
	} else if a.Key == "" {
		return false
	}
	return true
}

func (a *SshAuthConfig) ToPublicModel() {
	a.Key = ""
	a.KeyPassword = ""
}

type GitConfig struct {
	HttpAuth *HttpAuthConfig `json:"httpAuth"`
	SshAuth  *SshAuthConfig  `json:"sshAuth"`
}

// todo test me
func (g GitConfig) IsValid() bool {
	if (g.HttpAuth == nil && g.SshAuth == nil) || (g.HttpAuth != nil && g.SshAuth != nil) {
		return false
	}
	if g.HttpAuth != nil {
		return g.HttpAuth.IsValid()
	} else {
		return g.SshAuth.IsValid()
	}
	return true
}

func (g *GitConfig) ToPublicModel() {
	if g.HttpAuth != nil {
		g.HttpAuth.ToPublicModel()
	}
	if g.SshAuth != nil {
		g.SshAuth.ToPublicModel()
	}
}
