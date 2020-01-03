package repository

import "os"

type FromCMD struct{}

func (cmd *FromCMD) Check(argsNeed int) bool {
	return argsNeed == len(os.Args)
}
