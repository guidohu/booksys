export class Configuration {

  static getCustomizationParameters(cbSuccess, cbFailure){
    fetch('/api/configuration.php?action=get_customization_parameters', {
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
      console.log('Configuration/get_customization/parameters', error)
      cbFailure(
        error
      )
    })
  }

  static needsDbUpdate(cbSuccess, cbFailure){
    fetch('/api/backend.php?action=admin_check_database_update', {
      method: 'GET',
      cache: 'no-cache'
    })
    .then(response => {
      return response.json()
    })
    .then(data => {
      cbSuccess(data)
    })
    .catch(error =>{
      cbFailure(error)
    })
  }

  static getDbVersion(cbSuccess, cbFailure){
    fetch('/api/backend.php?action=get_version', {
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
      cbFailure(error)
    })
  }

  static updateDb(cbSuccess, cbFailure){
    fetch('/api/backend.php?action=update_database', {
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
      cbFailure(error)
    })
  }
}