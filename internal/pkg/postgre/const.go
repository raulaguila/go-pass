package postgre

const (
	Profile     string = "Profile"
	Permissions string = "Permissions"
	Operator    string = "Operator"
	Accounts    string = "Accounts"
	Site        string = "Site"
	Mail        string = "Mail"
	Phone       string = "Phone"

	PhoneOperator     string = Phone + "." + Operator
	ProfilePermission string = Profile + "." + Permissions
)
