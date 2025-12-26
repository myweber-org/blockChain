package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type UserData struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type Aggregator struct {
	mu    sync.RWMutex
	users map[int]UserData
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		users: make(map[int]UserData),
	}
}

func (a *Aggregator) AddUser(user UserData) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.users[user.ID] = user
}

func (a *Aggregator) GetUser(id int) (UserData, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	user, exists := a.users[id]
	return user, exists
}

func (a *Aggregator) GetAllUsers() []UserData {
	a.mu.RLock()
	defer a.mu.RUnlock()
	users := make([]UserData, 0, len(a.users))
	for _, user := range a.users {
		users = append(users, user)
	}
	return users
}

func (a *Aggregator) AggregateFromSource(source func() []UserData) {
	users := source()
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(u UserData) {
			defer wg.Done()
			a.AddUser(u)
		}(user)
	}
	wg.Wait()
}

func mockUserSource() []UserData {
	return []UserData{
		{ID: 1, Name: "John Doe", Email: "john@example.com", CreatedAt: time.Now().Format(time.RFC3339)},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com", CreatedAt: time.Now().Add(-24 * time.Hour).Format(time.RFC3339)},
		{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", CreatedAt: time.Now().Add(-48 * time.Hour).Format(time.RFC3339)},
	}
}

func main() {
	aggregator := NewAggregator()
	aggregator.AggregateFromSource(mockUserSource)

	users := aggregator.GetAllUsers()
	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return
	}

	fmt.Println("Aggregated user data:")
	fmt.Println(string(jsonData))

	if user, exists := aggregator.GetUser(2); exists {
		fmt.Printf("\nRetrieved user ID 2: %s (%s)\n", user.Name, user.Email)
	}
}