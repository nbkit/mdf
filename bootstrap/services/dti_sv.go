package services

//interface
type IDtiSv interface {
}

func DtiSv() IDtiSv {
	return dtiSv
}

var dtiSv IDtiSv = newDtiSvImlp()

//impl
type dtiSvImpl struct {
}

func newDtiSvImlp() *dtiSvImpl {
	return &dtiSvImpl{}
}
