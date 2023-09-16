# Opslink [Opsgenie - Gsuite/Slack]

![bpmn drawio](https://github.com/thallesdaniell/bot-opsgenie/assets/20145406/d3658169-1542-4a65-bec6-1bf40eb7aa11)

## Objetivo
- Atualizar usuários que estão no on-call(opsgenie) automaticamente no grupo xxx-oncall@corp.com(google) com permissão de edição e no canal oncall-incidentes.

## Processo:
- Recupera os membros do oncall tanto primário como secundário no ```Opsgenie```.
- Atualiza os usuários caso precise no ```Google``` .
- Atualiza os usuários caso precise no canal  ```@oncall-incidentes```  do ```Slack``` .

## Implantação:
- Implementado duas cronjobs, uma que executa a cada 10minutos e outra que executa toda segunda as 09:00.

## Execução local:
- Para executar o projeto crie um arquivo `.env` e insira as evns conforme abaixo e tenha acesso a service account.
```go
APP_ENV="production"
OPSGENIE_API_KEY="xxxxx-xxxxxx-xxxxxx"
OPSGENIE_ON_CALL_SCHEDULE="xxxxx-xxxxxx-xxxxxx,xxxxx-xxxxxx-xxxxxx"
GOOGLE_APPLICATION_CREDENTIALS="./service-account.json"
GOOGLE_SUBJECT_EMAIL="demo-demo@corp.com"
GOOGLE_GROUP_KEY="oncall@corp.com"
ADDITIONAL_USERS_SLACK_GROUP="XXXXXXX"
OPSGENIE_ON_CALL_SCHEDULE_PRIMARY="xxxxx-xxxxxx-xxxxxx,xxxxx-xxxxxx-xxxxxx"
SLACK_GROUP_ID_UPDATE_ONCALL="XXXXXXX"
SLACK_GROUP_ID_NEXT_ONCALL="XXXXXXX"
SLACK_BOT_TOKEN="xxxxx-xxxxxx-xxxxxx,xxxxx-xxxxxx-xxxxxx"
```

- Execute o docer compose
```docker
docker-compose run opslink
```

- Dentro do container
```go
go run cmd/main.go
```

## Execução produção:
- Crie uma aplicação no spinnaker
- Criar uma secrets com os mesmos valores do `.env`
- Criar uma cronjob apartir da pasta manifestos/cronjob.yaml
