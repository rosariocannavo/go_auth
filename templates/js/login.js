var metamaskAddress = null;

document.getElementById("loginForm").addEventListener("submit", async function (event) {
    event.preventDefault();

    const formData = new FormData(this); // 'this' refers to the form element
    metamaskAddress = await getMetaMaskAddress();

    if (metamaskAddress !== null) {
        formData.append("metamaskAddress", metamaskAddress);

        var object = {};
        formData.forEach(function (value, key) {
            object[key] = value;
        });
        var json = JSON.stringify(object);

        console.log(json)

        // Send POST request to Go Gin server
        fetch("/login", {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
            },
            body: json
        })
            .then(response => {
                if (response.status === 401) {
                    console.log("Unauthorized");
                    document.getElementById('password').style.border = '2px solid red';
                    document.getElementById('password').value = '';
                    document.getElementById('response').innerHTML = '<p>Passwords Incorrect. Please retry.</p>';

                    throw new Error("Unauthorized");

                } else if (response.status === 400) {
                    console.log("user not present");
                    document.getElementById('password').style.border = '2px solid red';
                    document.getElementById('username').style.border = '2px solid red';
                    document.getElementById('password').value = '';
                    document.getElementById('username').value = '';
                    document.getElementById('response').innerHTML = '<p>User not found. Please retry.</p>';

                    throw new Error("User not present");

                } else {
                    document.getElementById('password').style.border = '2px solid green';
                    document.getElementById('response').innerHTML = '<p>Passwords Checked. Metamask redirect.</p>';
                }

                return response.json();
            })
            .then(data => {
                const nonce = data.Nonce;
                console.log("Nonce: " + nonce);
                requestMetaMaskSignature(nonce);
            })
            .catch(error => {
                if (error.message === "Unauthorized") {
                    console.log("Unauthorized request");
                    // Handle unauthorized error here (e.g., show a message to the user)
                } else {
                    console.error("Generic error occurred:", error);
                    // Handle other generic errors (e.g., display a generic error message)
                    // Inform the user or perform necessary actions for unexpected errors
                }
            });


    } else {
        console.error("Metamask address not available");
        // Handle the case when Metamask address is not available
    }
});

async function getMetaMaskAddress() {
    if (typeof window.ethereum !== 'undefined') {
        // Metamask is available
        const provider = window.ethereum;

        try {
            // Request access to accounts
            const accounts = await provider.request({ method: 'eth_requestAccounts' });
            const accountAddress = accounts[0]; // Get the first account
            console.log('Account Address:', accountAddress);

            return accountAddress
        } catch (error) {
            console.error('Error:', error);
        }
    } else {
        // Metamask is not available
        console.error('Metamask extension not detected');
        return null
    }
}

async function requestMetaMaskSignature(nonce) {
    // Metamask is available
    //const nonce = "{{.Nonce}}";
    console.log(nonce)

    const provider = window.ethereum;

    try {
        // Request access to accounts
        const accounts = await provider.request({ method: 'eth_requestAccounts' });
        const accountAddress = accounts[0]; // Get the first account

        sessionStorage.setItem('accountAddress', accountAddress);

        const encodedMessage = stringToHex(nonce);

        const signature = await provider.request({
            method: 'personal_sign',
            params: [encodedMessage, accountAddress],
        });

        console.log(signature)

        // Send the signed message and Ethereum address to the backend
        const requestOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ message: nonce, address: accountAddress, signature: signature }),
        };

        const response = await fetch('/verify-signature', requestOptions)
        if (response.ok) {
            const data = await response.json();
            const token = data.token;
            const role = data.role;

            localStorage.setItem('jwtToken', token);
            localStorage.setItem('role', role);


            console.log('Verification Response:', data);


        } else {
            throw new Error('Network response was not ok.');
        }

        try {

            const jwtToken = localStorage.getItem('jwtToken');
            const role = localStorage.getItem('role');
            console.log(jwtToken)
            if (role == "admin") {
                const response = await fetch('/admin/admin_home', {
                    method: 'GET',
                    headers: {
                        'Authorization': `${jwtToken}`
                    }
                });
                if (response.ok) {
                    // If the response status is in the range 200-299
                    const htmlContent = await response.text();

                    //Update the content of a specific HTML element with the fetched HTML
                    document.open();
                    document.write(htmlContent);
                    document.close();
                } else if (response.status >= 400 && response.status < 500) {
                    // Handle client-side errors (4xx errors)
                    // For example, display an error message or handle accordingly
                    console.error('Client-side error:', response.status);
                } else if (response.status >= 500 && response.status < 600) {
                    // Handle server-side errors (5xx errors)
                    // For example, display an error message or handle accordingly
                    console.error('Server-side error:', response.status);
                } else {
                    // Handle other cases where response.ok is false but the status code is not in 4xx or 5xx range
                    console.error('Unexpected error:', response.status);
                }

            } else {

                const response = await fetch('/user/user_home', {
                    method: 'GET',
                    headers: {
                        'Authorization': `${jwtToken}`
                    }
                });
                if (response.ok) {
                    // If the response status is in the range 200-299
                    const htmlContent = await response.text();

                    //Update the content of a specific HTML element with the fetched HTML
                    document.open();
                    document.write(htmlContent);
                    document.close();

                } else if (response.status >= 400 && response.status < 500) {
                    // Handle client-side errors (4xx errors)
                    // For example, display an error message or handle accordingly
                    console.error('Client-side error:', response.status);
                } else if (response.status >= 500 && response.status < 600) {
                    // Handle server-side errors (5xx errors)
                    // For example, display an error message or handle accordingly
                    console.error('Server-side error:', response.status);
                } else {
                    // Handle other cases where response.ok is false but the status code is not in 4xx or 5xx range
                    console.error('Unexpected error:', response.status);
                }
            }


        } catch (error) {
            // Handle any errors during fetch or navigation
            console.error('Fetch error:', error);
        }
        // Handle the response from the backend as needed
    } catch (error) {
        console.error('Error:', error);
    }

}

function stringToHex(str) {
    let hex = '';
    for (let i = 0; i < str.length; i++) {
        const charCode = str.charCodeAt(i).toString(16);
        hex += charCode.length === 1 ? '0' + charCode : charCode;
    }
    return '0x' + hex;
}

