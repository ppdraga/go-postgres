package service

type Controller struct {
	ff FilesFinder
}

type FilesFinder interface {
	FindUserFiles(username string) (*[]string, error)
	FindServerFiles(servername string) (*[]string, error)
	FindAreaFiles(areaname string) (*[]string, error)
}

func NewController(ff FilesFinder) *Controller {
	return &Controller{
		ff: ff,
	}
}

func (c *Controller) GetUserFiles(username string) (*[]string, error) {
	return c.ff.FindUserFiles(username)
}

func (c *Controller) GetServerFiles(servername string) (*[]string, error) {
	return c.ff.FindServerFiles(servername)
}

func (c *Controller) GetAreaFiles(areaname string) (*[]string, error) {
	return c.ff.FindAreaFiles(areaname)
}
