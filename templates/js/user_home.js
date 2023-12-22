const searchBar = document.querySelector('.search-bar');

document.getElementById("getButton").addEventListener('click', async function () {
    try {
        let account = null;
        let token = null;

        const response = await fetch('/get-cookie', {
            method: "GET",
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const data = await response.json();
        account = data.account;
        token = data.token;

        console.log("Account:", account);
        console.log("Token:", token);

        const productId = parseInt(searchBar.value);
        searchBar.value = '';
        if (productId !== 0) {
            document.getElementById('bar').style.border = '2px solid green';
            const url = `http://localhost:8080/user/app/getProduct?productId=${productId}`;

            const secondResponse = await fetch(url, {
                method: "GET",
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `${token}`
                },
            });

            if (!secondResponse.ok) {
                throw new Error('Network response was not ok');
            }

            const responseData = await secondResponse.json();
            console.log(responseData);
            if (response.isRegistered == false) {
                document.getElementById("productId").textContent = "";
                document.getElementById("productName").textContent = "";
                document.getElementById("manufacturer").textContent = "";
                document.getElementById("isRegistered").textContent = "false";
                document.getElementById('response').innerHTML = '<p>Product not registered</p>';

            } else {
                document.getElementById("productId").textContent = responseData.productId;
                document.getElementById("productName").textContent = responseData.productName;
                document.getElementById("manufacturer").textContent = responseData.manufacturer;
                document.getElementById("isRegistered").textContent = responseData.isRegistered;
            }

        } else {
            document.getElementById('bar').style.border = '2px solid red';

            document.getElementById('response').innerHTML = '<p>Invalid id</p>';
        }

    } catch (error) {
        // Handle errors here
        console.error('There was a problem with the fetch operation:', error);
    }
});