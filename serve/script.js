const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
});

function loadChart(series) {
    document.getElementById("status").style.display = "none";
    Highcharts.chart('container', {
        chart: {
            zoomType: 'x'
        },
        title: {
            text: 'Occurrence of WhatsApp messages over time'
        },
        xAxis: {
            type: 'datetime'
        },
        yAxis: {
            title: {
                text: '# of messages'
            }
        },
        legend: {
            enabled: false
        },
        plotOptions: {
            area: {
                fillColor: {
                    linearGradient: {
                        x1: 0,
                        y1: 0,
                        x2: 0,
                        y2: 1
                    },
                    stops: [
                        [0, Highcharts.getOptions().colors[0]],
                        [1, Highcharts.color(Highcharts.getOptions().colors[0]).setOpacity(0).get('rgba')]
                    ]
                },
                marker: {
                    radius: 2
                },
                lineWidth: 1,
                states: {
                    hover: {
                        lineWidth: 1
                    }
                },
                threshold: null
            }
        },

        series: [{
            type: 'area',
            name: '# of messages',
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