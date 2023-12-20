const searchBar = document.querySelector('.search-bar');

document.getElementById("setButton").addEventListener('click', async function() {
    document.getElementById('response').innerHTML = '';

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

       // const newValue = Math.floor(Math.random() * (100 - 1)) + 1;
        const newValue = parseInt(searchBar.value);
        searchBar.value = '';

        if (Number.isInteger(newValue) && newValue > 0 && Number.MAX_SAFE_INTEGER >= newValue) {
            document.getElementById('bar').style.border = '2px solid green';

            const url = `http://localhost:8080/admin/app/setContractValue?account=${account}&newvalue=${newValue}`;

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
            console.log("updated value " + responseData.value)
            document.getElementById("blockValue").textContent = responseData.value;

            document.getElementById('response').innerHTML = '<p>Value updated on contract!</p>';

        } else {
            document.getElementById('bar').style.border = '2px solid red';

            document.getElementById('response').innerHTML = '<p>Value must be an integer greater than 0!</p>';
        }

    } catch (error) {
        // Handle errors here
        console.error('There was a problem with the fetch operation:', error);
    }
});

document.getElementById("getButton").addEventListener('click', async function() {
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

        const url = `http://localhost:8080/admin/app/getContractValue`;

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
        console.log("updated value " + responseData.value)
        document.getElementById("blockValue").textContent = responseData.value;
    } catch (error) {
        // Handle errors here
        console.error('There was a problem with the fetch operation:', error);
    }
});