import Chart from 'chart.js';

let fb: any = document.
    getElementById('testButton').
    addEventListener('click', sayHello);

function sayHello(): void {
    console.log('hello, there');
}

let chartDiv: any = document.getElementById("ordersChart")
let myLineChart = new Chart(ctx, {
    type: 'line',
    data: data,
    options: options
});