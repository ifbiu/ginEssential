package main

import (
	"fmt"
	_ "database/sql"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main()  {
	db := InitDB();
	defer db.Close();
	r := gin.Default()
	r.GET("/api/auth/register", func(ctx *gin.Context) {
		// 获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")
		// 数据验证
		if len(telephone) != 11 {
			ctx.JSON(422,gin.H{"code":422,"msg":"手机号必须为11位"})
			return
		}
		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"密码不能小于6位"})
			return
		}
		if len(name) == 0{
			name = RandomString(10)
		}
		if isTelephoneExist(db,telephone) {
			ctx.JSON(422,gin.H{"code":422,"msg":"用户已经存在"})
			return
		}
		log.Println(name,password,telephone)

		// 创建用户
		newUser := User{
			Name : name,
			TelePhone: telephone,
			Password: password,
		}
		db.Create(&newUser)


		ctx.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

type User struct{
	gorm.Model
	Name string `gorm:"type:varchar(20);NOT NULL"`
	Password string `gorm:"type:varchar(255);NOT NULL"`
	TelePhone string `gorm:"type:varchar(11);NOT NULL UNIQUE"`
}

func RandomString(n int)string{
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result := make([]byte,n)
	rand.Seed(time.Now().Unix())
	for i:=range result{
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func InitDB() *gorm.DB  {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "123456"
	charset := "utf8mb4"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset,
		)
	db,err := gorm.Open(driverName,args)
	if err != nil {
		panic("failed to connect database,err:" + err.Error())
	}
	db.AutoMigrate(&User{})
	return db;
}

func isTelephoneExist(db *gorm.DB,telephone string) bool {
	var user User
	db.Where("telephone = ?",telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}