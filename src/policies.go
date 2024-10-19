package src

func ProductPolicy(user User, data any) bool {
	// return user.Role == "admin"
	return true
}

func WalletCreatePolicy(user User, data any) bool {
	return user.ID != 0
}
func WalletReadPolicy(user User, data any) bool {
	return user.ID != 0
}
func WalletDetailPolicy(user User, data any) bool {
	con, ok := data.(Wallet)
	if !ok {
		return false
	}
	return user.ID == con.UserID && user.ID != 0
}

func CreateTransactionPolicy(user User, data any) bool {
	return user.ID != 0
}

func UpdateTransactionPolicy(user User, data any) bool {
	con, ok := data.(Transaction)
	if !ok {
		return false
	}
	return user.ID == con.UserID && user.ID == 0
}
