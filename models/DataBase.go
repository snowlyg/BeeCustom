package models

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/astaxie/beego"
)

func DataReset() (bool, error) {

	dbType := beego.AppConfig.DefaultString("db_type", "mysql")
	//数据库名称
	//dbName := beego.AppConfig.DefaultString(dbType+"::db_name", "bee_custom")
	//数据库连接用户名
	dbUser := beego.AppConfig.DefaultString(dbType+"::db_user", "root")
	//数据库连接用户名
	dbPwd := beego.AppConfig.DefaultString(dbType+"::db_pwd", "")
	//数据库IP（域名）
	dbHost := beego.AppConfig.DefaultString(dbType+"::db_host", "127.0.0.1")
	//数据库端口
	dbPort := beego.AppConfig.DefaultString(dbType+"::db_port", "3306")

	arv := []string{fmt.Sprintf("--host=%s", dbHost), fmt.Sprintf("--port=%d", dbPort), fmt.Sprintf("-u%s", dbUser), fmt.Sprintf("-p%s", dbPwd), ">", "bee_custom.sql"}

	cmd := exec.Command("mysql ", arv...)

	var out bytes.Buffer

	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {

		fmt.Println(err)
	}

	fmt.Printf(out.String())

	//f, err := os.Open("bee_custom.sql")
	//if err != nil {
	//	return false, err
	//}
	//
	//defer f.Close()
	//
	//o := orm.NewOrm()
	//buf := bufio.NewReader(f)
	//for {
	//	line, err := buf.ReadString(';')
	//	if err != nil {
	//		if err == io.EOF {
	//			return true, nil
	//		}
	//		return false, err
	//	}
	//	_, err = o.Raw(line).Exec()
	//	if err != nil {
	//		return false, err
	//	}
	//}

}
