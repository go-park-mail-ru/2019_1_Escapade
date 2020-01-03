package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
)

// Open connection and connect every usecase in chain to
//  it. If an error occurs, the function wil stop and
//  return it
func Open(dbi Interface, c config.Database, ucases ...UserCaseI) error {
	err := dbi.Open(c)
	if err != nil {
		return re.Wrap(err)
	}
	for _, ucase := range ucases {
		err = ucase.Use(dbi)
		if err != nil {
			break
		}
	}
	return re.Wrap(err)
}
