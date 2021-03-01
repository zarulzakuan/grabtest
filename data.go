package main

type UserToken struct {
	UID      string
	UserName string
}

// HTTPResponse2 Http response
type HTTPResponse struct {
	Status string      `json:"status"` // success, fail, error
	Data   interface{} `json:"data"`
}

type User struct {
	ID    string `firestore:"id" json:"id" swaggerignore:"true"`
	Email string `firestore:"email" json:"email"`
	Name  string `firestore:"name" json:"name"`
}

type Job struct {
	ID          string   `firestore:"id" json:"id" swaggerignore:"true"`
	Title       string   `firestore:"title" json:"title"`
	Description string   `firestore:"description" json:"description"`
	SalaryMin   string   `firestore:"salarymin" json:"salarymin"`
	SalaryMax   string   `firestore:"salarymax" json:"salarymax"`
	SearchKeys  []string `firestore:"searchkeys" json:"searchkeys"`
}

type Employee struct {
	ID         string   `firestore:"id" json:"id" swaggerignore:"true"`
	Name       string   `firestore:"name" json:"name"`
	DOB        string   `firestore:"dob" json:"dob"`
	JobTitle   string   `firestore:"jobtitle" json:"jobtitle"`
	Salary     string   `firestore:"salary" json:"salary"`
	SearchKeys []string `firestore:"searchkeys" json:"searchkeys"`
}

type Report struct {
	ID         string `firestore:"id" json:"id" swaggerignore:"true"`
	UserID     string `firestore:"userid" json:"userid"`
	PageName   string `firestore:"pagename" json:"page"`
	AccessFreq string `firestore:"accessfreq" json:"accessfreq"`
}

type RecordLog struct {
	ID       string `firestore:"id" json:"id" swaggerignore:"true"`
	UserID   string `firestore:"userid" json:"userid"`
	UserName string `firestore:"username" json:"username"`
	Type     string `firestore:"type" json:"type" example:"login|logout|page visit"`
	PageName string `firestore:"pagename" json:"pagename"`
	DateTime string `firestore:"datetime" json:"datetime"`
}
