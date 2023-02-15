package auth

import (
  "database/sql"

  "github.com/jimmc/jraceman/dbrepo/conn"

  "github.com/golang/glog"

  "github.com/jimmc/auth/permissions"
  "github.com/jimmc/auth/users"
)

// PwDB implements the Store interface to load and store data in our SQL database.
// Data is stored in a table called "user" with three string columns,
// id, saltword, and permissions,
// where the permissions value is a comma-separated list of permission names.
type PwDB struct {
    db conn.DB
}

func NewPwDB(db conn.DB) *PwDB {
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

func (pdb *PwDB) User(username string) *users.User {
  glog.V(2).Infof("Looking for user %s", username)
  query := `
    SELECT saltword, GROUP_CONCAT(COALESCE(permission.name,''),' ') as permissions
    FROM user
         LEFT JOIN userrole ON user.id=userrole.userid
         LEFT JOIN rolepermission ON userrole.roleid=rolepermission.roleid
         LEFT JOIN permission ON rolepermission.permissionid=permission.id
    WHERE user.username = :username
    GROUP BY user.id`
  var saltword string
  var perms string
  err := pdb.db.QueryRow(query,sql.Named("username", username)).Scan(&saltword, &perms)
  if err == sql.ErrNoRows {
    glog.Warningf("No matching rows found for user %s", username)
    return nil          // No matching username found
  }
  if err != nil {
    glog.Errorf("Error scanning for user %q: %v\n", username, err)
    return nil
  }
  user := users.NewUser(username, saltword, permissions.FromString(perms))
  glog.V(2).Infof("Found user %v with perms %v", user, perms)
  return user
}

func (pdb *PwDB) SetSaltword(username, saltword string) {
  if username == "" {
    glog.Errorf("Can't SetSaltword with no username\n")
    return
  }
  // Assume row does not exist, try to insert it.
  iQuery := `INSERT into user(id,username,saltword) values(:id, :username, :cw);`
  _, err := pdb.db.Exec(iQuery,
        sql.Named("cw", saltword),
        sql.Named("id", username),
        sql.Named("username", username))
  if err == nil {
    return      // Succeeded
  }
  glog.Infof("INSERT returned err=%v; will try UPDATE\n", err)   // Expected if the user already exists.
  // If the INSERT failed, assume it was because the row already exists, so try updating it.
  query := "UPDATE user SET saltword = :cw WHERE username = :username;"
  _, err = pdb.db.Exec(query,sql.Named("cw", saltword),sql.Named("username", username))
  if err != nil {
    glog.Errorf("Error setting saltword for user %q: %v\n", username, err)
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
