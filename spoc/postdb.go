package spoc

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/dminGod/ZooGuard/zgCode"
	// lib/pq is required to run postgres queries
	_ "github.com/lib/pq"
)

//PostDB stores the configuration details of Postgres-XL database
type PostDB struct {
	Name            string
	User            string
	Pass            string
	Host            string
	Port            string
	ParentAppName   string
	LinkedComponent string
	ParentType      string
	ComponentRole   string
	Identity        string
	Info            string
	Conn            *sql.DB
}

//PostConns has information regarding the various Postgres-XL databases connected
type PostConns struct {
	Connections []*PostDB
}

//GetByServerAndRole is used to get the database configuration
//from the server name and the role of the node
func (p *PostConns) GetByServerAndRole(server string, role string) (con PostDB, err error) {

	for _, v := range p.Connections {
		if server == v.Host && role == v.LinkedComponent {
			con.Name = v.Name
			con.Host = v.Host
			con.Port = v.Port
			con.ParentAppName = v.ParentAppName
			con.LinkedComponent = v.LinkedComponent
			con.ParentType = v.ParentType
			con.ComponentRole = v.ComponentRole
			con.Identity = v.Identity
			con.Conn = v.Conn
		}
	}

	if con.Name == "" {
		err = errors.New("Could not find server")
	}

	return

}

func connectPostgres(v zgConfig.Database) {

	var postdb PostDB
	postdb.Name = v.DatabaseName
	postdb.User = v.Username
	postdb.Pass = v.Password
	postdb.ParentAppName = v.ParentAppName
	postdb.LinkedComponent = v.LinkedComponent
	postdb.ParentType = v.ParentType
	postdb.ComponentRole = v.ComponentRole
	postdb.Identity = v.DbIdentity

	rand.Seed(time.Now().UTC().UnixNano())

	for _, curDB := range v.Host {

		// Get the host and port
		curDBSplit := strings.Split(curDB, ":")

		if len(curDBSplit) == 1 {

			//logger.Error( "ErrorType : CONFIG_ERROR, Got server configured without port information skipping server, Server:", curDB)
			continue
		}

		postdb.Host = curDBSplit[0]
		postdb.Port = curDBSplit[1]

		var err error

		if len(postdb.Pass) == 0 {

			postdb.Info = fmt.Sprintf("user=%s dbname=%s sslmode=disable host=%s port=%s",
				postdb.User, postdb.Name, postdb.Host, postdb.Port)
		} else {

			postdb.Info = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
				postdb.User, postdb.Pass, postdb.Name, postdb.Host, postdb.Port)
		}

		postdb.Conn, err = sql.Open("postgres", postdb.Info)
		_ = postdb.Conn

		//dbpoolConn.SetMaxOpenConns(Conf.Postgresxl.MaxOpenConns)
		//dbpoolConn.SetMaxIdleConns(Conf.Postgresxl.MaxIdleConns)
		//dbpoolConn.SetConnMaxLifetime(time.Duration(Conf.Postgresxl.ConnMaxLifetime) * time.Second)

		if err != nil {

			//logger.Error( "ErrorType : INFRA_ERROR, Not able to connect to DB Server:",  Host,"Got error", err.Error())
			fmt.Println("Unable to connect", err)
			continue
		} else {

			//logger.Info( "Adding connection to pool, Server : ", Host, " Port", dbPort)
			//dbpool = append(dbpool, dbpoolConn)
			PostConnections.Connections = append(PostConnections.Connections, &postdb)
			fmt.Println("Added db connection to dbpool")
		}
	}
}

//Query is used to run Select query on Postgres-XL database
func (p_db *PostDB) Query(s string) (retVal []map[string]interface{}) {

	rows, err := p_db.Conn.Query(s)

	if err != nil {
		fmt.Println("QUERY_ERROR, Postgres query failed", err)
		return
	}

	cols, err := rows.Columns()

	if err != nil {
		fmt.Println("QUERY_ERROR, Postgres query failed", err)

		return
	}

	data := make([]interface{}, len(cols))
	args := make([]interface{}, len(data))

	for i := range data {
		args[i] = &data[i]
	}

	for rows.Next() {

		var rowData = make(map[string]interface{})

		if err = rows.Scan(args...); err != nil {
			//logger.Error("ErrorType : QUERY_ERROR, Postgres query failed,  Request ID : ", dbAbstract.RequestID, ", Error when fetching data. Error: ", err.Error())
			fmt.Println("Query failed")
			return
		}

		for i := range data {

			rowData[cols[i]] = data[i]
		}

		retVal = append(retVal, rowData)

	}

	rows.Close()

	if err != nil {

		fmt.Println("Error in closing rows")
	} else {
		fmt.Println("Success")
		fmt.Println(retVal)
	}
	return
}

//Execute is used to run Insert query on Postgres-XL database
func (p_db *PostDB) Execute(s string) (err error) {

	_, err = p_db.Conn.Exec(s)

	if err != nil {
		fmt.Println("QUERY_ERROR, Postgres query failed")

	} else {
		fmt.Println("Postgres inserted succesfully")

	}
	return
}
