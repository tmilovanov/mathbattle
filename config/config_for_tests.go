package config

type TestConfig struct {
	DatabasePath  string `yaml:"db_path"`
	ProblemsPath  string `yaml:"problems_path"`
	SolutionsPath string `yaml:"solutions_path"`
}
