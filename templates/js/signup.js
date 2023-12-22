//registration form
document.getElementById("registerForm").addEventListener("submit", async function (event) {
    event.preventDefault();

    const formData = new FormData(this); // 'this' refers to the form element
    const metamaskAddress = await getMetaMaskAddress();


    const password = formData.get('password');
    const confirmPassword = formData.get('passwordConfirm');

    console.log("pwd:" + password);
    console.log("confirm pwd:" + confirmPassword);

    if (password === confirmPassword) {
        document.getElementById('response').innerHTML = '<p>Passwords match!</p>';

        document.getElementById('passwordRegister').style.border = '2px solid green';
        document.getElementById('passwordConfirm').style.border = '2px solid green';

        if (metamaskAddress !== null) {
            formData.append("metamaskAddress", metamaskAddress);

            var object = {};
            formData.forEach(function (value, key) {
                if (key != 'passwordConfirm')
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
                    if (response.status === 403) {
                        console.log("user already present");
                        document.getElementById('username').style.border = '2px solid red';
                        document.getElementById('passwordRegister').style.border = '2px solid red';
                        document.getElementById('passwordConfirm').style.border = '2px solid red';


                        document.getElementById('response').innerHTML = '<p>User or Address already present.</p>';

                    } else {
                        window.location.href = '/';
                    }
                })
                .catch(error => {
                    console.error("Error:", error);
                });


        } else {
            console.error("Metamask address not available");
            // Handle the case when Metamask address is not available
        }


    } else {
        document.getElementById('passwordRegister').style.border = '2px solid red';
        document.getElementById('passwordConfirm').style.border = '2px solid red';

        // Display error message
        document.getElementById('response').innerHTML = '<p>Passwords do not match. Please retry.</p>';
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
