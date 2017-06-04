package main

type UserModel struct {
    UserId   int
    Username string
    Pword    string
}

type GroupModel struct {
    GroupId    int      `json:"groupId"`
    GroupName  string   `json:"groupName"`
    InviteCode string   `json:"inviteCode"`
    CreatedBy  int      `json:"createdBy"`
}

type PersonModel struct {
    PersonId   int      `json:"personId"`
    GroupId    int      `json:"groupId"`
    PersonName string   `json:"personName"`
    Balance    float64  `json:"balance"`
}

type TransactionModel struct {
    TransId   int     `json:"transId,omitempty"`
    Title     string  `json:"title,omitempty"`
    Amount    float64 `json:"amount,omitempty"`
    TransDate string  `json:"transDate,omitempty"`
    GroupId   int     `json:"groupId,omitempty"`
    PaidBy    int     `json:"paidBy,omitempty"`
    SharedBy  []int   `json:"sharedBy,omitempty"`
}

// type TransactionModel struct {
//     TransId   int     `json:transId,omitempty`
//     Title     string  `json:title,omitempty`
//     Amount    float64 `json:amount,omitempty`
//     TransDate string  `json:transDate,omitempty`
//     GroupId   int     `json:groupId,omitempty`
//     PaidBy    int     `json:paidBy,omitempty`
//     SharedBy  []int   `json:sharedBy,omitempty`
// }
