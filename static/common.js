
// AUTH
var firebaseConfig = {
    apiKey: "AIzaSyCaUioWzhgjKbTj9nFEhtmaN4KbRXxotbk",
    authDomain: "grabtest-828e6.firebaseapp.com",
    projectId: "grabtest-828e6",
    storageBucket: "grabtest-828e6.appspot.com",
    messagingSenderId: "1050687107203",
    appId: "1:1050687107203:web:bd2a510b1266299e922bb6"
};


// Initialize Firebase
firebase.initializeApp(firebaseConfig);

const auth = firebase.auth();

auth.onAuthStateChanged(function (user) {
    if (user) {

        var email = user.email;
        user.getIdToken().then(function (token) {
            localStorage.setItem("token", token);
        });

        console.log("User In")
        


    } else {
        console.log("User Out")
        localStorage.removeItem("token");
        window.location.href = "/";
    }



});

function signOut() {
    auth.signOut();
    var token = localStorage.getItem("token");
    localStorage.removeItem("token");
    axios.get('/log/api/out', { 'headers': { 'Authorization': 'Bearer ' + token } });
}

function getGridData() {
    const zgRef = document.querySelector('zing-grid');
    const gridData = zgRef.getData({
        csv: true
    });
    console.log('--- Getting Data From Grid: ---', gridData);
    alert('Check console for exported data!');
}