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
}