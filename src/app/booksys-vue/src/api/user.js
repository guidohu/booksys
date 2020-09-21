
export default class User {

    static getHeats(cbSuccess, cbFailure){
      fetch('/api/user.php?action=get_my_user_heats', {
        method: 'GET',
        cache: 'no-cache'
      })
      .then(response => {
        return response.json()
      })
      .then(data => {
        cbSuccess(data)
      })
      .catch(error => {
        console.log('User/getHeats', error)
        cbFailure(
          error
        )
      })
    }

    static changeUserProfile(profileData, cbSuccess, cbFailure){
      fetch('/api/user.php?action=change_my_user_data', {
        method: 'POST',
        cache: 'no-cache',
        body: JSON.stringify(profileData)
      })
      .then(response => {
        return response.json()
      })
      .then(data => {
        cbSuccess(data)
      })
      .catch(error => {
        console.log('User/changeUserProfile', error)
        cbFailure(error)
      })
    }
}