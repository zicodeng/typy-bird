package models

import (
	"fmt"
	"testing"
	"gopkg.in/mgo.v2"
)



func TestMongoStore(t *testing.T) {

	newTypie1 := &TypieBird {
		UserName     : "eric",
		
	}

	// newTypie2 := &TypieBird{
	// 	UserName     : "Zico",
	// }

	// newTypie3 := &TypieBird{}

	session, _ := mgo.Dial("192.168.99.100");
	store := NewMongoStore(session, "db", "typiebirds");
	typie, err := store.InsertTypieBird(newTypie1)
	if err != nil {
		t.Errorf("error getting typie bird: %v", err)
	}
	fmt.Println(newTypie1)
	typie, err = store.GetByUserName("eric")
	fmt.Println(typie)
	if err != nil {
		t.Errorf("error getting typie bird: %v", err)
	}
	if typie.UserName != "eric" {
		t.Errorf("retrieving wrong typie bird")
	}

}
