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
pouca complexidade foi possível ler uma conversa com mais de 40k linhas em menos de meio
segundo, o que foi... satisfatório :)

A única dificuldade (ou obstáculo) foi conciliar os tipos de JS com Go sem perder
a performance do último. Com isso, o código tem uma ordem de `O(n + y)`, sendo
`n` o número de linhas no arquivo, que corresponde à decodificação inicial (devidamente
tipada) e `y` o número de dias, que corresponde à conversão aos tipos de JavaScript.
