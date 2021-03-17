package commands

import "github.com/ugent-library/momo/internal/theme"

func init() {
	theme.Register(theme.New("ugent"))
	theme.Register(theme.New("orpheus"))
}
