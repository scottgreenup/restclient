# restclient

## Example Usage

```golang
package main

import (
    "encoding/base64"
    "fmt"
    "github.com/scottgreenup/restclient"
)

func main() {
    type ExampleRequest struct {
        Version uint64 `json:"number"`
        Hash    string `json:"hash"`
    }

    requestData := &ExampleRequest{
        Version: 1,
        Hash:    "1c76f1f33f12c14a63026f71c8d17ab2",
    }

    rm := restclient.NewRequestMutator(
        restclient.BaseURL("https://scottgreenup.com/"),
    )

    BasicAuthMutator := func(username, password string) restclient.RequestMutation {
        return func(req *http.Request) error {
            code := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
            req.Header.Add("Authorization", fmt.Sprintf("Basic %s", code))
            return nil
        }
    }

    req, _ := rm.NewRequest(
        restclient.ResolvePath("/api/example"),
        restclient.Method(http.MethodPost),
        restclient.BasicAuthMutator("MyUsername", "4RuwRmDkLm990qkXMK6obWK88S7pW3K3"),
        restclient.BodyFromJSON(requestData),
    )

    fmt.Println("URL:", req.URL.String())
    fmt.Println("Content Length:", req.ContentLength)
    fmt.Println("Authorization Header:", req.Header.Get("Authorization"))
}
```
