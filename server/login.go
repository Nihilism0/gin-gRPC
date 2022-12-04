package server

import (
	"context"
	"gin_demo/model"
	"gin_demo/pb/proto"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"net"
	"time"
)

func Loginserver() {
	// 监听端口
	log.Println("监听端口中")
	lis, err := net.Listen("tcp", ":50051")
	log.Println("已监听50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer() //获取新服务示例
	proto.RegisterLoginServer(s, &loginserver{})

	// 开始处理
	err = s.Serve(lis)
	log.Println("开始处理")
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type loginserver struct {
	proto.UnimplementedLoginServer // 用于实现proto包里LoginServer接口
}

func (s *loginserver) mustEmbedUnimplementedLoginServer() {
	//TODO implement me
	panic("implement me")
}

func (s *loginserver) Login(ctx context.Context, req *proto.UserReq) (*proto.UserResp, error) {
	resp := &proto.UserResp{}
	log.Println("recv:", req.UserName, req.PassWord)
	if req.PassWord != GetPassWord(req.UserName) {
		resp.OK = false
		return resp, nil
	}
	resp.OK = true
	return resp, nil
}

func GetPassWord(userName string) (password string) {
	log.Println("传入", userName)
	var u struct {
		Password string
	}
	db, _ := gorm.Open(mysql.New(mysql.Config{ //配置
		DSN: "root:123456@tcp(127.0.0.1:3306)/gindemo?charset=utf8mb4&parseTime=True&loc=Local",
	}), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "gindemo_",
			SingularTable: false,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxLifetime(10) //数据池
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	log.Println("gRPC数据库连接成功")
	db.Model(&model.User{}).Where("username = ?", userName).Find(&u)
	log.Println("已经接受到username返回password", u.Password)
	return u.Password
}
