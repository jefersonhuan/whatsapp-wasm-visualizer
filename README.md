[Demonstração](https://wp-viz-wasm-4sdd6inywq-rj.a.run.app/)

Esta é uma simples aplicação-experimento utilizando WebAssembly com Go, que tem por 
intuito gerar um gráfico em série (timeseries) de mensagens diárias 
de uma conversa WhatsApp.


### Como funciona

O processo se inicia com a seleção do arquivo .txt exportado do WhatsApp; esse evento
ativa uma função JS que converte o conteúdo em um vetor de bytes. Por que isso? 
Bom, apesar das possibilidades do WASM, acessar arquivos do sistema não é uma delas; sendo assim,
essa é uma alternativa barata e pouca custosa neste contexto.

Após a conversão o código JS então invoca uma função de Go denominada
`parseChat`; essa função recebe o conteúdo e o coverte para o contexto de Go (em outras 
palavras, um `[]byte`).

Esse conteúdo é então decodificado, extraindo-se os dias e a ocorrência de mensagens,
para posteriormente ser convertido para o contexto de JS (um vetor 2D, 
que corresponde a um `[]interface{}`) e chamando a renderização do gráfico.

### Considerações

Este foi apenas um experimento com WASM, mas bem interessante diga-se de passagem. Com 
pouca complexidade foi possível ler uma conversa com mais de 40k linhas em menos de um segundo,
o que foi... bem satisfatório :)

A única dificuldade (ou obstáculo) foi conciliar os tipos de JS com Go sem perder
a performance do último. Com isso, o código tem uma ordem de processamento
sequencial* de `O(n + y)`, sendo
`n` o número de linhas no arquivo, que corresponde à decodificação inicial (devidamente
tipada) e `y` o número de dias, que corresponde à conversão aos tipos de JavaScript.

\* A exportação tem como base a ordem de chegada das mensagens no celular, dessa forma, acontece 
uma situação como a tal:
- Usuário A envia mensagem em 01/01/2020, às 23h50, porém perde a conexão com a internet;
- Usuário B envia uma mensagem em 02/01/2020, às 01h03;
- Usuário A volta a ter conexão, às 01h10, chegando finalmente a primeira mensagem
  ao celular de B;
- Ao exportar a primeira mensagem aparece **depois** da segunda.

Para contornar isso (e é algo que pode acontecer com frequência num grupo),
o programa tenta fazer uma busca invertida por até 100ms (o que é suficiente para a maioria
dos casos), corrigindo a contagem final.

Caso isso aconteça, é possível conferir a correção abrindo o console do navegador.