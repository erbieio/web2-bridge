package mysql

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/erbieio/web2-bridge/config"
	"github.com/erbieio/web2-bridge/utils/logger"

	"net"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var once sync.Once

// func init() {
// 	GetDB()
// }

type ViaSSHDialer struct {
	client *ssh.Client
	_      *context.Context
}

func (sshdial *ViaSSHDialer) Dial(context context.Context, addr string) (net.Conn, error) {
	return sshdial.client.Dial("tcp", addr)
}

// DailWithKey connect  remote ssh server through CA key
func DailWithKey(addr, user, keyfile string) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return ssh.Dial("tcp", addr, config)
}

// GetDB get mysql db instance by sync.Once
func GetDB() *gorm.DB {
	once.Do(func() {
		// Connect to the SSH Server
		var dsn string
		if config.GetMysqlConfig().UseCA {
			hostip := fmt.Sprintf("%s:%d", config.GetMysqlConfig().IP, config.GetMysqlConfig().SSHPort)
			client, err := DailWithKey(hostip, config.GetMysqlConfig().SSHAccount, config.GetMysqlConfig().SSHKey)
			if err != nil {
				logger.Logrus.Fatal(fmt.Sprintf("Fatal err dial ssh server: %v", err))
			}

			// Now we register the ViaSSHDialer with the ssh connection as a parameter
			mysql.RegisterDialContext("mysql+tcp", (&ViaSSHDialer{client, nil}).Dial)
			dsn = fmt.Sprintf("%s:%s@mysql+tcp(127.0.0.1:%d)/%s?charset=utf8&parseTime=True&loc=UTC", config.GetMysqlConfig().Account, config.GetMysqlConfig().Password, config.GetMysqlConfig().ConnPort, config.GetMysqlConfig().SqlName)
		} else {
			network := "@tcp"
			// ip := config.GetMysqlConfig().IP
			// network := "@mysql+tcp"
			// if ip == "127.0.0.1" || ip == "localhost" {
			// 	network = "@tcp"
			// }
			dsn = fmt.Sprintf("%s:%s%s(%s:%d)/%s?charset=utf8&parseTime=True&loc=UTC", config.GetMysqlConfig().Account, config.GetMysqlConfig().Password, network, config.GetMysqlConfig().IP, config.GetMysqlConfig().ConnPort, config.GetMysqlConfig().SqlName)
		}

		// And now we can use our new driver with the regular mysql connection string tunneled through the SSH connection

		mdb, err := gorm.Open(gmysql.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.Logrus.Fatal(fmt.Sprintf("Fatal error open mysql: %v", err))
		}

		db = mdb

		//setup mysql conn pool
		sqlDB, err := db.DB()
		if err != nil {
			logger.Logrus.Fatal(fmt.Sprintf("Fatal error convert db: %v", err))
		}

		sqlDB.SetMaxOpenConns(int(config.GetMysqlConfig().MaxOpenConns))
		sqlDB.SetMaxIdleConns(int(config.GetMysqlConfig().MaxIdleConns))

		sqlDB.SetConnMaxLifetime(time.Duration(config.GetMysqlConfig().MaxLifetime) * time.Second)
	})

	return db
}
