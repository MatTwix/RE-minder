package migrations

import "github.com/jackc/pgx/v5/pgxpool"

func Migrate(db *pgxpool.Pool) {
	CreateUsersTable(db)
	CreateHabitsTable(db)
	AddIsAdminField(db)
	AddStartDateField(db)
	CreateNotificationSettingsTable(db)
	MakeUserIdFieldUniqueNotificationsSettingsTable(db)
	СhangeVkFieldsToGoogleNotificationsTable(db)
}
