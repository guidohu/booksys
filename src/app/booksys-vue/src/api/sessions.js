
export default class Sessions {

  static getSessions(dateStart, dateEnd, cbSuccess, cbFailure){
    let query = {
      start: dateStart.format("X"),
      end: dateEnd.format("X")
    }
    fetch('/api/booking.php?action=get_booking_day', {
      method: 'POST',
      cache: 'no-cache',
      body: JSON.stringify(query)
    })
    .then(response => {
      return response.json()
    })
    .then(data => {
      cbSuccess(data)
    })
    .catch(error => {
      console.log('Sessions/getSessions', error)
      cbFailure(
        error
      )
    })
  }
}