# Aplicação Go Weather

Este projeto faz parte do laboratório da pós-graduação e fornece informações meteorológicas com base em um CEP (código postal brasileiro) fornecido.

## Pré-requisitos

- Docker e Docker Compose instalados
- Go 1.20+ instalado (para executar os testes localmente)

## Variáveis de Ambiente

Certifique-se de que as seguintes variáveis de ambiente estejam configuradas:

- `PORT`: A porta na qual a aplicação será executada (ex.: `8080`).
- `WEATHER_API_KEY`: Sua chave de [API para o serviço de clima](https://www.weatherapi.com/).

## Executando a Aplicação

1. Altere o arquivo `.env` com as variáveis de ambiente necessárias ( `PORT` e `WEATHER_API_KEY`).
  
```bash
cp .env.example .env
```

2. Construa e inicie a aplicação usando o seguintes comando:
  
```bash
make run
```

3. A aplicação estará disponível em `http://localhost:<PORT>`.

## Executando os Testes unitários

Para executar os testes unitários você pode usar o comando `make` para facilitar o processo. Siga os passos abaixo:

1. Execute o seguinte comando:
```bash
make test
```

Isso executará todos os testes do projeto.

## Construindo a Imagem Docker

Para construir a imagem Docker manualmente:

```bash
make build
```

## Endpoints

- **Health Check**: `GET /health`  
  Retorna o status de saúde da aplicação.

- **Obter Clima por CEP**: `GET /weather/{cep}`  
  Recupera informações meteorológicas para o CEP fornecido.

Para testar os endpoints, você pode usar ferramentas como `curl` ou Postman. Por exemplo:

```bash 

curl http://localhost:<PORT>/weather/20270150

curl https://goweather-109794580457.us-east1.run.app/weather/20270150

```


## Endereço da aplicação no Google Cloud Run

https://goweather-109794580457.us-east1.run.app/weather/25030170 