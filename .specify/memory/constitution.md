<!--
Sync Impact Report:
- Version change: 1.0.0 → 1.1.0
- List of modified principles:
  - Technology & Tools → Pivotado de TypeScript/Node para linguagem de sistemas (Go/Rust/Python-LowLevel) focada em engine de banco de dados imutável.
  - Performance Requirements → Redefinido para focar em inserção O(1), paginação de disco e concorrência ACID.
  - Simplicity & Scalability → Foco em isolamento via Docker e armazenamento baseado em Log Appender estruturado.
- Added sections:
  - Cryptographic Constraints & Integrity Gates
- Removed sections: None
- Templates requiring updates (✅ updated / ⚠ pending):
  - .specify/templates/plan-template.md (✅ updated)
  - .specify/templates/spec-template.md (✅ updated)
-->

# avocato-db Constitution

## Core Principles

### Code Quality & Determinism
O código do avocato-db deve ser estritamente idiomático, focado em gerenciar recursos de forma eficiente e previsível. Como lidamos com hashing criptográfico, **o determinismo é mandatório**: qualquer operação de serialização de dados (ex: JSON/MsgPack para Hash) deve garantir ordenação fixa de chaves. O código deve priorizar legibilidade e tratamento explícito de erros, eliminando falhas silenciosas que possam corromper o Ledger.

### Testing & Verification Standards
Testes automatizados são a nossa única linha de defesa contra corrupção de dados. É obrigatória a existência de testes de mutação e testes de integração que simulem:
1. Tentativas de injeção de dados corrompidos diretamente nos arquivos de dados.
2. Quedas repentinas de energia (crash-consistency) durante operações de escrita.
3. Validação matemática da integridade da cadeia de blocos de ponta a ponta.

### Performance & Memory Requirements
Performance é o core do avocato-db. 
* **Escrita:** A inserção de novos registros deve operar em tempo constante $O(1)$ através do append direto no fim do arquivo físico e atualização do ponteiro do último hash válido em cache de memória volátil.
* **Leitura:** Buscas por bloco por meio de hash ou índice devem utilizar índices estruturados eficientes para manter a latência abaixo de 10ms.
* **Alocação:** Minimizar alocações dinâmicas de memória no caminho crítico de escrita.

### Simplicity & Scalability via Containerization
Favorecemos a simplicidade de infraestrutura. O avocato-db deve nascer e operar de forma 100% isolada e autocontida em ambientes Docker. Ele deve ser agnóstico ao sistema operacional do host. A escalabilidade do storage deve ser garantida via volumes Docker otimizados para alta taxa de I/O de escrita sequencial.

## Technology & Tools
A stack tecnológica é selecionada estritamente para garantir concorrência segura, performance de baixo nível e manipulação direta de I/O de arquivos:

* **Core Engine:** **Go (Golang)** — Selecionadas pela ausência de runtime pesado, excelente suporte a primitivas de concorrência (goroutines/async), manipulação nativa de ponteiros/arquivos e bibliotecas criptográficas robustas na standard library.
* **Criptografia:** SHA-256 (ou Blake3) para geração de hashes estruturais de blocos.
* **API Layer:** Protocolos leves e de alta performance como gRPC ou HTTP/REST puro (sem frameworks pesados).
* **Environment:** **Docker & Docker Compose** para padronização do ambiente produtivo, garantindo que o banco de dados inicialize com privilégios restritos de escrita direta ao disco da engine.

## Cryptographic Constraints & Integrity Gates
O motor do avocato-db deve seguir as seguintes regras invioláveis a nível de arquitetura de software:
1. **Append-Only no Engine Level:** O código do driver de storage não deve possuir ou expor rotinas para `UPDATE`, `DELETE` ou `TRUNCATE` de blocos persistidos.
2. **Chain Coupling:** Cada bloco persistido deve conter obrigatoriamente o hash do bloco imediatamente anterior.
3. **Boot Integrity Check:** A cada inicialização do container Docker, o avocato-db executará uma varredura completa de integridade física e criptográfica na base de dados antes de abrir a porta da API para conexões.

## Quality Gates
Cada Pull Request deve passar por validações automatizadas em pipelines de CI/CD:
* Verificação estrita de linters e detectores de race conditions.
* Passar por 100% dos testes de quebra de integridade criptográfica.
* Revisão de código obrigatória por ao menos um mantenedor focando na garantia de imutabilidade do fluxo de dados implementado.

## Governance
Esta constituição dita as regras de engenharia do avocato-db. Qualquer alteração de design que viole os princípios de imutabilidade ou mude a stack de sistemas aqui definida exige justificativa técnica profunda em formato RFC e aprovação unânime dos mantenedores core do projeto.

**Version**: 1.1.0 | **Ratified**: 2026-05-20 | **Last Amended**: 2026-05-20