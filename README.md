# restclient

## Example Usage

```golang
type AuthenticationStruct struct {
		Username string
		Password string
}

BasicAuthenticationMethod := func(req *http.Request, v ...interface{}) (*http.Request, error) {
  auth := v[0].(AuthenticationStruct)
  code := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", auth.Username, auth.Password)))
  req.Header.Add("Authorization", fmt.Sprintf("Basic %s", code))
  return req, nil
}

config := NewConfig().
		WithEndpoint("https://api.trello.com/1/").
		SetAuthenticationMethod(BasicAuthenticationMethod)

operation, err := config.NewOperation(http.MethodPost)

if err != nil {
    return err
}

req, err := operation.
		WithPath("boards/{boardid}/cards/{cardid}").
		WithPathVar("boardid", "12345").
		WithPathVar("cardid", "12345").
		BodyFromJSONString(`{"name": "Seymour Butts"}`).
		Authenticate(AuthenticationStruct{Username: "john", Password: "random1234"}).
		BuildRequest()

if err != nil {
    return err
}

// Don't actually use the DefaultClient in production...
http.DefaultClient.Do(req)
```
