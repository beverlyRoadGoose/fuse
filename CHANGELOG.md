# Change Log

## v0.9.2

- Bump go version up to 1.20
- Bump github.com/sirupsen/logrus from 1.9.0 to 1.9.2
- Bump github.com/sirupsen/logrus from 1.9.2 to 1.9.3
- Bump github.com/stretchr/testify from 1.8.1 to 1.8.4

## v0.9.1

- Return error if send message request receives a non-200 status
- Return error if webhook registration/deletion requests receive a non-200 status

## v0.9.0

- Introduced Conversation Handler
- [breaking change] added error return type to handler functions

## v0.8.0

- Changed poll scheduling setup to allow for more frequent polls