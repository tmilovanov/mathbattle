package repositorytest

import (
	"os"
	"path"
)

func DeleteTempDatabase() {
	if _, err := os.Stat(TestStoragePath()); !os.IsNotExist(err) {
		os.RemoveAll(TestStoragePath())
	}
}

func TestStoragePath() string {
	return path.Join(os.TempDir(), "mathbattle_test_storage")
}

func TestDbPath() string {
	return path.Join(TestStoragePath(), "test_mathbattle.sqlite")
}

func TestProblemsPath() string {
	return path.Join(TestStoragePath(), "test_problems")
}

func TestSolutionsPath() string {
	return path.Join(TestStoragePath(), "test_solutions")
}
