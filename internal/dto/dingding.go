package dto

type DingdingBase struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type DingdingResult struct {
	DingdingBase
	Result interface{}
}

type DingdingGetTokenResult struct {
	DingdingBase
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type DingdingDept struct {
	DeptId          int64  `json:"dept_id"`
	ParentId        int64  `json:"parent_id"`
	Name            string `json:"name"`
	AutoAddUser     bool   `json:"auto_add_user"`
	CreateDeptGroup bool   `json:"create_dept_group"`
}

type DingdingDeptListResult struct {
	DingdingBase
	Result []DingdingDept `json:"result"`
}

type DingdingDeptUser struct {
	Name   string `json:"name"`
	UserId string `json:"userid"`
}

type DingdingDeptUserListResult struct {
	DingdingBase
	Result struct {
		HasMore bool               `json:"has_more"`
		List    []DingdingDeptUser `json:"list"`
	}
}

type DingdingGetUserResult struct {
	DingdingBase
	Result struct {
		Email  string `json:"email"`
		Title  string `json:"title"`
		Mobile string `json:"mobile"`
	}
}

type DingdingSendAppMessageResult struct {
	DingdingBase
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}
