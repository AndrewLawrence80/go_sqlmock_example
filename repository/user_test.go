package repository

import (
	"database/sql"
	"example/entity"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	// 使用sqlmock创建一个模拟数据库连接，使用sqlmock.QueryMatcherRegexp匹配查询语句
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)

	// 使用gorm连接模拟数据库连接
	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	// 创建gorm数据库连接
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return db, mock, gormDB
}

func TestNewUserRepository(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want UserRepository
	}{
		{
			name: "TestNewUserRepository",
			args: args{
				db: nil,
			},
			want: &userRepository{db: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepository_Create(t *testing.T) {
	db, mock, gormDB := newMockDB(t)
	defer db.Close()

	newUser := &entity.User{
		ID:    1,
		Name:  "John",
		Email: "john@example.com",
		Age:   20,
	}
	mock.ExpectBegin()
	// 正则表达式的所有空格与特殊字符都使用 .* 代替
	// .* 表示匹配任意字符
	// 返回最后插入的ID和受影响的行数
	mock.ExpectExec("INSERT.*INTO.*users.*(`name`,`email`,`age`,`id`).*VALUES.*").WithArgs(newUser.Name, newUser.Email, newUser.Age, newUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		user *entity.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create",
			fields: fields{
				db: gormDB,
			},
			args: args{
				user: newUser,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepository{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepository_GetByID(t *testing.T) {
	db, mock, gormDB := newMockDB(t)
	defer db.Close()

	// 创建一个模拟的用户
	rows := sqlmock.NewRows([]string{"id", "name", "email", "age"}).AddRow(1, "John", "john@example.com", 20)
	// 设置期望的查询语句和返回结果
	// gorm 会默认带 LIMIT 子句查询，因此参数为 1, 1
	// SELECT * FROM `users`` WHERE `users`.`id` = 1 ORDER BY `users`.`id` LIMIT 1
	// 正则表达式的所有空格与特殊字符都使用 .* 代替
	// .* 表示匹配任意字符
	mock.ExpectQuery("SELECT.*FROM.*users.*WHERE.*id.*=.*ORDER.*BY.*LIMIT.*").WithArgs(1, 1).WillReturnRows(rows)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "GetUserById",
			fields: fields{
				db: gormDB,
			},
			args: args{
				id: 1,
			},
			want: &entity.User{
				ID:    1,
				Name:  "John",
				Email: "john@example.com",
				Age:   20,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepository{
				db: tt.fields.db,
			}
			got, err := r.GetByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepository.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepository_Update(t *testing.T) {
	db, mock, gormDB := newMockDB(t)
	defer db.Close()

	updatedUser := &entity.User{
		ID:    1,
		Name:  "john",
		Email: "john@example.com",
		Age:   20,
	}

	mock.ExpectBegin()
	// 返回最后插入的ID和受影响的行数
	mock.ExpectExec("UPDATE.*users.*SET.*name=?.*email=?.*age=?.*WHERE.*id.*=.*").WithArgs(updatedUser.Name, updatedUser.Email, updatedUser.Age, updatedUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		user *entity.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Update",
			fields: fields{
				db: gormDB,
			},
			args: args{
				user: updatedUser,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepository{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepository_Delete(t *testing.T) {
	db, mock, gormDB := newMockDB(t)
	defer db.Close()

	deletedUser := &entity.User{
		ID: 1,
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE.*FROM.*users.*WHERE.*id.*=.*").WithArgs(deletedUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{

		{
			name: "Delete",
			fields: fields{
				db: gormDB,
			},
			args: args{
				id: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepository{
				db: tt.fields.db,
			}
			if err := r.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("userRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
