package users_transport_http

import "github.com/AMmetro/ODRProcesing/internal/core/domain"

type UserDTOResponse struct {
	ID          int    `json:"id" example:"10"`
	Version     int    `json:"version" example:"3"`
	FullName    string `json:"full_name" example:"John Doe"`
	PhoneNumber string `json:"phone_number" example:"+1234567890"`
}

func UserDtoFromDomain(userAgent domain.UserAgent) UserDTOResponse {
	var phone string
	if userAgent.PhoneNumber != nil {
		phone = *userAgent.PhoneNumber
	}
	return UserDTOResponse{
		ID:          userAgent.ID,
		Version:     userAgent.Version,
		FullName:    userAgent.FullName,
		PhoneNumber: phone,
	}
}

func UsersDtoFromDomains(users []domain.UserAgent) []UserDTOResponse {
	res := make([]UserDTOResponse, 0, len(users))
	for _, userAgent := range users {
		dto := UserDtoFromDomain(userAgent)
		res = append(res, dto)
	}
	return res
}
