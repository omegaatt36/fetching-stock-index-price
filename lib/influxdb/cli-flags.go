package influxdb

import (
	"github.com/urfave/cli"
)

// CliFlags returns flags to setup influxdb package.
func (p *connPool) CliFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, cli.StringFlag{
		Name:   "influx-host",
		Value:  "http://influxdb:8086",
		EnvVar: "INFLUX_HOST",
	})
	flags = append(flags, cli.StringFlag{
		Name:   "influx-db",
		Value:  "db0",
		EnvVar: "INFLUX_DB",
	})
	flags = append(flags, cli.StringFlag{
		Name:   "influx-user",
		Value:  "admin",
		EnvVar: "INFLUX_USER",
	})
	flags = append(flags, cli.StringFlag{
		Name:   "influx-pass",
		Value:  "aaaa1234",
		EnvVar: "INFLUX_PASS",
	})
	flags = append(flags, cli.IntFlag{
		Name:   "influx-max-conn",
		Value:  6,
		EnvVar: "INFLUX_MAX_CONN",
	})

	return flags
}

// InitFromCliFlags inits influxdb package from cli flags.
func (p *connPool) InitFromCliFlags(c *cli.Context) {
	host := c.String("influx-host")
	db := c.String("influx-db")
	user := c.String("influx-user")
	pass := c.String("influx-pass")
	maxConn := c.Int("influx-max-conn")

	p.host = host
	p.db = db
	p.user = user
	p.pass = pass
	p.maxConn = maxConn
	p.init()
}

func (p *connPool) Finalize() {
	// no need.
}
