//registration form
document.getElementById("registerForm").addEventListener("submit", async function(event) {
    event.preventDefault();
    
    const formData = new FormData(this); // 'this' refers to the form element
    const metamaskAddress = await getMetaMaskAddress(); 

    //TODO: check password and confirm password
    if (metamaskAddress !== null) {
        formData.append("metamaskAddress", metamaskAddress);

        var object = {};
        formData.forEach(function(value, key){
            if(key != 'passwordConfirm') 
                object[key] = value;
        });
        var json = JSON.stringify(object);
        
        console.log(json)

        // Send POST request to your Go Gin server
        fetch("/registration", {
            method: "POST",
            //body: formData
            headers: {
                'Content-Type': 'application/json'
            },
            body: json
 
        })
        .then(response => {
            // Handle the response as needed
            console.log(response);
            // You can redirect or perform other actions based on the response
        })
        .catch(error => {
            console.error("Error:", error);
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
