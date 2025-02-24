document.addEventListener('DOMContentLoaded', () => {
    if (window.location.href.includes("/orders")) {
        startDate();
    }
    else if (window.location.href.includes("/couriers")) {
        ShowCourier();
    }
});

let orders = [];


function startDate(){
    let options = {
        method: 'GET',
        headers: {
            "Content-Type": "application/json"
        },
    }; 
    fetch('http://localhost/orders', options)
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(allOrders => {
        orders = allOrders;
        showOrders(allOrders);
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
}
function showOrders(jsonOrders)
{
    table = document.getElementById('customer-table');
    
    let tableContent = ''; // Создаем переменную для хранения содержимого таблицы

    for(let object of jsonOrders) {
        tableContent += `<tr> 
            <td>${object.date_order}</td> 
            <td>${object.date_pay}</td> 
            <td>${object.type_pay}</td> 
            <td>${object.client_name}</td> 
            <td>${object.name_product}</td> 
            <td>${object.count_product}</td> 
            <td>${object.price_product}</td> 
            <td>${object.count_warehouse}</td> 
            <td>${object.courier_name}</td> 
            <td>${object.price_delivery}</td> 
        </tr>`;
    }

    table.innerHTML = tableContent; // Вставляем все строки в таблицу
}

document.getElementById('addCourierForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Предотвращаем перезагрузку страницы

    const firstName = document.getElementById('firstName').value;
    const lastName = document.getElementById('lastName').value;

    fetch('http://localhost/couriers', {
        method: 'POST',
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ first_name: firstName, last_name: lastName })
    })
    .then(response => response.json())
    .then(data => {
        console.log('Courier added:', data);
        // Обновляем таблицу
        const tableBody = document.getElementById('customer-table');
        tableBody.innerHTML += `<tr><td>${firstName} ${lastName}</td></tr>`;
    })
    .catch(error => console.error('Error adding courier:', error));
});

function ShowCourier(){
    fetch('http://localhost/couriers', {
        method: 'GET',
        headers: {
            "Content-Type": "application/json"
        }
    })
    .then(response => response.json())
    .then(couriers => {
        const tableBody = document.getElementById('customer-table');
        couriers.forEach(courier => {
            tableBody.innerHTML += `<tr><td>${courier.first_name} ${courier.last_name}</td></tr>`;
        });
    })
    .catch(error => console.error('Error fetching couriers:', error));
}