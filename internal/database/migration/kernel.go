package migration

var ArrMigrations = []*Migrate{
	{
		Table:     "channels",
		Function: ChannelMigration,
	},
}

func MigrationList(migrate []*Migrate) {
	ArrMigrations = append(ArrMigrations, migrate...)
}
