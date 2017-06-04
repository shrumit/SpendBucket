package main

// SQL implementation of the DataController interface,
// written with the intention of being used with a MySQL driver.

import (
    "database/sql"
    "errors"
    // "fmt"
)

// Initialized with a sql.DB driver
type DataControllerMysql struct {
    db *sql.DB
}

// Users

func (dc *DataControllerMysql) CreateUser(model *UserModel) (err error) {
    result, err := dc.db.Exec("insert into Users (username, pword) values (?,?)", model.Username, model.Pword)
    var temp int64
    if err == nil {
        temp, err = result.LastInsertId()
        model.UserId = int(temp)
    }
    return err
}

func (dc *DataControllerMysql) GetUser(username string) (model *UserModel, err error) {
    model = &UserModel{}
    err = dc.db.QueryRow("select * from Users where username=?", username).Scan(&model.UserId, &model.Username, &model.Pword)
    return
}

// Group

func (dc *DataControllerMysql) CreateGroup(model *GroupModel, userId int) (err error) {
    tx, err := dc.db.Begin()
    if err != nil {
        return
    }
    defer tx.Rollback()

    res, err := tx.Exec("insert into Groups (inviteCode, groupName, createdBy) values (?,?,?)", model.InviteCode, model.GroupName, userId)
    if err != nil {
        return
    }

    var temp int64
    temp, err = res.LastInsertId()
    if err != nil {
        return
    }
    model.GroupId = int(temp)

    _, err = tx.Exec("insert into CanAccess values (?,?)", userId, model.GroupId)
    if err != nil {
        return
    }
    tx.Commit()
    return
}

func (dc *DataControllerMysql) GetGroupById(groupId int) (model *GroupModel, err error) {
    model = &GroupModel{}
    err = dc.db.QueryRow("select * from Groups where groupId=?", groupId).Scan(&model.GroupId, &model.InviteCode, &model.GroupName)
    return
}

func (dc *DataControllerMysql) GetGroupByInvite(inviteCode string) (model *GroupModel, err error) {
    model = &GroupModel{}
    err = dc.db.QueryRow("select * from Groups where inviteCode=?", inviteCode).Scan(&model.GroupId, &model.GroupName, &model.InviteCode, &model.CreatedBy)
    return
}

func (dc *DataControllerMysql) UpdateGroup(model *GroupModel) (err error) {
    _, err = dc.db.Exec("update Groups set inviteCode=?, groupName=? where groupId=?", model.InviteCode, model.GroupName, model.GroupId)
    return
}

// User-Group Access

func (dc *DataControllerMysql) CreateAccess(userId int, groupId int) (err error) {
    _, err = dc.db.Exec("insert into CanAccess values (?,?)", userId, groupId)
    return
}

func (dc *DataControllerMysql) CheckAccess(userId int, groupId int) (err error) {
    var dummy int
    err = dc.db.QueryRow("select 1 from CanAccess where userId=? and groupId=?", userId, groupId).Scan(&dummy)
    return
}

func (dc *DataControllerMysql) GetAllAccessGroups(userId int) (models []*GroupModel, err error) {
    rows, err := dc.db.Query("select g.GroupId, inviteCode, groupName from CanAccess c join Groups g on c.GroupId=g.GroupId where userId=?", userId)
    if err != nil {
        return
    }
    for rows.Next() {
        m := GroupModel{}
        err = rows.Scan(&m.GroupId, &m.InviteCode, &m.GroupName)
        if err != nil {
            return nil, err
        }
        models = append(models, &m)
    }
    return
}

func (dc *DataControllerMysql) GetAllAccessUsernames(groupId int) (usernames []string, err error) {
    rows, err := dc.db.Query("select username from CanAccess c join Users u on c.UserId=u.UserId where groupId=?", groupId)
    if err != nil {
        return
    }
    for rows.Next() {
        var username string
        err = rows.Scan(&username)
        if err != nil {
            return nil, err
        }
        usernames = append(usernames, username)
    }
    return
}

// Person

func (dc *DataControllerMysql) CreatePerson(model *PersonModel) (err error) {
    result, err := dc.db.Exec("insert into Persons (groupId, personName) values (?,?)", model.GroupId, model.PersonName)
    var temp int64
    if err == nil {
        temp, err = result.LastInsertId()
        model.PersonId = int(temp)
    }
    return
}

func (dc *DataControllerMysql) GetAllPersons(groupId int) (persons []*PersonModel, err error) {
    rows, err := dc.db.Query("select * from Persons where groupId=?", groupId)
    if err != nil {
        return
    }
    for rows.Next() {
        m := PersonModel{}
        err = rows.Scan(&m.PersonId, &m.GroupId, &m.PersonName)
        if err != nil {
            return
        }

        // get Balance
        err = dc.db.QueryRow("select amount from Balance where personId=?", m.PersonId).Scan(&m.Balance)
        if err == sql.ErrNoRows { // no rows means balance=0
            m.Balance = 0
            err = nil
        } else if err != nil {
            return nil, err
        }

        persons = append(persons, &m)
    }
    return
}

// Transaction

func (dc *DataControllerMysql) CreateTransaction(model *TransactionModel) (err error) {

    tx, err := dc.db.Begin()
    if err != nil {
        return
    }
    defer tx.Rollback()

    // update Transactions table
    res, err := tx.Exec("insert into Transactions (title, amount, transDate, groupId, paidBy) values (?,?,?,?,?)",
        model.Title,
        model.Amount,
        model.TransDate,
        model.GroupId,
        model.PaidBy,
    )

    if err != nil {
        tx.Rollback()
        return
    }

    // retrieve transId
    var temp int64
    temp, err = res.LastInsertId()
    if err != nil {
        tx.Rollback()
        return
    }
    model.TransId = int(temp)

    // update SharedBy table
    for _, personId := range model.SharedBy {
        _, err = tx.Exec("insert into SharedBy (transId, personId) values (?,?)", model.TransId, personId)
        if err != nil {
            tx.Rollback()
            return
        }
    }
    tx.Commit()
    return
}

func (dc *DataControllerMysql) GetAllTransactions(groupId int) (transactions []*TransactionModel, err error) {
    rows, err := dc.db.Query("select * from Transactions where groupId=? order by transDate desc", groupId)
    if err != nil {
        return
    }
    for rows.Next() {
        m := TransactionModel{}

        err = rows.Scan(&m.TransId, &m.Title, &m.Amount, &m.TransDate, &m.GroupId, &m.PaidBy)
        if err != nil {
            return nil, err
        }

        rowsShared, err := dc.db.Query("select personId from SharedBy where transId=?", m.TransId)
        if err != nil {
            return nil, err
        }
        for rowsShared.Next() {
            var p int
            err = rowsShared.Scan(&p)
            if err != nil {
                return nil, err
            }
            m.SharedBy = append(m.SharedBy, p)
        }

        transactions = append(transactions, &m)
    }
    return
}

func (dc *DataControllerMysql) DeleteTransaction(transId int, groupId int) (err error) {
    res, err := dc.db.Exec("delete from Transactions where transId=? and groupId=?", transId, groupId)
    rowCnt, err := res.RowsAffected()
    if err != nil {
        return
    } else if rowCnt == 0 {
        err = errors.New("No such transaction found to delete.")
    }
    return
}