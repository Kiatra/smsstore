# smsstore
Minimalistic server that provides REST endpoints to receive and serve sms messages.

### Endpoints:
/messages/?user="username"&pass="password"
POST | GET

  Example Payload (POST):

  {"subject": "some subject", "message": "A sms Message"}

  Example Response (GET):

  Return last recived message as text.
