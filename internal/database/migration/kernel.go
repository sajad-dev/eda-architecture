package migration

var ArrMigrations = []*Migrate{
	&Migrate{
		Table:     "channels",
		Function: ChannelMigration,
	},
}

func MigrationList(migrate []*Migrate) {
	ArrMigrations = append(ArrMigrations, migrate...)
}
