# 🗳️ Plataforma de Votação Cooperativa

Sistema FullStack para gerenciamento de sessões de votação em assembleias de cooperativas. Desenvolvido como solução para desafio técnico, permite que associados votem de forma remota, segura e transparente.

---

## 📚 Tecnologias Utilizadas

### Frontend
- [React](https://react.dev/) + [Vite](https://vitejs.dev/)
- [TypeScript](https://www.typescriptlang.org/)
- [Redux Toolkit](https://redux-toolkit.js.org/)

### Backend
- [Go (Golang)](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin) (Framework web)
- [JWT](https://jwt.io/) (Autenticação)
- [PostgreSQL](https://www.postgresql.org/)

### DevOps
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Nginx](https://nginx.org/) (Proxy reverso)

---

## 🏗️ Arquitetura

A aplicação utiliza uma arquitetura containerizada com:
- **Nginx** como proxy reverso (porta 80)
- **Frontend** React + Vite (porta 5173)
- **Backend** Go + Gin (porta 8080)
- **PostgreSQL** como banco de dados

---

## 📚 API Endpoints

### Autenticação
- `POST /register` - Cadastrar usuário
- `POST /login` - Fazer login

### Pautas
- `POST /topics` - Criar pauta (protegido)
- `GET /topics` - Listar pautas

### Votação
- `POST /topics/{id}/session` - Abrir sessão (protegido)
- `POST /topics/{id}/vote` - Registrar voto (protegido)
- `GET /topics/{id}/result` - Ver resultados

> 📁 **Para testes detalhados**: Importe a collection `postman_collection.json` no Postman

---

## 🚀 Como Executar

1. **Clone o repositório**
```bash
git clone https://github.com/seu-usuario/seu-repositorio.git
cd desafio-tecnico-fullstack
```

2. **Execute a aplicação com Docker**
```bash
docker-compose up --build
```

3. **Acesse a aplicação**
- **Aplicação principal**: http://localhost (via Nginx)
- Frontend direto: http://localhost:5173
- Backend direto: http://localhost:8080

> 💡 **Nginx** está configurado como proxy reverso para servir a aplicação de forma integrada.

> ⚠️ **Observação**: O arquivo `.env` foi mantido no repositório para facilitar o teste da aplicação.

---

## 🚧 Dívidas Técnicas

### 1. Integração com API Externa de Verificação de CPF
**Problema**: A verificação da elegibilidade do associado para votar através da API `https://user-info.herokuapp.com/users/{cpf}` não foi implementada.

**Motivo**: O endpoint estava fora do ar ou inacessível durante o desenvolvimento.

**Solução Ideal**: 
- Implementar uma chamada HTTP para esse endpoint antes de registrar o voto
- Validar se o CPF é elegível para votar
- Implementar cache para melhor performance
- Tratar erros de timeout e indisponibilidade da API

### 2. Notificações em Tempo Real via MQTT
**Problema**: Sistema de notificações em tempo real não foi implementado.

**Motivo**: Priorização de funcionalidades core devido ao tempo limitado.

**Solução Ideal**:
- Configurar broker MQTT
- Publicar eventos quando sessões abrem/fecham
- Frontend subscribir para receber atualizações em tempo real

### 3. Interface do Usuário e Experiência
**Problema**: Interface funcional, com design ainda básico e focado na usabilidade mínima viável.

**Motivo**: Tempo limitado foi priorizado para implementar as funcionalidades core e aprender Redux.

**Solução Ideal**:
- Melhorar responsividade para diferentes dispositivos
- Adicionar um feedback visual mais rico
- Implementar validações de formulário mais elegantes
- Refatorar lógica de estado para ser mais robusta

📌 Observações Pessoais

 Este projeto foi meu primeiro desenvolvimento prático com Go. Já havia estudado a linguagem anteriormente, mas ainda não tinha tido a oportunidade de aplicá-la em um sistema completo. Foi um ótimo exercício para reforçar conceitos e estrutura de projeto em Go.

 Também foi minha primeira vez utilizando Redux com React. Estudei o básico durante o desafio para garantir a correta separação e gerenciamento de estado, mas reconheço que a implementação ainda pode evoluir.

Por conta do tempo do desafio, acredito que a base entregue atende bem aos requisitos, com algumas partes que ainda podem evoluir.