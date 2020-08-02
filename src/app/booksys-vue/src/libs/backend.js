import axios from 'axios'

export class BooksysBackend {

    static getStatus(cb) {
        axios.get('backend/api/backend.php?action=get_status')
            .then(response => {
                cb(response.data)
            })
            .catch(error => {
                console.error("Cannot connect to backend: ", error)
                cb({ ok: false, message: "Cannot connect to backend"})
            })
    }

    static isUnavailable(responseData) {
        if (responseData.message != "no database configured"
            && responseData.message != "no user configured"
            && responseData.message != "no config")
        {
            return true
        }
    }

    static needsSetup(responseData) {
        if (responseData.message == "no database configured"
            || responseData.message == "no user configured"
            || responseData.message == "no config")
        {
            return true
        }
    }
}