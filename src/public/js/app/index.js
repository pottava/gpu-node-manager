function signout() {
    firebase.auth().signOut().then(() => {
        console.log('Signed out');
    })
}
