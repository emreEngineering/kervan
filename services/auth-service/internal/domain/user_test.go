package domain

import "testing"

func TestNewUser_Success(t *testing.T) {
	user, err := NewUser("test@example.com", "hashedpasswword123")

	if err != nil {
		t.Fatalf("hata beklenmiyordu: %v", err)
	}

	if user == nil {
		t.Fatalf("user nil olamaz")
	}

	if user.Email != "test@example.com" {
		t.Errorf("email= %s, beklenen test@example.com", user.Email)
	}
}

func TestNewUser_EmptyEmail(t *testing.T) {
	_, err := NewUser("", "hashedpassword123")
	if err == nil {
		t.Fatal("boş mail için hata bekleniyordu")
	}
}

func TestNewUser_InvalidEmail(t *testing.T) {
	_, err := NewUser("bu-bir-email-değil", "hashedpassword123")

	if err == nil {
		t.Fatal("geçersiz email için hata bekleniyordu")
	}
}

func TestNewUser_EmptyPasswordHash(t *testing.T) {
	_, err := NewUser("test@example.com", "")

	if err == nil {
		t.Fatal("boş şifre hash'i için hata bekleniyordu")
	}
}
