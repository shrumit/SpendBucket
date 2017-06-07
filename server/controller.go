package main

// This interface allows the server to access the underlying database.

// Create and Update functions take a *Model as input 
// Create functions update the input *Model's index field upon successful insertion
// Get functions take one identifying parameter and return *Model

type DataController interface {
    // User
    CreateUser(model *UserModel) error
    GetUser(username string) (*UserModel, error)

    // Group metadata
    CreateGroup(model *GroupModel, userId int) error // also update access
    GetGroupById(groupId int) (*GroupModel, error)
    GetGroupByInvite(inviteCode string) (*GroupModel, error)
    UpdateGroup(model *GroupModel) error

    // User-Group access
    CreateAccess(userId int, groupId int) error
    CheckAccess(userId int, groupId int) error // nil error is success
    GetAllAccessUsernames(groupId int) ([]string, error)
    GetAllAccessGroups(userId int) ([]*GroupModel, error)

    // Group transaction data
    CreatePerson(model *PersonModel) error
    GetAllPersons(groupId int) ([]*PersonModel, error)
    CreateTransaction(model *TransactionModel) error
    GetAllTransactions(groupId int) ([]*TransactionModel, error)
    DeleteTransaction(transId int, groupId int) error // groupId required for security
}