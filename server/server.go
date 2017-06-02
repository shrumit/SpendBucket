package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/dgrijalva/jwt-go"
     _ "github.com/go-sql-driver/mysql"
     "github.com/gorilla/schema"
    "log"
    "net/http"
    "strconv"
)

type ctxKey string
var dc DataController
var decoder = schema.NewDecoder()

func index(w http.ResponseWriter, r *http.Request) {
    if (r.URL.Path == "/") {
        http.ServeFile(w,r,"apiTest.html")
    } else {
        failureResponder(w,r,nil,"Nonexistent endpoint")
    }
}

func getCtxUid(r *http.Request) int {
    val, _ := r.Context().Value(ctxKey("uid")).(int)
    return val
}

func getCtxGid(r *http.Request) int {
    val, _ := r.Context().Value(ctxKey("gid")).(int)
    return val
}

// ----------------------------------

// Parse and validate JWT
func secureWrapper(h http.HandlerFunc) http.HandlerFunc {
    fn := func (w http.ResponseWriter, r *http.Request) {
        tokenString := r.PostFormValue("token")
        if tokenString == "" {
            failureResponder(w, r, nil, "No token supplied.")
            return
        }

        // parse and check token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
            }
            return HmacSecret, nil
        })
        var uid int
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            if v, ok := claims["uid"].(float64); ok {
                uid = int(v)
            } else {
                failureResponder(w, r, nil, "Problem parsing claim (type assertion).")
                return
            }
        } else {
                failureResponder(w, r, err, "Problem parsing claim.")
                return
        }
        fmt.Println("sW", uid)
        // store in Context
        key := ctxKey("uid")
        ctx := context.WithValue(r.Context(), key, uid)

        h(w, r.WithContext(ctx))
    }
    return fn
}

// Verify that user has access to the requested group
func groupWrapper(h http.HandlerFunc) http.HandlerFunc {
    fn := func (w http.ResponseWriter, r *http.Request) {
        groupId, err := strconv.Atoi(r.PostFormValue("groupId"))
        if err != nil {
            failureResponder(w,r,err,"Error parsing groupId.")
            return
        }
        fmt.Println("gW", getCtxUid(r), groupId)
        err = dc.CheckAccess(getCtxUid(r), groupId)
        if err != nil {
            failureResponder(w,r,err, "User doesn't have access to supplied Group.")
            return
        }

        // store in Context
        key := ctxKey("gid")
        ctx := context.WithValue(r.Context(), key, groupId)
        h(w, r.WithContext(ctx))
    }
    return fn
}

func main() {
    log.Println("log test")
    fmt.Println("fmt test")

    // Start database connection
    db, err := sql.Open("mysql", DSN)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer db.Close()
    err = db.Ping()
    if err != nil {
        log.Fatal(err.Error())
    }

    dc = &DataControllerMysql{db}

    // Bind API endpoints to middleware and handler functions

    http.HandleFunc("/",            index)
    http.HandleFunc("/login",       login)
    http.HandleFunc("/register",    register)

    // Token verified endpoints
    http.HandleFunc("/createGroup",         secureWrapper(createGroup))
    http.HandleFunc("/enterGroupInvite",    secureWrapper(enterGroupInvite))
    http.HandleFunc("/getGroups",           secureWrapper(getGroups))

    // Token and groupId verified endpoints
    http.HandleFunc("/getGroupData",        secureWrapper(groupWrapper(getGroupData)))
    http.HandleFunc("/getUsernames",        secureWrapper(groupWrapper(getUsernames)))
    http.HandleFunc("/addPerson",           secureWrapper(groupWrapper(addPerson)))
    http.HandleFunc("/getPersons",          secureWrapper(groupWrapper(getPersons)))
    http.HandleFunc("/getTransactions",     secureWrapper(groupWrapper(getTransactions)))
    http.HandleFunc("/addTransaction",      secureWrapper(groupWrapper(addTransaction)))
    http.HandleFunc("/deleteTransaction",   secureWrapper(groupWrapper(deleteTransaction)))

    log.Fatal(http.ListenAndServe(":8081", nil))
}

func login(w http.ResponseWriter, r *http.Request) {
    username := r.PostFormValue("username")
    pword := r.PostFormValue("password")
    fmt.Println(username, pword)
    m, err := dc.GetUser(username)
    if err == sql.ErrNoRows { 
        // username not found
        failureResponder(w,r,err,"Incorrect credentials.")
    } else if err != nil {
        failureResponder(w,r,err,"Server error.")
    } else if m.Pword != pword {
        // password doesn't match
        failureResponder(w,r,err,"Incorrect credentials.")
    } else {
        tokenResponder(w, r, m.UserId)
    }
}

func register(w http.ResponseWriter, r *http.Request) {
    model := &UserModel{Username: r.PostFormValue("username"), Pword: r.PostFormValue("password")}
    err := dc.CreateUser(model)
    if err != nil {
        failureResponder(w, r, err, "Registration error.")
        return
    }

    tokenResponder(w, r, model.UserId)
}

func tokenResponder(w http.ResponseWriter, r *http.Request, userId int) {
    fmt.Println("tokenResponder userId", userId)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
        "uid": userId,
    })
    tokenString, err := token.SignedString(HmacSecret)
    if err != nil {
        failureResponder(w, r, err, "Problem making token.")
        return
    }

    res := make(map[string]interface{})
    res["token"] = tokenString
    successResponder(w, r, res, "")
}

func createGroup(w http.ResponseWriter, r *http.Request) {
    var code string
    // generate random invite string and test for uniqueness
    for i := 0; i < 10; i++ {
        code = randomString(10)
        _,err := dc.GetGroupByInvite(code)
        if err == nil {
            break
        }
    }

    fmt.Println("code", code)
    fmt.Println("getCtxUid(r)", getCtxUid(r))

    m := GroupModel{GroupName: r.PostFormValue("groupName"), InviteCode: code}
    err := dc.CreateGroup(&m, getCtxUid(r))
    if err != nil {
        failureResponder(w, r, err, "Problem creating group.")
        return
    }

    res := make(map[string]interface{})
    res["group"] = &m
    successResponder(w, r, res, "")
}

func getGroups(w http.ResponseWriter, r *http.Request) {
    groups, err := dc.GetAllAccessGroups(getCtxUid(r))
    // ErrNoRows error is not a true failure, just return empty
    if err != nil && err != sql.ErrNoRows {
        failureResponder(w, r, err, "Problem.")
        return
    }

    res := make(map[string]interface{})
    res["groups"] = groups
    successResponder(w,r,res,"")
}

func enterGroupInvite(w http.ResponseWriter, r *http.Request) {
    m, err := dc.GetGroupByInvite(r.PostFormValue("inviteCode"))
    if err != nil {
        failureResponder(w,r,err, "Unable to enter group.")
        return
    }
    
    dc.CreateAccess(getCtxUid(r), m.GroupId)       

    res := make(map[string]interface{})
    res["group"] = m
    successResponder(w,r,res,"")
}

func getGroupData(w http.ResponseWriter, r *http.Request) {

    model, err := dc.GetGroupById(getCtxGid(r))
    if err != nil {
        failureResponder(w,r,err, "Unable to get group data.")
        return
    }

    res := make(map[string]interface{})
    res["group"] = model
    successResponder(w,r,res,"")
}

func getUsernames(w http.ResponseWriter, r *http.Request) {
    usernames, err := dc.GetAllAccessUsernames(getCtxGid(r))
    if err != nil {
        failureResponder(w,r,err, "Unable to get users.")
        return
    }

    res := make(map[string]interface{})
    res["usernames"] = usernames
    successResponder(w,r,res,"")    
}

func getPersons(w http.ResponseWriter, r *http.Request) {
    persons, err := dc.GetAllPersons(getCtxGid(r))
    if err != nil {
        failureResponder(w,r,err, "Unable to get persons.")
        return
    }

    res := make(map[string]interface{})
    res["persons"] = persons
    successResponder(w,r,res,"")
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
    transactions, err := dc.GetAllTransactions(getCtxGid(r))
    if err != nil {
        failureResponder(w,r,err, "Unable to get transactions.")
        return
    }

    res := make(map[string]interface{})
    res["transactions"] = &transactions
    successResponder(w,r,res,"")
}

func addTransaction(w http.ResponseWriter, r *http.Request) {

    // Parse string to required types in struct
    m := TransactionModel{}
    m.Title = r.PostFormValue("transaction[title]")
    m.TransDate = r.PostFormValue("transaction[transDate]")
    var err error
    m.Amount, err = strconv.ParseFloat(r.PostFormValue("transaction[amount]"), 64)
    if err != nil {
        failureResponder(w,r,err, "Unable to add transaction. Cannot parse amount.")
        return
    }
    m.PaidBy, err = strconv.Atoi(r.PostFormValue("transaction[paidBy]"))
    if err != nil {
        failureResponder(w,r,err, "Unable to add transaction. Cannot parse paidBy.")
        return
    }
    for _, val := range r.PostForm["transaction[sharedBy]"] {
        j, err := strconv.Atoi(val)
        if err != nil {
            failureResponder(w,r,err, "Unable to add transaction. Cannot parse sharedBy.")
            return
        }
        m.SharedBy = append(m.SharedBy, j)
    }
    if len(m.SharedBy) == 0 {
        failureResponder(w,r,err, "Unable to add transaction. Nothing supplied to sharedBy.")
        return
    }
    m.GroupId = getCtxGid(r)
    
    
    // Add transaction
    fmt.Println(m)
    err = dc.CreateTransaction(&m)
    if err != nil {
        failureResponder(w,r,err, "Unable to add transaction.")
        return
    }
    res := make(map[string]interface{})
    res["transaction"] = m
    successResponder(w,r,res,"") 
}

func addPerson(w http.ResponseWriter, r *http.Request) {
    m := PersonModel{PersonName: r.PostFormValue("personName"), GroupId: getCtxGid(r), Balance: 0}
    err := dc.CreatePerson(&m)
    if err != nil {
        failureResponder(w,r,err,"Unable to add person")
        return
    }
    res := make(map[string]interface{})
    res["person"] = m    
    successResponder(w,r,res,"")
}

func deleteTransaction(w http.ResponseWriter, r *http.Request) {
    transId, err := strconv.Atoi(r.PostFormValue("transId"))
    if err != nil {
        failureResponder(w,r,err, "Unable to delete transaction. Cannot parse transId.")
        return
    }

    err = dc.DeleteTransaction(transId, getCtxGid(r))
    if err != nil {
        failureResponder(w,r,err, "Unable to delete transaction.")
        return
    }

    successResponder(w,r,nil,"")
}




// predefined response handlers

func failureResponder(w http.ResponseWriter, r *http.Request, e error, msg string) {
    fmt.Println(r.URL.Path)
    w.Header().Set("Content-Type", "application/json")
    log.Printf("Error:%v;Path:%v;Msg:%v", e, r.URL.Path, msg)
    res := make (map[string]interface{})
    res["success"] = false
    res["message"] = msg

    err := json.NewEncoder(w).Encode(res)
    if err != nil {
        log.Println(err)
    }
}

func successResponder(w http.ResponseWriter, r *http.Request, res map[string]interface{}, msg string) {
    fmt.Println(r.URL.Path)
    w.Header().Set("Content-Type", "application/json")
    log.Printf("successResponder")
    if res == nil {
        res = make(map[string]interface{})
    }

    res["success"] = true
    res["message"] = msg

    err := json.NewEncoder(w).Encode(res)
    if err != nil {
        log.Println(err)
    }
}
