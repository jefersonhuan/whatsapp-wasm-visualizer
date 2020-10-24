const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
});

function loadChart(series) {
    document.getElementById("status").style.display = "none";
    Highcharts.chart('container', {
        lang: {
            months: ['Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho', 'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'],
            shortMonths: ['Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun', 'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez'],
            weekdays: ['Domingo', 'Segunda', 'Terça', 'Quarta', 'Quinta', 'Sexta', 'Sábado'],
            loading: ['Atualizando o gráfico...aguarde'],
            contextButtonTitle: 'Exportar gráfico',
            decimalPoint: ',',
            thousandsSep: '.',
            downloadJPEG: 'Baixar imagem JPEG',
            downloadPDF: 'Baixar arquivo PDF',
            downloadPNG: 'Baixar imagem PNG',
            downloadSVG: 'Baixar vetor SVG',
            printChart: 'Imprimir gráfico',
            rangeSelectorFrom: 'De',
            rangeSelectorTo: 'Para',
            rangeSelectorZoom: 'Zoom',
            resetZoom: 'Limpar Zoom',
            resetZoomTitle: 'Voltar Zoom para nível 1:1',
        },
        chart: {
            zoomType: 'x'
        },
        title: {
            text: 'Ocorrências de mensagens por tempo'
        },
        xAxis: {
            type: 'datetime'
        },
        yAxis: {
            title: {
                text: 'Número de mensagens'
            }
        },
        legend: {
            enabled: false
        },
        plotOptions: {
            column: {
                pointPadding: 0,
                borderWidth: 0,
                groupPadding: 0,
                shadow: false,
            },
            series: {
                borderColor: "#303030"
            }
        },

        series: [{
            type: 'column',
            name: 'Número de mensagens',
            data: series
        }]
    });
}

function handleFile() {
    document.getElementById("status").style.display = "block";
    const input = document.getElementById("file").files;
    const fileData = new Blob([input[0]]);

    const reader = new FileReader();
    reader.readAsArrayBuffer(fileData);
    reader.onload = function(){
        const bytes = new Uint8Array(reader.result);
        global.parseChat(bytes);
    }
}