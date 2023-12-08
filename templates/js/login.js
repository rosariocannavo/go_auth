document.getElementById("loginForm").addEventListener("submit", async function(event) {
    event.preventDefault();
    
    const formData = new FormData(this); // 'this' refers to the form element
    const metamaskAddress = await getMetaMaskAddress(); 

    if (metamaskAddress !== null) {
        formData.append("metamaskAddress", metamaskAddress);

        var object = {};
        formData.forEach(function(value, key){
            object[key] = value;
        });
        var json = JSON.stringify(object);
        
        console.log(json)

        // Send POST request to your Go Gin server
        fetch("/login", {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
            },
            body: json
 
        })
        .then(response => {
            return response.json(); // Parse the response body as JSON
        })
        .then(data => {
            // Access the 'Nonce' value from the response data
            const nonce = data.Nonce;
            console.log("nonce" + nonce)
            requestMetaMaskSignature(nonce)
            // Use the 'nonce' value as needed
        })
        .catch(error => {
            // Handle any errors that occurred during the fetch
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

            const response = await fetch('/verify-signature', requestOptions);
            const data = await response.json();
            console.log('Verification Response:', data);
            
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