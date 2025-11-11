package tools

type mockDB struct{}

var mockUserDetails = []UserDetail{
	{ID: 1, Name: "salil", Email: "salil@contactspace.com"},
	{ID: 2, Name: "deepa", Email: "deepamnc12@gmail.com"},
	{ID: 3, Name: "john", Email: "john@gmail.com"},
}

func (db *mockDB) ValidateToken(token string) (*Token, error) {
	if token == "Bearer secrettoken" {
		return &Token{Token: token}, nil
	}
	return nil, nil
}
func (db *mockDB) SetupDatabase() error {
	// Mock setup, do nothing
	return nil
}

func (db *mockDB) GetAllUsers() (*[]UserDetail, error) {
	return &mockUserDetails, nil
}

func (db *mockDB) AddUser(user *UserDetail) (*[]UserDetail, error) {
	mockUserDetails = append(mockUserDetails, *user)
	return &mockUserDetails, nil
}

func (db *mockDB) UpdateUser(updatedUser UserDetail, id int) (*[]UserDetail, error) {
	for i, user := range mockUserDetails {
		if user.ID == id {
			mockUserDetails[i] = updatedUser
			break
		}
	}
	return &mockUserDetails, nil
}

func (db *mockDB) DeleteUser(id int) (*[]UserDetail, error) {
	for i, user := range mockUserDetails {
		if user.ID == id {
			mockUserDetails = append(mockUserDetails[:i], mockUserDetails[i+1:]...)
			break
		}
	}
	return &mockUserDetails, nil
}
