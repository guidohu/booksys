
/*
  Formats engine hours to either hh.h or hh:mm depending
  on the format input specified.
*/
export const formatEngineHour = (value, format) => {
  if(format == null || format == "hh.h"){
    return value;
  }else if(format == "hh:mm"){
    const hours = parseInt(value);
    const minutes = (parseFloat(value) - hours) * 60;
    return hours + ":" + minutes;
  }
}