package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/dreamans/syncd"
	"github.com/dreamans/syncd/model"
	"github.com/jinzhu/gorm"
)

var (
	SqliteFile string
	Hostname   string
	Port       int
	Username   string
	Password   string
	DbName     string
)

func init() {
	flag.StringVar(&SqliteFile, "sqlite", "", "set using sqlite database")
	flag.StringVar(&Hostname, "hostname", "127.0.0.1", "mysql host ip")
	flag.IntVar(&Port, "port", 3306, "port")
	flag.StringVar(&Username, "user", "sms", "mysql username")
	flag.StringVar(&Password, "password", "password", "mysql password")
	flag.StringVar(&DbName, "dbname", "sms", "mysql database name")
}

func getMd5(s string) string {
	m := md5.New()
	io.WriteString(m, s)
	return hex.EncodeToString(m.Sum(nil))
}

func initRole(db *gorm.DB) {
	privstr := `2001,2002,2003,2004,2100,2101,2102,2201,2202,2203,2204,2205,2206,2207,3001,3002,3004,3003,3101,3102,3103,3104,4001,4002,4003,4004,4101,4102,4103,4104,1001,1002,1006,1003,1004,1005`
	now := int(time.Now().Unix())
	role := &model.UserRole{ID: 1, Name: "管理员", Privilege: privstr, Ctime: now}
	db.Create(role)

	if err := db.Error; err != nil {
		log.Fatal(err)
		return
	}
}

func initUser(db *gorm.DB) {
	salt := "u0EMxuE6qh"
	now := int(time.Now().Unix())
	user := &model.User{
		RoleId:        1,
		Username:      "admin",
		Password:      getMd5(getMd5("111111") + salt),
		Salt:          salt,
		Truename:      "管理员",
		Mobile:        "13800012345",
		Email:         "admin@admin.com",
		Status:        1,
		LastLoginIp:   "1.1.1.1",
		LastLoginTime: 15,
		Ctime:         now,
	}
	db.Create(user)
	if err := db.Error; err != nil {
		log.Fatal(err)
		return
	}
}

func insertData(db *gorm.DB) {
	initRole(db)
	initUser(db)

	log.Println("init table data done")
}

func run() {

	var dbCfg = &syncd.DbConfig{
		SqliteFile: SqliteFile,
		Host:       Hostname,
		User:       Username,
		Pass:       Password,
		Port:       Port,
		DbName:     DbName,
		Charset:    "utf8mb4",
	}
	var db = syncd.NewDatabase(dbCfg)
	err := db.Open()
	if err != nil {
		log.Println(err)
		return
	}

	var objs = []interface{}{
		new(model.UserToken),
		new(model.UserRole),
		new(model.User),
		new(model.DeployApply),
		new(model.DeployBuild),
		new(model.DeployTask),
		new(model.Project),
		new(model.ProjectMember),
		new(model.ProjectSpace),
		new(model.Server),
		new(model.ServerGroup),
	}
	for _, o := range objs {
		var c *gorm.DB

		c = db.DbHandler.DropTableIfExists(o)
		if err := c.Error; err != nil {
			log.Println(err)
		}

		dbType := db.DbHandler.Dialect().GetName()
		if strings.ToLower(dbType) == "mysql" {
			c = db.DbHandler.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
		}
		c = c.AutoMigrate(o)
		if err := c.Error; err != nil {
			log.Println(err)
		}

		db.DbHandler.CreateIndexes(o)
	}

	insertData(db.DbHandler)
}

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Llongfile)

	fmt.Println("This app will drop all tables and re-create them!")

	var ok string
	for {
		fmt.Printf("Are you sure? [Yes/No]:")
		fmt.Scanf("%s", &ok)

		ok = strings.ToLower(ok)

		switch ok {
		case "n":
			fallthrough
		case "no":
			return

		case "y":
			fallthrough
		case "ye":
			fallthrough
		case "yes":
			break

		default:
			continue
		}

		run()
		return
	}

}
