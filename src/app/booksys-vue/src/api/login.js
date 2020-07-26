import axios from 'axios'
import sha256 from 'crypto-js/sha256';
import hex from 'crypto-js/enc-hex';

export class Login {

    // perform backend login
    static login(username, password, cbSuccess, cbFailed){
        console.log("login called")
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
        axios.post('http://localhost/api/login.php?action=login', {
            username: username,
            password: pwHash,
          })
          .then(function (response) {
            console.log(response);
            if(cbSuccess != null) {
                cbSuccess()
            }
          })
          .catch(function (error) {
            console.log(error);
            if(cbFailed != null) {
                cbFailed()
            }
          });
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