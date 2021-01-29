import { sprintf } from 'sprintf-js';

/**
 * Formats engine hours to either hh.h or hh:mm depending
 * on the format input specified.
 * @param {*} value the engine hour in hh.h format (e.g. 9.7)
 * @param {*} format either hh.h or hh:mm
 */
export const formatEngineHour = (value, format) => {
  if(format == null || format == "hh.h"){
    console.log("value is:", value);
    return sprintf('%.1f', value);
  }else if(format == "hh:mm"){
    const hours = parseInt(value);
    const minutes = (parseFloat(value) - hours) * 60;
    return hours + ":" + sprintf('%02d', minutes);
  }
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
}