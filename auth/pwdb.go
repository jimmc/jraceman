package auth

import (
  "database/sql"

  "github.com/golang/glog"

  "github.com/jimmc/auth/permissions"
  "github.com/jimmc/auth/users"
)

// PwDB implements the Store interface to load and store data in our SQL database.
// Data is stored in a table called "user" with three string columns,
// id, cryptword, and permissions,
// where the permissions value is a comma-separated list of permission names.
type PwDB struct {
    db *sql.DB
}

func NewPwDB(db *sql.DB) *PwDB {
  return &PwDB{
    db: db,
  }
}

// Load does nothing when we are using a database.
func (pdb *PwDB) Load() error {
  return nil
}

// Save does nothing when we are using a database.
func (pdb *PwDB) Save() error {
  return nil
}

func (pdb *PwDB) User(userid string) *users.User {
  glog.V(2).Infof("Looking for user %s", userid)
  query := `
    SELECT cryptword, GROUP_CONCAT(COALESCE(permission.name,''),',') as permissions
    FROM user
         LEFT JOIN userrole ON user.id=userrole.userid
         LEFT JOIN rolepermission ON userrole.roleid=rolepermission.roleid
         LEFT JOIN permission ON rolepermission.permissionid=permission.id
    WHERE user.username = :username
    GROUP BY user.id`
  var cryptword string
  var perms string
  err := pdb.db.QueryRow(query,sql.Named("username", userid)).Scan(&cryptword, &perms)
  if err == sql.ErrNoRows {
    glog.Warningf("No matching rows found for user %s", userid)
    return nil          // No matching userid found
  }
  if err != nil {
    glog.Errorf("Error scanning for user %q: %v\n", userid, err)
    return nil
  }
  user := users.NewUser(userid, cryptword, permissions.FromString(perms))
  glog.V(2).Infof("Found user %v", user)
  return user
}

func (pdb *PwDB) SetCryptword(userid, cryptword string) {
  if userid == "" {
    glog.Errorf("Can't SetCryptword with no userid\n")
    return
  }
  // Assume row does not exist, try to insert it.
  iQuery := `INSERT into user(id,username,cryptword) values(:id, :username, :cw);`
  _, err := pdb.db.Exec(iQuery,
        sql.Named("cw", cryptword),
        sql.Named("id", userid),
        sql.Named("username", userid))
  if err == nil {
    return      // Succeeded
  }
  glog.Infof("INSERT returned err=%v\n", err)   // Expected if the user already exists.
  // If the INSERT failed, assume it was because the row already exists, so try updating it.
  query := "UPDATE user SET cryptword = :cw WHERE username = :username;"
  _, err = pdb.db.Exec(query,sql.Named("cw", cryptword),sql.Named("username", userid))
  if err != nil {
    glog.Errorf("Error setting cryptword for user %q: %v\n", userid, err)
  }
}

func (pdb *PwDB) UserCount() int {
  var count int
  sql := "SELECT count(*) from user;"
  err := pdb.db.QueryRow(sql).Scan(&count)
  if err != nil {
    glog.Errorf("Error counting user in database: %v\n", err)
    return 0
  }
  return count
}
