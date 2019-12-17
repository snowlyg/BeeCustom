package models

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/astaxie/beego"
)

func DataReset() (bool, error) {

	dbType := beego.AppConfig.String("db_type")
	//数据库名称
	//dbName := beego.AppConfig.String(dbType+"::db_name")
	//数据库连接用户名
	dbUser := beego.AppConfig.String(dbType + "::db_user")
	//数据库连接用户名
	dbPwd := beego.AppConfig.String(dbType + "::db_pwd")
	//数据库IP（域名）
	dbHost := beego.AppConfig.String(dbType + "::db_host")
	//数据库端口
	dbPort := beego.AppConfig.String(dbType + "::db_port")

	arv := []string{fmt.Sprintf("--host=%s", dbHost), fmt.Sprintf("--port=%d", dbPort), fmt.Sprintf("-u%s", dbUser), fmt.Sprintf("-p%s", dbPwd), ">", "bee_custom_clearances.sql"}

	cmd := exec.Command("mysql", arv...)

	var out bytes.Buffer

	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {

		return false, err
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
	return true, nil
}
