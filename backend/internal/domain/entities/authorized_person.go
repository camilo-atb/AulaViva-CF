package entities

import (
	"strings"
	"unicode"
)

type DocumentType string

const (
	DocumentCC       DocumentType = "CC"
	DocumentTI       DocumentType = "TI"
	DocumentCE       DocumentType = "CE"
	DocumentPassport DocumentType = "PASSPORT"
	DocumentOther    DocumentType = "OTHER"
)

type AuthorizedPerson struct {
	ID             int64        `json:"id"`
	StudentID      int64        `json:"studentId"`
	FullName       string       `json:"fullName"`
	DocumentType   DocumentType `json:"documentType"`
	DocumentNumber string       `json:"documentNumber"`
	Relationship   string       `json:"relationship"`
	Phone          string       `json:"phone"`
}

func (p AuthorizedPerson) Validate() bool {
	if p.StudentID <= 0 || strings.TrimSpace(p.FullName) == "" ||
		strings.TrimSpace(p.DocumentNumber) == "" || strings.TrimSpace(p.Relationship) == "" ||
		!validDocumentType(p.DocumentType) {
		return false
	}
	return validPhone(p.Phone)
}

func validDocumentType(documentType DocumentType) bool {
	switch documentType {
	case DocumentCC, DocumentTI, DocumentCE, DocumentPassport, DocumentOther:
		return true
	default:
		return false
	}
}

func validPhone(phone string) bool {
	trimmed := strings.TrimSpace(phone)
	digits := 0
	for index, char := range trimmed {
		if unicode.IsDigit(char) {
			digits++
			continue
		}
		if (char == '+' && index == 0) || char == ' ' || char == '-' || char == '(' || char == ')' {
			continue
		}
		return false
	}
	return digits >= 7 && digits <= 15
}

func NormalizeAuthorizedPerson(person AuthorizedPerson) AuthorizedPerson {
	person.FullName = strings.Join(strings.Fields(person.FullName), " ")
	person.DocumentNumber = strings.Join(strings.Fields(person.DocumentNumber), " ")
	person.Relationship = strings.Join(strings.Fields(person.Relationship), " ")
	person.Phone = strings.TrimSpace(person.Phone)
	return person
}
