package errorhandling

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}
