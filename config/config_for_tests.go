package config

type TestConfig struct {
	DatabaseType             string `yaml:"db_type"`
	DatabaseConnectionString string `yaml:"db_connection_string"`
	ProblemsPath             string `yaml:"problems_path"`
	SolutionsPath            string `yaml:"solutions_path"`
}
