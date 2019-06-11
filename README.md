# dito-challenge

Este repositório contém uma aplicação desenvolvida para um teste técnico, onde as funcionalidades exigidas eram uma api coletora, um autocomplete e uma timeline.

# Tecnologias utilizadas
## Armazenamento: 
* Elasticsearch na versão 7 - https://www.elastic.co/pt/products/elasticsearch
## API: 
* Golang - https://golang.org/
* Framework Gin - https://gin-gonic.com/
## Provisionamento do Elasticsearch: 
* Docker - https://www.docker.com/
* Docker Compose - https://docs.docker.com/compose/

# Demo

* **API Coletora**: poderá ser acessada através da url http://localhost:5000/v1/events **METHOD:POST**  

Objeto do request:

```json
{
"event": "buy",
"timestamp": "2016-09-22T13:57:31.2311892-04:00"
}
```
* **Autocomplete**: poderá ser acessada através da url http://localhost:5000/v1/autocomplete?q=termopesquisado **METHOD:GET**

* **Timeline**: poderá ser acessada através da url http://localhost:5000/v1/timeline **METHOD:GET**

## Requisitos: 
* Sitema operacional baseado em Unix
* Git - https://git-scm.com/downloads
* Docker - https://www.docker.com/products/docker-desktop
* Docker Compose - https://docs.docker.com/compose/install/

## Passos a passo: 
1 - Clonar o repositório
```bash
git clone https://github.com/jfredericon/dito-challenge && cd dito-challenge
```
2 - Subir o container com o Elasticsearch: 
```bash
cd elasticsearch && make
```
3 - Subir a aplicação
```bash
cd app && make
```

É importante que o container do Elasticseach esteja de pé para que a aplicação rode de forma correta. 
