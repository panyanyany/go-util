package go_util

func Must(err error)  {
	if err != nil {
		panic(err)
	}
}
