package database

import (
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
)

// Use calls the use method on each object in the UseCase
//  chain. If an error occurs, the function immediately
//  stops working and returns an error
func Use(db Interface, ucases ...UserCaseI) error {
	var err error
	for _, ucase := range ucases {
		err = ucase.Use(db)
		if err != nil {
			break
		}
	}
	return re.Wrap(err)
}
