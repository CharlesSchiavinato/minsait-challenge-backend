## Descrição da API
API de Fluxo de Caixa
Está API tem como objetivo registrar lançamentos de crédito/débito e demonstrar o saldo diário.
É possível realizar todas as operações de CRUD de lançamentos e consultar o saldo passando uma data de referencia.
Decidi retornar o saldo com valor zero quando não existir nenhum lançamento cadastrado para a data informada mas poderia devolver um erro informando que não foi encontrado nenhum lançamento para o data informada.
A documentação dos endpoints estão disponíveis na própria API [localhost:9000/api/docs](localhost:9000/api/docs).

## Etapas para poder executar a API na máquina local
1. Docker instalado. [Documentação](https://docs.docker.com/engine/install/)
2. Docker-compose instalado. [Documentação](https://docs.docker.com/compose/install/linux/)
3. Go instalado. [Documentação](https://go.dev/doc/install)

#### Obs: Os comandos a seguir devem ser executados na pasta raiz do projeto.

4. Subir os serviços de banco de dados e de cache utilizados pela API
    ```
    make docker-compose-up
   
    ou
    
    docker-compose up -d
    ```

5. Executar a API

    - Direto na máquina local
        ```
        make go-run

        ou
        
        go run server.go
        ```

    - Dentro do Docker
        ```
        make docker-build
        make docker-run
        ```

6. Endpoint da API [localhost:9000/api](localhost:9000/api)


## Descrição Técnica
1. Banco de dados Postgres por ser um serviço de banco relacional robusto, completo e open source que atende perfeitamente desde pequenas aplicações até aplicações robustas e compatível com serviços de banco
2. Cache Redis por ser um serviço de cache robusto, open source, amplamente utilizado e compatível com serviço de cache em nuvem como o memory store do GCP. Apesar de não ter utilizado o cache nos endpoints atuais já deixei implementado toda estrutura necessária demonstrando que tenho o hábito de utilizar sempre que necessário.
3. Migration para versionamento de alterações no banco de dados.
4. Health Check [localhost:9000/api/healthz](localhost:9000/api/healthz) para monitorar se a aplicação está no ar e se os serviços de banco de dados e cache estão funcionando.
5. Swagger para geração automática da documentação a partir de tags adicionadas ao código e também para exibição dos exemplos de requisições para os endpoints da API.
6. CORS para poder permitir requisições de origem diferente da API.
7. Middleware para logar o resultado de todas as requisições contendo informações da origem da requisição e request id para ajudar no troubleshooting da aplicação. O request id ajuda a rastrear uma mesma requisição por diversos micro serviços e as informações da origem ajudam a identificar se o problema está relacionado a uma origem ou dispositivo específico.
8. Parametrização das configurações da API por meio de variáveis de ambiente ou arquivo de configuração "config.env" na pasta raiz da aplicação.
9. Testes Unitários e Integração. Nesse projeto o foco foi no padrão teste de integração de componentes mas fiz alguns testes unitários utilizando mock e alguns de integração completo demostrando que tenho total conhecimento para trabalhar com o que a empresa preferir. Implementei todas as funcionalidades da API utilizando a técnica de TDD começando pela implementação de testes e funcionalidades das regras de negócios e depois expandindo para testes e funcionalidades das demais camadas da aplicação.
- Teste unitário utilizando Mock;
- Teste de integração de componentes utilizando o conceito database in memory;
- Teste de integração completo;

10. Desenhos contendo o fluxo de funcionalidade dos enpoints da solução utlizando o Mermaid Markdown e disponibilizados na pasta /docs.
11. O projeto já contém um arquivo config.env com todas as configurações necessárias para poder executar a API no ambiente local.
12. Docker-compose para poder subir os serviços de banco de dados e cache para poder rodas a API no ambiente local.
13. Dockerfile para poder realizar o build da API e gerar imagem docker para rodar no ambiente local.
14. Github Action para validar PR e Push para a branch main iniciando um processo de CI/CD contemplando a execução dos testes unitários, testes de integração e build do projeto.
15. Makefile para poder executar de forma simples diversos comandos.
16. É possível trocar o manipulador de rotas, banco de dados e serviço de cache facilmente devido à utilização do Clean Architecture no projeto.

## Observação
1. Faltou implementar um filtro de período de data na listagem de lançamentos
2. Apesar de não ter aplicado nesse projeto também tenho conhecimento do padrão conventional commits.


## Geração da Documentação da API - Swagger

Executar o comando abaixo na pasta raiz da API para atualizar a documentação
    ```
    make swag
    ```

## Sugestões de Melhorias que Normalmente Implemento nas APIs

1. Incluir cabeçalhos HTTP de segurança
2. Incluir controle de auditoria
3. Incluir na documentação da API a relação dos erros que podem ser retornado
4. Paginação no endpoint de listagem.
5. Substituição do ID sequencial por UUID.


    #### **Obs:** Com certeza tem mais melhorias a ser feita tanto no código quanto na documentação. Melhoria contínua deve fazer parte da vida útil de toda aplicação.

# Espero que gostem bastante do projeto que entreguei ;)
