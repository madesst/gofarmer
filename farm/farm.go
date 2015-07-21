package farm

import (
	"github.com/codegangsta/cli"
)

func Create(ctx *cli.Context) {
	/*
		1. Check and prepare internal dirs
		2. Check and read global config
		3. Check cli input if auth info does not exist in global config
		4. Create new dir with name from input
		5. Save typical farm config in new dir from step 4
	*/
}
