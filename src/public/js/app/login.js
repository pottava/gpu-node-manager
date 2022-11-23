var ui = new firebaseui.auth.AuthUI(firebase.auth());
ui.start('#firebase-auth', {
    signInSuccessUrl: '/',
    signInOptions: [
        firebase.auth.EmailAuthProvider.PROVIDER_ID
    ]
});
