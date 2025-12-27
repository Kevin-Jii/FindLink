package router

var AdminAuthWhiteList = map[string]bool{
	"/ping":                               true,
	"/metrics":                            true,
	"/admin/v1/user/verify/captcha/check": true,
	"/admin/v1/user/verify/captcha":       true,
	"/admin/v1/user/login":                true,
	"/api/app/customer/v1/user/login":     true,
	"/api/app/customer/v1/user/register":  true,
}
