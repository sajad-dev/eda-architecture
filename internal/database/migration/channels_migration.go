package migration

func ChannelMigration() []string {
	return []string{
		IntPrimary("id", true),
		Char("secret_key", 255, false, "", true, false),
		Char("public_key", 255, false, "", true, false),
		Timestamp("created_at", false, "0000-00-00 00:00:00", false, false),
	}
}
