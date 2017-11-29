package target_applications

import (
	"fmt"

	"github.com/dminGod/ZooGuard/configParsers"
	"github.com/dminGod/ZooGuard/spoc"
)

var PgCluster configParsers.Pgctl_parser
var PgConf configParsers.PgConf

//Load_cluster gets the information from the pgxc_ctl.conf file
//and loads the clusters of the application
func Load_cluster() {

	PgCluster.Init()

	for _, v := range spoc.AppConnections.Connections {

		if v.ApplicationType == "postgresxl" {

			t := `( set -o posix ; set) >/tmp/variables.before
				source %v
				(set -o posix ; set) >/tmp/variables.after
				diff /tmp/variables.before /tmp/variables.after
			`

			command := fmt.Sprintf(t, v.ConfigFile)

			c := spoc.ClientConnections.GetServerByName(v.Server)

			if c == nil {
				fmt.Println("Skipping nil connectionn")
				continue
			} else {

				fmt.Println("Going forward")
			}

			str := spoc.RunCommand(c, command)

			fmt.Println("response from server::::", str)
			PgCluster.ParseString(str)

			for i := range PgCluster.Cluster.Datanodes {

				PgNode_details(&(PgCluster.Cluster.Datanodes[i]))
				fmt.Printf("Datanode Server Configuration:\n \n Datanode Ident configuration\n %+v \n Datanode HBAConfiguration\n %+v\n", PgCluster.Cluster.Datanodes[i].IdentConfiguration, PgCluster.Cluster.Datanodes[i].HbaConfiguration)
			}

			/*for i, _ := range PgCluster.Cluster.Coord {
				PgNode_details(&(PgCluster.Cluster.Coord[i]))
				//fmt.Println("Coord server config:", PgCluster.Cluster.Coord[i].ServerConfiguration)

			}

			for i, _ := range PgCluster.Cluster.DatanodeSlaves {

				PgNode_details(&(PgCluster.Cluster.DatanodeSlaves[i]))
				//fmt.Println("Datanode slave configuration", PgCluster.Cluster.DatanodeSlaves[i].ServerConfiguration)
			}

			for i, _ := range PgCluster.Cluster.CoordSlaves {

				PgNode_details(&(PgCluster.Cluster.CoordSlaves[i]))

			}*/

			/*for i, v := range PgCluster.Cluster.GTMProxies {

				cmd := fmt.Sprintf("cat %v/gtm_proxy.conf", v.GtmProxyDir)

				if v.ServerConn != nil {
					fmt.Println("proxy dir:", v.GtmProxyDir)
					fmt.Println(v.ServerConn.Server_ip)
					kk := spoc.RunCommand(v.ServerConn, cmd)
					fmt.Println("printing kk of gtmproxy", kk)

					var pp configParsers.Pg_conf

					pp.Set_contents(kk)
					pp.Parse()

					PgCluster.Cluster.GTMProxies[i].ServerConfiguration = pp
					//fmt.Printf("Printing pp for proxies %+v \n", pp)
				} else {
					fmt.Println("Server configuration unavailable")
				}

			}

			v := PgCluster.Cluster.GtmMaster

			cmd := fmt.Sprintf("cat %v/gtm.conf", v.GtmMasterDir)

			if v.ServerConn != nil {
				kk := spoc.RunCommand(v.ServerConn, cmd)

				var pp configParsers.Pg_conf

				pp.Set_contents(kk)
				pp.Parse()

				PgCluster.Cluster.GtmMaster.ServerConfiguration = pp
				fmt.Printf("Printing pp for master %+v \n", PgCluster.Cluster.GtmMaster.ServerConfiguration)
			} else {
				fmt.Println("Server configuration unavailable")
			}

			vv := PgCluster.Cluster.GtmSlave

			cmdi := fmt.Sprintf("cat %v/gtm.conf", vv.GtmSlaveDir)

			if vv.ServerConn != nil {

				kkk := spoc.RunCommand(vv.ServerConn, cmdi)

				var ppp configParsers.Pg_conf

				ppp.Set_contents(kkk)
				ppp.Parse()

				PgCluster.Cluster.GtmSlave.ServerConfiguration = ppp
				fmt.Printf("printing pp for gtm slave %+v \n", ppp)
			} else {
				fmt.Println("Server configuration not available")
			}*/

		}
	}
}

//PgNode_details gets the configuration of all the nodes in the postgres-xl database
func PgNode_details(p configParsers.PgNode) {

	k := p.GetPgConfig()

	//postgresql.conf
	var pp configParsers.Pg_conf
	//fmt.Println("Directory is:", k.PgDir, k.ServerIp)
	cmd := fmt.Sprintf("cat %v/postgresql.conf", k.PgDir)
	s := spoc.RunCommand(k.Conn, cmd)
	pp.Set_contents(s)
	pp.Parse()
	p.SetPgConfig(pp)

	//pg_ident.conf
	var pg configParsers.Pg_ident
	cmdi := fmt.Sprintf("cat %v/pg_ident.conf", k.PgDir)
	si := spoc.RunCommand(k.Conn, cmdi)
	pg.Set_contents(si)
	pg.Parse()
	p.SetIdentConfig(pg)

	//pg_hba.conf
	var ph configParsers.Pg_hba
	cmdh := fmt.Sprintf("cat %v/pg_hba.conf", k.PgDir)
	sh := spoc.RunCommand(k.Conn, cmdh)
	ph.Set_contents(sh)
	ph.Parse()
	p.SetHbaConfig(ph)

}