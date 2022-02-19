package utilities

import (
	"errors"
	"net/http"

	"github.com/ddld93/auth/model"
)

func UserModelValidate(user *model.User)  (*model.User, error){
	// checking required fields 
	if user.FullName == "" {
		return user, errors.New("full name field cannot be empty")
	}
	if user.Email == "" {
		return user, errors.New("email field cannot be empty")
	}
	if user.Password == "" {
		return user, errors.New("password field cannot be empty")
	}
	
	if user.Phone == "" {
		return user, errors.New("phone field cannot be empty")
	}
	
	// assigning default value

	user.Role = "client"
	return user, nil
}
func VerifyPayment(ref string) error{
	apiKey:= "sk_test_91dcbd0fc948c4670f12b9384402e87a56927c27"

	
	url := "https://api.paystack.co/transaction/verify/" + ref


    // Create a Bearer string by appending string access token
    var bearer = "Bearer " + apiKey

    // Create a new request using http
    req, _ := http.NewRequest("GET", url, nil)

    // add authorization header to the req
    req.Header.Add("Authorization", bearer)

    // Send req using http Client
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        
		return errors.New("error making request to paystack")
    }
    defer resp.Body.Close()
	if resp.StatusCode != 200{
		
		return errors.New("payment not verified")
	}
	
	return nil
}
