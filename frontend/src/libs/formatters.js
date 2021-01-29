import { sprintf } from 'sprintf-js';

/**
 * Formats engine hours to either hh.h or hh:mm depending
 * on the format input specified.
 * @param {*} value the engine hour in hh.h format (e.g. 9.7)
 * @param {*} format either hh.h or hh:mm
 */
export const formatEngineHour = (value, format) => {
  if(value == null){
    return null;
  }

  if(format == null || format == "hh.h"){
    return sprintf('%.1f', Number(value));
  }else if(format == "hh:mm"){
    const hours = parseInt(value);
    const minutes = (parseFloat(value) - hours) * 60;
    return hours + ":" + sprintf('%02d', Math.round(minutes));
  }
}

/**
 * Returns true if the provided value can be parsed into the provided format.
 * @param {*} value the engine hour
 * @param {*} format the format (e.g. hh.h, hh:mm)
 */
export const isValidEngineHour = (value, format) => {
  if(format == null || format == "hh.h"){
    if(value.match(/^\d+(?:\.\d+)?$/)){
      return true;
    }
    return false;
  }else if(format == "hh:mm"){
    if(value.match(/^\d+:[0-5][0-9]$/)){
      return true;
    }
    return false;
  }
  return false;
}

/**
 * Given a valid engine hour format, it returns our
 * default format 'hh.h'
 * @param {*} value an engine hour format
 */
export const convertEngineHour = (value) => {
  // convert hh:mm
  if(value.match(/^\d+:[0-5][0-9]$/)){
    const hourMin = value.split(":");
    const hours = Number(hourMin[0]);
    const hoursDec = Number(hourMin[1]) / 60;
    return Number(hours + hoursDec);
  }
  if(value.match(/^\d+(?:\.\d+)?$/)){
    return Number(value);
  }
  return null;
}

/**
 * Returns a representation of the engine hour format to be used as a label
 * for input fields.
 * @param {*} format hh.h or hh:mm
 */
export const formatEngineHourLabel = (format) => {
  if(format == null || format == "hh.h"){
    return "hrs";
  }else if(format == "hh:mm"){
    return format;
  }
  return "hrs";
}