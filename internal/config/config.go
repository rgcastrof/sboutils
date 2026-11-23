package config

/* Configuration constants */
const (
	path = "/usr/sbo"
	url = "https://github.com/Ponce/slackbuilds.git"
	branch = "current"
	tmp = "/tmp"
	gitBin = "/usr/bin/git"
	shellBin = "/bin/sh"
	numJobs = 1
)

type Config struct {
	RepoPath string
	RepoUrl string
	RepoBranch string
	Tmp string
	GitBin string
	ShellBin string
	NumJobs int
}

func LoadConfig() *Config {
	return &Config{
		RepoPath: path,
		RepoUrl: url,
		RepoBranch: branch,
		Tmp: tmp,
		GitBin: gitBin,
		ShellBin: shellBin,
		NumJobs: numJobs,
	}
}

func (c *Config) SyncArgs() []string {
	return []string{
		"clone", "--depth", "1", "--no-tags",
		"--branch", c.RepoBranch, "--", c.RepoUrl, c.RepoPath,
	}
}

func (c *Config) PullArgs() []string {
	return []string{"-C", c.RepoPath, "pull", "--rebase"}
}

func (c *Config) CleanArgs() []string {
	return []string{"-C", c.RepoPath, "clean", "-fdx"}
}

func (c *Config) BuildArgs() []string {
	return []string{"--login"}
}
