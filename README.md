# simplehydraflow
## step 1

***Install ory cli by following the steps given in the link below***
[ory cli installation docs-- please click this link to install the ory cli](https://www.ory.sh/docs/guides/cli/installation)


## step 2
***once the ory cli is installed create a ory project by running the below line***
*`ory create project`*

## step 3
***once the project is created run the below line to get the project ID***
*`ory list projects`*

## step 4
***please note down the project ID you got from the above step***

***now run the below line to create a oauth client in Hydra , use the Project ID that you got from the Previous step***

    ory create oauth2-client --project <Your Project ID> \
    
    --name "nodeapp" \
    
    --grant-type authorization_code,refresh_token,client_credentials \
    
    --response-type code,id_token \
    
    --scope openid --scope offline_access --scope email \
    
    --redirect-uri http://127.0.0.1:5555/callback

***The response will be similar to the one provided below, note down all those details***

    CLIENT ID 505409xa-8ccd-4259-9444-c4f7281d8de7
    
    CLIENT SECRET cSYiFwbbOlqrgkjashsdf~ksGDt
    
    GRANT TYPES authorization_code, refresh_token, client_credentials
    
    RESPONSE TYPES code, id_token
    
    SCOPE openid offline_access email
    
    AUDIENCE
    
    REDIRECT URIS http://127.0.0.1:5555/callback

## step 5

***Execute the below line  to get the config file***

    ory get oauth2-config <Your project ID> --format yaml > config.yaml
## step 6 

***modify the config as given below: ***

    urls:
    
    consent: http://localhost:3000/consent
    
    error: /ui/error
    
    login: http://localhost:3000/login
    
    post_logout_redirect: /oauth2/fallbacks/logout
    
   ## step 7
   ***after modifying the config files login and consent endpoints as given above execute the below line***

    ory update oauth2-config <Your project ID> --file config.yaml

## step 8

***Run the go program which will serve the login and consent endpoints***

    go run main.go

## step 9

***now execute the below lines , this will create a oauth client application that will run in http://localhost:5555***
***This application will mimic the oauth client application in real world.***
***Ory has created this awsome solution for the developers to test everything***

***

    ory perform authorization-code \
      --client-id <Your client ID> \
      --client-secret <Your client secret> \
      --project <your Project ID> \
      --port 5555 \
      --scope openid

***

## step 10

***now open http://localhost:5555***
***voila you got the id token ***
[https://jwt.io/](https://jwt.io/) use this link to decode the id token
