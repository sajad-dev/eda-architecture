package command

import (
	"fmt"

	"github.com/sajad-dev/eda-architecture/internal/database/migration"
)

func Migrate(args []string) {
	fmt.Println(args)
	switch args[1]{
	case "create":
		migration.CreateAll()
	case "drop":
		migration.DropTable()
	}

}