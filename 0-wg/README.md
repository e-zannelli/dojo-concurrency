La fonction process reçoit un Client et un channel recevant des strings.

Le client à l'interface
```go
type Client interface {
    Call(string) Response
}
type Response struct {
    Body string
}
```

La fonction process doit appeler la fonction Call du client pour chaque élément de generator, le retour n'est pas vérifié.

Modifier le code pour faire les appels de manière concurrente.