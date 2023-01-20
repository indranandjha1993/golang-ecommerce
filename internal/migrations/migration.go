package migrations

import (
	"golang-ecommerce/utils"

	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB, action string) utils.ErrorWarper {
	var errs utils.ErrorWarper
	switch action {
	case "up":
		if err := createUsersTable(db); err != nil {
			errs.Add("create user: ", err.Error())
		}
		if err := createUsersTokenTable(db); err != nil {
			errs.Add("create user token: ", err.Error())
		}
	case "down":
		if err := dropUsersTable(db); err != nil {
			errs.Add("drop user: ", err.Error())
		}
		if err := dropUsersTokenTable(db); err != nil {
			errs.Add("drop user: ", err.Error())
		}
	default:
		errs.Add("invalid action provided: ", action)
	}
	return errs
}
