package command

import (
	"github.com/sajad-dev/eda-architecture/internal/database/migration"
)

func Migrate(args []string) {
	switch args[1]{
	case "create":
		migration.CreateAll()
	case "drop":
		migration.DropTable()
	}

}