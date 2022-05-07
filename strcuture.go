package main

type Session struct {
	num       int
	startTime string
	endTime   string
}

type Dentist struct {
	name string
}

type Patient struct {
	name      string
	mobileNum int
}

func (s *Session) GetSession() int {
	return s.num
}

func (s *Session) GetStartTime() string {
	return s.startTime
}

func (s *Session) GetEndTime() string {
	return s.endTime
}

func (d *Dentist) GetName() string {
	return d.name
}

func (p *Patient) GetName() string {
	return p.name
}

func (p *Patient) GetMobileNum() int {
	return p.mobileNum
}
