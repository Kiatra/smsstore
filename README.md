# smsstore
Minimalistic server written in go that provides REST endpoints to receive messages as a json paylod as 
defined by the android sms forwarding app "SMS an PC/Telefon - Automatische Umleitung" and
expose the recived message. 

Supports simplisitc user and autentification support.

## API Specifiction
### POST /messages Stores a message for a user
Required parameters: user, pass

Example Payload:
`{"subject": "foo", "message": "bar"}`

Test API: 
```
curl -X POST -H 'Content-Type: application/json' -d '{"subject": "foo","message": "bar"}' https://localhost/messages/?user=test&pass=1234
``` 

### POST /messages Provides last recived message (without subject) for a speified user as text
Required parameters: user, pass

Test API: `curl "https://localhost/messages/?user=test&pass=1234"`

## Usage
### Install go (macos)
Install Brew (skip if you already did) 
```
ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```
Install go
`brew update&& brew install golang`
### Install go (linux) 
`sudo apt install golang-go`

### Build
`go build smsserve.go`

### Run
Create the user file and add a user
```
touch users
echo "test:1234" >> users
```
Run the server

`./smsserve`






