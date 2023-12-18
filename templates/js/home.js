document.getElementById("setButton").addEventListener('click', function() {
    
   // const token = localStorage.getItem('jwtToken');

    const queryParams = new URLSearchParams(window.location.search);
    const encodedData = queryParams.get('data');

    // Decode the data
    const token = decodeURIComponent(encodedData);

    console.log("retrieved token" + token)
    
    const newValue = Math.floor(Math.random() * (100 - 1)) + 1;

    const url = `http://localhost:8080/app/interactWithContract?account=${account}&newvalue=${newValue}`;   
    fetch(url, {
        method: "GET",
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `${token}` // Add Authorization header with JWT token
        },
       // body: json
    })
    .then(response => {
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
    return response.json();
    })
    .then(data => {
    // Process the response data here
    console.log(data);
    })
    .catch(error => {
    // Handle errors here
    console.error('There was a problem with the fetch operation:', error);
    });
});