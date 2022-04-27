package actors

import (
	"astroterm/astro"
	"astroterm/db"
	"astroterm/project"
	"regexp"
	"strconv"
)

var portMatch = regexp.MustCompile("(localhost|127.0.0.1):([0-9]{4})(\\/.*)")

type DevServerActor struct {
	DB      *db.Database
	Model   *db.DevServerModel
	project *project.Project

	bstarting     *broker[bool]
	bstopped      *broker[any]
	binitialstate *broker[*db.DevServerModel]
	bhostinfo     *broker[*db.DevServerModel]
	blogs         *broker[[]byte]
}

func NewDevServerActor(project *project.Project) *DevServerActor {
	return &DevServerActor{
		DB:            db.NewDatabase(),
		Model:         &db.DevServerModel{},
		project:       project,
		bstarting:     newBroker[bool](),
		bstopped:      newBroker[any](),
		binitialstate: newBroker[*db.DevServerModel](),
		bhostinfo:     newBroker[*db.DevServerModel](),
		blogs:         newBroker[[]byte](),
	}
}

func (d *DevServerActor) Start() *DevServerActor {
	go d.startup()
	return d
}

func (d *DevServerActor) SubscribeToStarting() chan bool {
	return d.bstarting.Subscribe()
}

func (d *DevServerActor) SubscribeToStopped() chan any {
	return d.bstopped.Subscribe()
}

func (d *DevServerActor) SubscribeToInitialState() chan *db.DevServerModel {
	return d.binitialstate.Subscribe()
}

func (d *DevServerActor) SubscribeToHostInfo() chan *db.DevServerModel {
	return d.bhostinfo.Subscribe()
}

func (d *DevServerActor) SubscribeToLogs() chan []byte {
	return d.blogs.Subscribe()
}

func (d *DevServerActor) StartDevServer() chan bool {
	done := make(chan bool)
	go func() {
		d.bstarting.Publish(true)
		d.startServer()
	}()
	return done
}

func (d *DevServerActor) StopDevServer() chan bool {
	done := make(chan bool)
	go func() {
		d.shutdownServer()
		d.Model.Pid = 0
		d.bstopped.Publish(true)
		done <- true
	}()
	return done
}

func (d *DevServerActor) startup() {
	go d.bstarting.Start()
	go d.bstopped.Start()
	go d.binitialstate.Start()
	go d.bhostinfo.Start()
	go d.blogs.Start()

	d.LoadDevServerModel()
	d.binitialstate.Publish(d.Model)
}

func (d *DevServerActor) LoadDevServerModel() {
	projectDir := d.project.Dir
	model, err := d.DB.LoadDevServerModel(projectDir)
	if err != nil {
		d.Model.ProjectDir = projectDir
		return
	}
	if model != nil {
		d.Model = model
	} else {
		d.Model.ProjectDir = projectDir
	}
}

func (d *DevServerActor) startServer() error {
	cmd, err := astro.RunCommand(astro.Dev, d)
	if err != nil {
		return err
	}
	d.Model.Pid = cmd.Process.Pid
	d.DB.AddStartedDevServer(d.Model)
	return nil
}

func (d *DevServerActor) Write(p []byte) (int, error) {
	if d.Model.Port == 0 {
		hostname, port, subpath := d.parseHostInformation(p)
		if hostname != "" {
			d.Model.Port = port
			d.Model.Hostname = hostname
			d.Model.Subpath = subpath
			d.saveDevServerInformation()

			d.bhostinfo.Publish(d.Model)
		}
	}

	d.blogs.Publish(p)
	return 0, nil
}

func (d *DevServerActor) parseHostInformation(p []byte) (string, int, string) {
	part := string(p)
	rs := portMatch.FindStringSubmatch(part)
	if len(rs) > 1 {
		portString := rs[2]
		port, _ := strconv.Atoi(portString)
		hostname := rs[1]
		subpath := rs[3]
		return hostname, port, subpath
	}
	return "", 0, ""
}

func (d *DevServerActor) saveDevServerInformation() error {
	err := d.DB.SetDevServerInformation(d.Model)
	return err
}

func (d *DevServerActor) killServer() error {
	if d.Model.Pid != 0 {
		err1 := killPid(d.Model.Pid + 1)
		err2 := killPid(d.Model.Pid)

		if err1 != nil {
			return err1
		}
		return err2
	}
	return nil
}

func (d *DevServerActor) shutdownServer() error {
	e1 := d.killServer()
	e2 := d.DB.DeleteDevServer(d.Model)
	if e1 != nil {
		return e1
	}
	return e2
}
