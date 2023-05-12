package mdf

type Option struct {
	isMigrate        bool //默认为false，依据启动参数设置
	isUpgrade        bool //默认为false，依据启动参数设置
	enableRule       bool //默认为false
	enableWidget     bool //默认为false
	isBaseDataCenter bool //默认为false
	enableCron       bool //默认为false
	isRegistry       bool //默认为false
	enableAuthToken  bool //默认为true
}

func newOption() *Option {
	return &Option{enableAuthToken: true}
}
func (s *Server) WithOptionRule() func(*Option) {
	return func(r *Option) {
		r.enableRule = true
	}
}
func (s *Server) WithOptionWidget() func(*Option) {
	return func(r *Option) {
		r.enableWidget = true
	}
}
func (s *Server) WithOptionMigrate() func(*Option) {
	return func(r *Option) {
		r.isMigrate = true
	}
}
func (s *Server) WithOptionAuthToken(enable bool) func(*Option) {
	return func(r *Option) {
		r.enableAuthToken = enable
	}
}
func (s *Server) WithOptionBaseDataCenter() func(*Option) {
	return func(r *Option) {
		r.isBaseDataCenter = true
	}
}
