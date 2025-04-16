document.addEventListener('DOMContentLoaded', () => {
    if (window.location.href.includes("/orders")) {
        
        if (localStorage.getItem('jwtToken') !== null)
            startDate();
        else {
            window.location.href = 'login.html';
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
    fetch('http://localhost:8080/products', {
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
            <div class="product-price">${product.Product_price}$ count: ${product.Product_count}</div> </div>`;
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
    fetch('http://localhost:8080/orders', options)
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

function addCorier(){
    const firstName = document.getElementById('firstName').value;
    const lastName = document.getElementById('lastName').value;
    const id_warehouse = document.getElementById('id_warehouse').value;
    
    fetch('http://localhost:8080/courierAdd', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('jwtToken')}`,
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ first_name: firstName, last_name: lastName, id_warehouse: Number(id_warehouse)})
    })
    .then(response => response.json())
    .then(data => {
        console.log('Courier added:', data);
        // Обновляем таблицу
        /*const tableBody = document.getElementById('customer-table');
        tableBody.innerHTML += `<tr><td>${firstName} ${lastName}</td> <td> ${0} </td> <td>${String(id_warehouse)}</td></tr>`;*/
    //    alert("Courier added successfully")
        ShowCouriers()
    })
    .catch(error => {
    //    console.error('Error adding courier:', error)
        return;
    });

};

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

function addProduct(){
    const productName = document.getElementById('productName').value;
    const productPrice = document.getElementById('productPrice').value;
    const productDescription = document.getElementById('productDescription').value;
    const productCount = document.getElementById('productCount').value;
    const productWarehouse = document.getElementById('productWarehouse').value;

    fetch('http://localhost:8080/productAdd', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('jwtToken')}`,
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            product_name: productName,
            product_price: Number(productPrice),
            product_description: productDescription,
            product_count: Number(productCount),
            id_warehouse: Number(productWarehouse)
        })
    })
    .then(response => response.json())
    .then(data => {
        //console.log('Product added:', data);
        //alert("Product added successfully");
        ShowProducts(); // Функция обновления таблицы
    })
    .catch(error => {
        //console.error('Error adding product:', error);
        //alert("Ошибка добавления продукта");
    });
};