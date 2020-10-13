package modelstest

import (
	"testing"

	mathbattle "mathbattle/models"

	"github.com/stretchr/testify/require"
)

func multStr(base string, count int) string {
	if count <= 0 {
		return ""
	}

	result := ""
	for i := 0; i < count; i++ {
		result = result + base
	}

	return result
}

func TestIsNameValid(t *testing.T) {
	req := require.New(t)

	req.False(mathbattle.IsParticipantNameValid(""))
	req.False(mathbattle.IsParticipantNameValid(multStr("a", 31)))
	req.False(mathbattle.IsParticipantNameValid("12345"))
	req.False(mathbattle.IsParticipantNameValid("John!"))
	req.False(mathbattle.IsParticipantNameValid("Василий1234"))
	req.False(mathbattle.IsParticipantNameValid(multStr("ы", 31)))

	req.True(mathbattle.IsParticipantNameValid("John"))
	req.True(mathbattle.IsParticipantNameValid("John Doe"))
	req.True(mathbattle.IsParticipantNameValid(multStr("a", 30)))
	req.True(mathbattle.IsParticipantNameValid("Василий"))
	req.True(mathbattle.IsParticipantNameValid("Василий Кроссовкин"))
	req.True(mathbattle.IsParticipantNameValid(multStr("ы", 30)))
}
