Setup (waarschijnlijk weet je dit al maar just in case):
1. Download en unzip
2. Open de folder in een terminal (cd smartcare of openen in VSCode)
3. go mod tidy
4. go run main.go
5. Open localhost:8080

Als Gin raar doet kan je hem ook handmatig installeren: go get github.com/gin-gonic/gin

Wat al werkt:
- Frontend van de pagina
- Knoppen doen al daadwerkelijk calls naar een API
- Structs en dataopslag

Wat nog niet werkt:
- De API zelf, hiervoor moeten de microservices (zoals een DB of message broker) nog ingesteld worden.