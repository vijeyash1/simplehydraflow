package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AlekSi/pointer"
	"github.com/gin-gonic/gin"
	ory "github.com/ory/client-go"
)

func main() {
	r := gin.Default()
	handler := NewHandler()
	r.GET("/createclient", handler.CreateOauthClientHandler)
	r.GET("/login", handler.AcceptLoginHandler)
	r.GET("/consent", handler.AcceptConsentHandler)
	r.Run(":3000")
}

type Handler struct {
	ApiClient *ory.APIClient
	ctx       context.Context
	logger    *log.Logger
}

func NewHandler() *Handler {
	conf := ory.NewConfiguration()
	conf.Servers = ory.ServerConfigurations{{
		URL: "https://gauss-2yivnr1dcc.projects.oryapis.com",
	}}
	logger := log.New(os.Stderr, "app: ", log.LstdFlags|log.Lshortfile|log.Llongfile)
	oryAuthedContext := context.WithValue(context.Background(), ory.ContextAccessToken, "ory_pat_ACe9AwgTYTz1HoO08CyoVHW")
	return &Handler{
		ApiClient: ory.NewAPIClient(conf),
		ctx:       oryAuthedContext,
		logger:    logger,
	}
}

// CreateOauthClientHandler creates a new Oauth2.0 client
func (h *Handler) CreateOauthClientHandler(c *gin.Context) {
	log.Println(h.ctx.Value(ory.ContextAccessToken))
	_, res, err := h.ApiClient.OAuth2Api.CreateOAuth2Client(h.ctx).
		OAuth2Client(ory.OAuth2Client{
			ClientName:              pointer.ToString("auth-client-test"),
			ClientSecret:            pointer.ToString("secret"),
			GrantTypes:              []string{"authorization_code", "refresh_token"},
			RedirectUris:            []string{"http://localhost:3000/callback"},
			ResponseTypes:           []string{"code", "id_token"},
			Scope:                   pointer.ToString("openid offline"),
			TokenEndpointAuthMethod: pointer.ToString("client_secret_post"),
		}).Execute()
	defer res.Body.Close()
	if err != nil {
		h.logger.Println(res)
		c.JSON(500, gin.H{
			"response" : res,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Oauth2.0 client created successfully",
	})
}

// func (h *Handler) GetOauthClientHandler(c *gin.Context) {
// 	clientId := os.Getenv("CLIENT_ID")
// 	client, _, err := h.ApiClient.OAuth2Api.GetOAuth2Client(h.ctx, clientId).Execute()
// 	if err != nil {
// 		h.logger.Println(err.Error())
// 		c.JSON(500, gin.H{
// 			"message": "Error getting Oauth2.0 client",
// 		})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"message": client,
// 	})
// }

// func (h *Handler) DeleteOauthClientHandler(c *gin.Context) {
// 	clientId := os.Getenv("CLIENT_ID")
// 	_, err := h.ApiClient.OAuth2Api.DeleteOAuth2Client(h.ctx, clientId).Execute()
// 	if err != nil {
// 		h.logger.Println(err.Error())
// 		c.JSON(500, gin.H{
// 			"message": "Error deleting Oauth2.0 client",
// 		})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"message": "Oauth2.0 client deleted successfully",
// 	})
// }



func (h *Handler) AcceptLoginHandler(c *gin.Context) {
	challenge := c.Query("login_challenge")
	accept, _, err := h.ApiClient.OAuth2Api.AcceptOAuth2LoginRequest(h.ctx).LoginChallenge(challenge).AcceptOAuth2LoginRequest(ory.AcceptOAuth2LoginRequest{
		Remember:    pointer.ToBool(true),
		RememberFor: pointer.ToInt64(3600),
		Subject:     "vijey@gmail.com",
	}).Execute()

	if err != nil {
		fmt.Println("error", err.Error())
	}
	c.Redirect(http.StatusFound, accept.RedirectTo)
}

func (h *Handler) AcceptConsentHandler(c *gin.Context) {
	challenge := c.Query("consent_challenge")
	fmt.Println("consent_challenge", challenge)
	acceptConsentRes, _, err := h.ApiClient.OAuth2Api.AcceptOAuth2ConsentRequest(h.ctx).
		ConsentChallenge(challenge).AcceptOAuth2ConsentRequest(
		ory.AcceptOAuth2ConsentRequest{
			GrantScope:  []string{"openid"},
			Remember:    pointer.ToBool(true),
			RememberFor: pointer.ToInt64(3600),
			Session: &ory.AcceptOAuth2ConsentRequestSession{
				IdToken: PersonSchemaJsonTraits{Email: "vijesh@gmail.com", Name: &PersonSchemaJsonTraitsName{First: pointer.ToString("vijesh"), Last: pointer.ToString("kumar")}},
			},
		}).Execute()
	if err != nil {
		fmt.Println("error", err.Error())
	}
	c.Redirect(http.StatusFound, acceptConsentRes.RedirectTo)
}

type PersonSchemaJsonTraits struct {
	// Email corresponds to the JSON schema field "email".
	Email string `json:"email" yaml:"email"`

	// Name corresponds to the JSON schema field "name".
	Name *PersonSchemaJsonTraitsName `json:"name,omitempty" yaml:"name,omitempty"`
}

type PersonSchemaJsonTraitsName struct {
	// First corresponds to the JSON schema field "first".
	First *string `json:"first,omitempty" yaml:"first,omitempty"`

	// Last corresponds to the JSON schema field "last".
	Last *string `json:"last,omitempty" yaml:"last,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PersonSchemaJsonTraits) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["email"]; !ok || v == nil {
		return fmt.Errorf("field email in PersonSchemaJsonTraits: required")
	}
	type Plain PersonSchemaJsonTraits
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PersonSchemaJsonTraits(plain)
	return nil
}

type PersonSchemaJson struct {
	// Traits corresponds to the JSON schema field "traits".
	Traits *PersonSchemaJsonTraits `json:"traits,omitempty" yaml:"traits,omitempty"`
}
