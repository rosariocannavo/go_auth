var metamaskAddress  = null;

document.getElementById("loginForm").addEventListener("submit", async function(event) {
    event.preventDefault();
    
    const formData = new FormData(this); // 'this' refers to the form element
    metamaskAddress = await getMetaMaskAddress(); 

    if (metamaskAddress !== null) {
        formData.append("metamaskAddress", metamaskAddress);

        var object = {};
        formData.forEach(function(value, key){
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
            
                localStorage.setItem('jwtToken', token);
            
                console.log('Verification Response:', data);
            
            
            } else {
                throw new Error('Network response was not ok.');
            }

            try {

                const jwtToken = localStorage.getItem('jwtToken');

                const response = await fetch('/home', {
                    method: 'GET',
                    headers: {
                    'Authorization': `${jwtToken}`
                    }
                });
                
                if (response.ok) {
                    // If the response is successful, you can navigate to the '/home' page
                    window.location.href = '/home';
                } else {
                    throw new Error('Network response was not ok.');
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

