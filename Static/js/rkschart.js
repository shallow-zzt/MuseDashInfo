
var ctx = document.getElementById('user-rks-chart').getContext('2d');
var myChart = new Chart(ctx, {
    type: 'line',
    data: {
        labels: [], 
        datasets: [{
            label: '玩家的RKS',
            data: [],
            borderColor: 'rgba(75, 192, 192, 1)',
            backgroundColor: 'rgba(75, 192, 192, 0.2)', 
            borderWidth: 2, 
            pointRadius: 5, 
            fill: true
        }]
    },
    options: {
        responsive: false,
        plugins: {
            legend: {
                display: false,
                position: 'top'
            }
        },
        scales: {
            x: {
                title: {
                    display: false,
                    text: '时间'
                }
            },
            y: {
                title: {
                    display: true,
                    text: 'RKS'
                },
                beginAtZero: false
            }
        }
    }
});