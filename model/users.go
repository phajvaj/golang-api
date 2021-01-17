package model

import (
	"database/sql"
	"fmt"
	"nan_api_main/entities"

	"github.com/bdwilliams/go-jsonify/jsonify"
)

type UserJson struct {
	Username string `json: "username"`
	Password string `json: "password"`
}

type UsersModel struct {
	Db *sql.DB
}

var tbl = "users"

func (model UsersModel) Find(id int) ([]string, error) { //entities.Users
	rows, err := model.Db.Query("select * from "+tbl+" where id = ?", id)
	if err != nil {
		return jsonify.Jsonify(rows), err
	} else {
		return jsonify.Jsonify(rows), nil
	}
	/*if err != nil {
		return entities.Users{}, err
	} else {
		var _user entities.Users
		for rows.Next() {
			var userId int64
			var username string
			var password string
			var firstName string
			var lastName string
			var userType string
			var isEnabled string
			err2 := rows.Scan(&userId, &username, &password, &firstName, &lastName, &userType, &isEnabled)
			if err2 != nil {
				return entities.Users{}, err2
			} else {
				_user = entities.Users{userId, username, password, firstName, lastName, userType, isEnabled}
			}
		}
		return _user, nil
	}*/
}

func (model UsersModel) FindLogin(_username, _password string) ([]entities.UsersLogin, error) { //[]entities.UsersLogin
	rows, err := model.Db.Query("select user_id,username,first_name,last_name,user_type from "+tbl+" where username = ? and password = ? and is_enabled='Y' limit 1", _username, _password)
	/*if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		//fmt.Println(jsonify.Jsonify(rows))
		return jsonify.Jsonify(rows), nil
	}*/

	if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		user := []entities.UsersLogin{}
		for rows.Next() {
			data := entities.UsersLogin{}
			err2 := rows.Scan(&data.UserId, &data.Username, &data.FirstName, &data.LastName, &data.UserType)
			if err2 != nil {
				return nil, err2
			} else {
				user = append(user, data)
			}
		}
		return user, nil
	}

}
