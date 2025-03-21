document.addEventListener('DOMContentLoaded', () => {
    if (window.location.href.includes("/orders")) {
        
        if (localStorage.getItem('jwtToken') !== null)
            startDate();
        else {
            window.location.href = 'login.html'
        }

    }
    else if (window.location.href.includes("/couriers")) {
        ShowCouriers();
    }
    else if (window.location.href.includes("/products")) {
        showProducts();
    }
});

let orders = [];

function showProducts(){
    fetch('http://localhost/products', {
        method: 'GET',
        headers: {
            "Content-Type": "application/json"
        }
    })
    .then(response => response.json())
    .then(products => {
        const tableBody = document.getElementById('products-container');
        products.forEach(product => {
        tableBody.innerHTML += `<div class="product-card">
            <div class="product-name">${product.Product_name}</div>
            <div class="product-description">${product.Product_description}</div>
            <div class="product-price">${product.Product_price}</div> </div>`;
        });
    })
    .catch(error => console.error('Error fetching couriers:', error));

}

function startDate(){
    let options = {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('jwtToken')}`,
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
    const id_warehouse = document.getElementById('id_warehouse').value;
    
    fetch('http://localhost/courierAdd', {
        method: 'POST',
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ id_courier: 0, first_name: firstName, last_name: lastName, total_salary: 0, warehouse_name: "warehouse", id_warehouse: Number(id_warehouse)})
    })
    .then(response => response.json())
    .then(data => {
        console.log('Courier added:', data);
        // Обновляем таблицу
        /*const tableBody = document.getElementById('customer-table');
        tableBody.innerHTML += `<tr><td>${firstName} ${lastName}</td> <td> ${0} </td> <td>${String(id_warehouse)}</td></tr>`;*/
        ShowCouriers()
    })
    .catch(error => {
        console.error('Error adding courier:', error)
        return;
    });
    alert("Courier added successfully")
});

function ShowCouriers(){
    fetch('http://localhost/couriers', {
        method: 'GET',
        headers: {
            "Content-Type": "application/json"
        }
    })
    .then(response => response.json())
    .then(couriers => {
        const tableBody = document.getElementById('customer-table');
        tableBody.innerHTML = "";
        couriers.forEach(courier => {
            tableBody.innerHTML += `<tr><td>${courier.first_name} ${courier.last_name}</td> <td>${courier.warehouse_name}</td> <td>${courier.total_salary}</td> </tr>`;
        });
    })
    .catch(error => console.error('Error fetching couriers:', error));
}
