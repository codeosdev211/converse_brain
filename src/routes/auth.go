package routes

import (
	"fmt"
	"encoding/json"
	def "../definitions"
	helper "../helpers"
)

func sendErrResponse(res *http.ResponseWriter, message string) {
	json.NewEncoder(*res).Encode(defs.Response{1, message, nil})
}

func Authenticate(res http.ResponseWriter, req *http.Request) {
	// allowing only post requests on this method else returning bad request;
	if req.Method != "POST" {
		sendErrResponse(&res, "Bad request!")
		return
	}
	var request def.Request // struct to parse request body into;

	// validating request json object 
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrResponse(&res, "Invalid json body")
		return
	}

	reqType := string(request.Meta["reqType"])
	if reqType == "SI" {
		signIn(&res, &request.Meta)
	} else if  reqType == "SU" {
		signUp(&res, &request.Data[0])
	} else if reqType == "CP" {
		changePassword(&res, &request.Meta)
	} else if reqType == "FP" {
		recoverPassword(&res, &request.Meta)
	} else {
		sendErrResponse(&res, "Invalid Request")
		return
	}

}

func signIn(res *http.ResponseWriter, data *map[string]interface{}) {
	var query string
	var err error
	var result []map[string]interface{}

	query = fmt.Sprintf("select count(*) as isThere from CU "+
		"where email='%v' and password='%v';",
		data["email"], data["extra1"])
	result, err = helper.OnDatabase(true, &query)
	if err != nil {
		sendErrResponse(res, "Database error!")
		return
	}

	if result[0]["isThere"] == "0" {
		sendErrResponse(res, "Invalid Email ID or Password!")
		return
	}
	query = fmt.Sprintf("select * from CU where email='%v';",
		data["email"])
	result, err = helper.OnDatabase(true, &query)
	if err != nil {
		sendErrResponse(res, "Database error!")
		return
	}
	json.NewEncoder(*res).Encode(def.Response{
		0, "Success", result})
}

func signUp(res *http.ResponseWriter, data *map[string]interface{}) {
	var query string
	var err error

	query = fmt.Sprintf("select count(*) as isThere from CU " +
		"where email='%v';",
		data["email"])
	result, err = helper.OnDatabase(true, &query)
	if err != nil {
		sendErrResponse(res, "Database error!")
		return
	}

	if result[0]["isThere"] != "0" {
		sendErrResponse(res, "Email ID already registered!")
		return
	}

	query = fmt.Sprintf("insert into CU values (0, '%v', '%v', '%v',"+
		"'%v', 'active');",
		data["name"], data["email"], data["password"],
		data["createdOn"])
	result, err = helper.OnDatabase(false, &query)
	if err != nil {
		sendErrResponse(res, "Database error")
		return
	}
	json.NewEncoder(*res).Encode(def.Response{
		0, "Success", nil})
}


func changePassword(res *http.ResponseWriter, data *map[string]interface{}) {
	var query string
	var err error

	query = fmt.Sprintf("select count(*) as isThere from CU "+
		"where email='%v' and password='%v';",
		data["email"], data["extra1"])
	result, err = helper.OnDatabase(true, &query)
	if err != nil {
		sendErrResponse(res, "Database error!")
		return
	}

	if result[0]["isThere"] == "0" {
		sendErrResponse(res, "Invalid old password")
		return
	}

	query = fmt.Sprintf("Update CU set password='%v' where email='%v';"+
		data["extra2"], data["email"])
	result, err = helper.OnDatabase(false, &query)
	if err != nil {
		sendErrResponse(res, "Database error")
		return
	}
	json.NewEncoder(*res).Encode(def.Response{
		0, "Success", nil})
}

func recoverPassword(res *http.ResponseWriter, data *map[string]interface{}) {
	var query string
	var err error

	query = fmt.Sprintf("select count(*) as isThere from CU " +
		"where email='%v';",
		data["email"])
	result, err = helper.OnDatabase(true, &query)
	if err != nil {
		sendErrResponse(res, "Database error!")
		return
	}

	if result[0]["isThere"] == "0" {
		sendErrResponse(res, "Invalid Email address!")
		return
	}

	code := helper.getRandomStr()

	err = helper.sendPlainEmail(data["email"], code)
	if err != nil {
		sendErrResponse(res, "Could not send confirmation code!")
		return
	}
	var values map[string]interface{}
	values["code"] = code
	json.NewEncoder(*res).Encode(def.Response{
		0, "Success", []map[string]interface{
			values
		}
	})
}





