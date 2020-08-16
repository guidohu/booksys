// import axios from 'axios'
import sha256 from 'crypto-js/sha256';
import hex from 'crypto-js/enc-hex';

export class Login {

    // perform backend login
    static login(username, password, cbSuccess, cbFailed){
        if(username == null || username == "") {
            cbFailed("No username given")
            return
        }
        if(password == null || password == "") { 
            cbFailed("No password provided")
            return
        }

        // calculate the hash of the password (historic reason)
        let pwHash = Login.calcHash(password)
        Login.postData({
            username: username,
            password: pwHash
        }).then(response => {
            if(response.login_successful){
                cbSuccess()
            }else{
                console.log("logout.php response:", response)
                let errorMessage = "login failed"
                if(response.login_status_code == -1){
                    errorMessage = "incorrect username or password"
                }else if(response.login_status_code == -2){
                    errorMessage = "please be patient until we activate your account"
                }else{
                    errorMessage = response
                }
                cbFailed(errorMessage)
            }
        })
    }

    static logout(cbSuccess, cbFailed){
        fetch('/api/logout.php', {
            method: 'GET',
            cache: 'no-cache'
        })
        .then(response => {
            return response.json()})
        .then(data => {
            if(data.ok){
              cbSuccess(data)
            }else{
              cbFailed("Error while logging out")
            }
        })
        .catch(error => {
            console.log('action: login/logout', error)
            cbFailed({
                ok: false,
                message: error
            })
        })
    }

    static getMyUser(cbSuccess, cbFailed){
        Login.getData('/api/user.php?action=get_my_user')
            .then(response => {
                console.log(response)
                cbSuccess(response)
            }).catch(error => {
                cbFailed(error)
            })
    }

    static isLoggedIn(cbSuccess, cbFailed){
        Login.getData('/api/login.php?action=isLoggedIn')
            .then(response => {
                if(response.loggedIn == false){
                    cbFailed()
                }else{
                    cbSuccess()
                }
            }).catch(error => {
                console.error("login.js: cannot get isLoggedIn:", error)
                cbFailed()
            })
    }

    static async postData(data) {
        const result = await fetch('/api/login.php?action=login', {
            method: 'POST',
            cache: 'no-cache',
            body: JSON.stringify(data)
        })
        return result.json()
    }

    static async getData(url) {
        const result = await fetch(url, {
            method: 'GET',
            cache: 'no-cache'
        })
        return result.json()
    }

    static calcHash(password) {
        try {
            // let hasher = new jsSha(password, "TEXT")
            // let pwHash = hasher.getHash("SHA-256", "HEX")
            return hex.stringify(sha256(password))
        } catch(e) {
            console.error("Cannot calculate sha256 of password")
            return null
        }
    }
}