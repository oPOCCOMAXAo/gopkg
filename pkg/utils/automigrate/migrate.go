package automigrate

import (
	pkgerr "github.com/pkg/errors"
	"gorm.io/gorm"
)

func CreateTables(
	migrator gorm.Migrator,
	modelsPtrs ...any,
) error {
	err := migrator.AutoMigrate(modelsPtrs...)
	if err != nil {
		return pkgerr.WithStack(err)
	}

	return nil
}

func DropTables(
	migrator gorm.Migrator,
	unusedTablesNames ...string,
) error {
	for _, table := range unusedTablesNames {
		err := migrator.DropTable(
			table,
		)
		if err != nil {
			return pkgerr.WithStack(err)
		}
	}

	return nil
}

func DropIndexesSingle(
	migrator gorm.Migrator,
	modelPtr any,
	indexes ...string,
) error {
	for _, index := range indexes {
		if migrator.HasIndex(modelPtr, index) {
			err := migrator.DropIndex(modelPtr, index)
			if err != nil {
				return pkgerr.WithStack(err)
			}
		}
	}

	return nil
}
