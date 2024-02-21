## Testes do projeto
- Realize as configurações desejadas de limitação no arquivo .env
- Execute para subir o redis: docker-compose up -d
- Execute a aplicação: go run main.go persistence.go rate_limiter_middleware.go rate_limit.go
- Execute a requisição sem header pelo arquivo ./api/get.http
- Execute a requisição com header pelo arquivo ./api/get-header.http
- Execute para rodar a limitação por IP: curl --location 'http://localhost:8080'
- Execute para rodar a limitação por API_KEY: curl --location 'http://localhost:8080' --header 'API_KEY:API_KEY_A'
- Execute a automação de testes em load_test.go

## Documentação main.go
- Visão Geral
Esta aplicação Go é um servidor web simples que utiliza o framework Gin para gerenciar rotas HTTP. O objetivo principal é demonstrar um middleware de controle de taxa (rate limiter) que limita o número de requisições por IP ou por token de acesso.

- Configuração: Variáveis de Ambiente
A aplicação utiliza um arquivo .env para configurar variáveis de ambiente. As seguintes variáveis são suportadas:

REDIS_ADDR: Endereço do servidor Redis.
REDIS_PASSWORD: Senha do servidor Redis (opcional).
BLOCKING_TIME_OUT_IN_SECONDS_PER_IP: Tempo de bloqueio em segundos para limitação por IP.
BLOCKING_TIME_OUT_IN_SECONDS_PER_API_KEY: Tempo de bloqueio em segundos para limitação por API_KEY no Header da requisição.
LIMIT_IP_: Limite de requisições por IP.
LIMIT_API_KEY_: Limite de requisições por API_KEY no Header da requisição.
Certifique-se de fornecer essas variáveis no arquivo .env antes de iniciar a aplicação.

- Estrutura do Código
O arquivo main.go é o ponto de entrada da aplicação e contém a lógica principal.

Função init(): Carrega as variáveis de ambiente do arquivo .env e inicializa o banco de dados escolhido chamando a função initialize().

Função main(): Inicializa o servidor web utilizando o framework Gin. Adiciona um middleware de controle de taxa (rateLimiterMiddleware) e define uma rota padrão que responde a requisições HTTP na raiz ("/"). O servidor é executado na porta 8080.

O arquivo middleware.go contém o middleware rateLimiterMiddleware para controle de taxa.

Função rateLimiterMiddleware(): Retorna um middleware Gin que verifica se o número de requisições excede o limite definido por IP ou token de acesso. Responde com um código 429 (Limite de Requisições Excedido) caso o limite seja atingido, com a mensagem "you have reached the maximum number of requests or actions allowed within a certain time frame"

- Executando a Aplicação
Certifique-se de ter o Go instalado no seu ambiente de desenvolvimento. Execute o seguinte comando para iniciar a aplicação: go run main.go

- Considerações Finais
Esta documentação fornece uma visão geral da estrutura e funcionamento do arquivo main.go. Certifique-se de ajustar as variáveis de ambiente conforme necessário para o seu ambiente específico.


## Documentação rate_limiter_middleware.go
Este middleware tem como objetivo limitar o número de requisições por IP ou token de acesso (API_KEY) em um servidor web construído com o framework Gin. Ele verifica se o número de requisições recebidas de um único endereço IP ou com um token de acesso específico excedeu o limite configurado.

- Funcionamento
Checagem do Limite de Requisições por IP ou Token de Acesso:

A função checkRateLimit é chamada para verificar se o número de requisições excedeu o limite configurado para o IP ou token de acesso (API_KEY).
Se o limite for excedido, o middleware responde com um código de status HTTP 429 (Limite de Requisições Excedido) e uma mensagem explicativa.
Caso contrário, o middleware permite a execução normal da requisição.

- Resposta em Caso de Limite Excedido:
Se o limite for excedido, o middleware envia uma resposta com o código de status HTTP 429 e a mensagem "you have reached the maximum number of requests or actions allowed within a certain time frame", indicando que o número máximo de requisições foi atingido dentro de um determinado intervalo de tempo.
A execução do middleware é interrompida usando context.Abort().

- Resposta em Caso de Limite Não Excedido:
Se o número de requisições por IP ou token de acesso não exceder o limite, o middleware responde com um código de status HTTP 200 (OK) e uma mensagem indicando que a requisição foi permitida. A mensagem inclui também o número total de requisições até o momento.

- Uso
O middleware pode ser utilizado adicionando a seguinte linha ao roteador Gin:

router := gin.Default()
router.Use(rateLimiterMiddleware())

Dessa forma, todas as rotas registradas no roteador serão submetidas ao controle de taxa implementado por este middleware.

O middleware é projetado para funcionar em conjunto com o framework Gin, aproveitando a estrutura de middleware do Gin para interceptar as requisições HTTP.
Este middleware oferece uma solução eficaz para limitar o tráfego de requisições em servidores web, proporcionando uma experiência mais robusta e segura.


# Documentação rate_limit.go
A função checkRateLimit é responsável por verificar se o número de requisições por IP ou token de acesso (API_KEY) excedeu o limite configurado. Ela é parte integrante do middleware de controle de taxa e desempenha um papel fundamental na decisão de permitir ou bloquear uma requisição com base no controle de taxa estabelecido.

- Funcionamento
Identificação do Tipo de Requisição (IP ou API_KEY):

A função verifica a presença do header "API_KEY" na requisição. Se presente, a requisição é tratada como uma requisição por token de acesso (API_KEY). Caso contrário, é considerada uma requisição por IP.

- Configuração do Bloqueio e Limite
Se a requisição é por token de acesso, a função busca o tempo de bloqueio (expiração) e o limite de requisições por API_KEY no arquivo .env.
Se a requisição é por IP, a função obtém o IP do cliente da requisição, substituindo "::1" por "127.0.0.1" caso seja localhost. Em seguida, ela busca o tempo de bloqueio e o limite de requisições por IP no arquivo .env.

- Consulta do Contador de Requisições
Utilizando a função getRequestCount, a função obtém o número atual de requisições feitas por IP ou API_KEY.

- Inicialização e Expiração do Contador
Se o número de requisições for 1 (indicando uma nova sequência de requisições), a função inicializa o contador de requisições e define o tempo de expiração com base nas configurações obtidas.
O tempo de expiração é definido pela função setTimeToExpireKey.

- Verificação do Limite
A função verifica se o número de requisições excede o limite configurado.
Se exceder, a função envia uma resposta de erro indicando que o número máximo de requisições foi atingido.

- Retorno de Resultados:
A função retorna um par de valores booleanos. O primeiro valor indica se o limite foi excedido, e o segundo valor é o número total de requisições.

- Uso
Esta função é utilizada internamente pelo middleware de controle de taxa. Não é recomendado chamar diretamente em outras partes do código.

- Observações
As configurações de limite, tempo de bloqueio e outras opções são fornecidas através de variáveis de ambiente ou de um arquivo .env.
O controle de taxa é baseado no uso de um banco de dados para armazenar e consultar informações sobre as requisições.
Esta função desempenha um papel crucial na implementação do controle de taxa, garantindo que as requisições sejam tratadas de acordo com as configurações especificadas. 


## Documentação persistence.go
Tem um conjunto de funções responsável por inicializar a conexão com o servidor Redis e gerenciar a expiração e contagem de requisições no contexto do controle de taxa.

- Funcionalidades
A função initialize é responsável por inicializar a conexão com o servidor Redis. Ela utiliza as informações de endereço e senha do Redis provenientes do arquivo de configuração .env.

- Uso
Esta função é chamada no momento da inicialização da aplicação para garantir que a conexão com o Redis esteja pronta para ser utilizada.

- setTimeToExpireKey
A função setTimeToExpireKey é responsável por definir o tempo de expiração para um determinado contador no Redis. Ela é utilizada para configurar a expiração do contador de requisições por IP ou API_KEY.

- getRequestCount
A função getRequestCount é responsável por incrementar o contador de requisições no Redis e retornar o número atual de requisições feitas por IP ou API_KEY.

- Observações
Todas as operações no Redis são realizadas utilizando a biblioteca github.com/go-redis/redis/v8.
As configurações do servidor Redis (como endereço e senha) são obtidas a partir do arquivo .env. 
Estas funções são parte essencial do controle de taxa, garantindo a contagem correta e a expiração adequada dos contadores no Redis. Certifique-se de que a conexão com o Redis esteja funcionando corretamente para um controle eficaz de taxa.