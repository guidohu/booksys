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
                cbFailed(response)
            }
        })
    }

    static getMyUser(cbSuccess, cbFailed){
        Login.getData('/api/user.php?action=get_my_user')
            .then(response => {
                console.log(response)
                cbSuccess(response)
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